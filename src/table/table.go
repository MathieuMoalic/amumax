package table

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"github.com/MathieuMoalic/amumax/src/fsutil"
	"github.com/MathieuMoalic/amumax/src/log"
	"github.com/MathieuMoalic/amumax/src/quantity"
	"github.com/MathieuMoalic/amumax/src/script"
	"github.com/MathieuMoalic/amumax/src/solver"
	"github.com/MathieuMoalic/amumax/src/utils"
)

// the Table is kept in RAM and used for the API
type Table struct {
	log            *log.Logs
	fs             *fsutil.FileSystem
	solver         *solver.Solver
	quantities     []quantity.Quantity
	columns        []column
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	mu             sync.Mutex
	initialized    bool
	lastSavedHash  string // Hash of the last saved state
}

type column struct {
	name   string
	unit   string
	writer *bufio.Writer
	file   *os.File
}

// NewTable creates a new empty table, does not save it to disk
func NewTable(solver *solver.Solver, log *log.Logs, fs *fsutil.FileSystem, script *script.ScriptParser) *Table {
	t := &Table{
		solver:         solver,
		log:            log,
		fs:             fs,
		quantities:     []quantity.Quantity{},
		columns:        []column{},
		Data:           make(map[string][]float64),
		AutoSavePeriod: 0,
		AutoSaveStart:  0,
		Step:           -1,
		mu:             sync.Mutex{},
		initialized:    false,
		lastSavedHash:  "",
	}
	script.RegisterFunction("TableAutoSave", t.tableAutoSave)
	script.RegisterFunction("TableAdd", t.tableAdd)
	script.RegisterFunction("TableAddAs", t.tableAddAs)
	script.RegisterFunction("TableSave", t.tableSave)
	return t
}

// generateHash creates a hash based on the current Step and column names.
func (ts *Table) generateHash() string {
	hash := sha256.New()
	fmt.Fprint(hash, ts.Step) // Include the Step in the hash
	for _, col := range ts.columns {
		hash.Write([]byte(col.name)) // Include column names
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// InitTable creates default columns and saves them to disk
func (ts *Table) initTable() {
	err := ts.fs.Remove("table")
	ts.log.PanicIfError(err)
	err = ts.fs.CreateZarrGroup("table/")
	ts.log.PanicIfError(err)
	ts.addColumn("step", "")
	ts.addColumn("t", "s")
	ts.initialized = true
}

func (ts *Table) writeToBuffer() {
	buf := []float64{}
	buf = append(buf, float64(ts.Step))
	buf = append(buf, ts.solver.Time)
	for _, q := range ts.quantities {
		buf = append(buf, q.Average()...)
	}

	ts.mu.Lock() // Lock before modifying shared data
	defer ts.mu.Unlock()

	for i, b := range buf {
		// Convert float64 to bytes
		data := utils.Float64ToBytes(b)

		// Write directly to the buffered writer
		_, err := ts.columns[i].writer.Write(data)
		if err != nil {
			ts.log.PanicIfError(err)
		}

		// Update in-memory data
		ts.Data[ts.columns[i].name] = append(ts.Data[ts.columns[i].name], b)
	}
}

// FlushToFile writes the buffered data to disk
func (ts *Table) FlushToFile() {
	if ts.Step == -1 {
		return
	}
	// Check if the table state has changed
	currentHash := ts.generateHash()
	if currentHash == ts.lastSavedHash {
		ts.log.Debug("Table state has not changed, skipping save.")
		return
	}
	for i := range ts.columns {
		// Update zarray if necessary, it is not buffered at the moment
		err := ts.fs.SaveFileTableZarray("table/"+ts.columns[i].name, ts.Step)
		if err != nil {
			ts.log.PanicIfError(err)
		}
		err = ts.columns[i].writer.Flush()
		if err != nil {
			ts.log.PanicIfError(err)
		}
	}
	ts.lastSavedHash = currentHash
}

func (ts *Table) exists(q quantity.Quantity, name string) bool {
	suffixes := []string{"x", "y", "z"}
	for _, col := range ts.columns {
		if q.NComp() == 1 {
			if col.name == name {
				return true
			}
		} else {
			for comp := 0; comp < q.NComp(); comp++ {
				if name+suffixes[comp] == col.name {
					return true
				}
			}
		}
	}
	return false
}

func (ts *Table) addColumn(name, unit string) {
	if !ts.initialized {
		ts.initTable()
	}
	err := ts.fs.Mkdir("table/" + name)
	ts.log.PanicIfError(err)
	writer, file, err := ts.fs.Create("table/" + name + "/0")
	ts.log.PanicIfError(err)
	ts.columns = append(ts.columns, column{name: name, unit: unit, writer: writer, file: file})
}

func (ts *Table) tableSave() {
	if len(ts.columns) == 0 {
		ts.log.Warn("No columns in table, not saving.")
	}
	ts.Step += 1
	ts.writeToBuffer()
}

func (ts *Table) tableAdd(q quantity.Quantity) {
	ts.tableAddAs(q, q.Name())
}

func (ts *Table) tableAddAs(q quantity.Quantity, name string) {
	suffixes := []string{"x", "y", "z"}
	if ts.Step != -1 {
		ts.log.Warn("You cannot add a new quantity to the table after the simulation has started. Ignoring.")
	}
	if len(ts.columns) == 0 {
		ts.log.Warn("No columns in table, not saving.")
	}

	if ts.exists(q, name) {
		ts.log.Warn("%s is already in the table. Ignoring.", name)
		return
	}
	ts.quantities = append(ts.quantities, q)
	if q.NComp() == 1 {
		ts.addColumn(name, q.Unit())
	} else {
		for comp := 0; comp < q.NComp(); comp++ {
			ts.addColumn(name+suffixes[comp], q.Unit())
		}
	}
}

func (ts *Table) tableAutoSave(period float64) {
	ts.AutoSaveStart = ts.solver.Time
	ts.AutoSavePeriod = period
}

func (ts *Table) Close() {
	for _, col := range ts.columns {
		err := col.writer.Flush()
		ts.log.PanicIfError(err)
		err = col.file.Close()
		ts.log.PanicIfError(err)
	}
}
