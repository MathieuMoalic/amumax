package engine

import (
	"time"

	"github.com/MathieuMoalic/amumax/httpfs"
	"github.com/MathieuMoalic/amumax/util"
	"github.com/MathieuMoalic/amumax/zarr"
)

func init() {
	DeclFunc("TableSave", ZTableSave, "Save the data table right now.")
	DeclFunc("TableAdd", ZTableAdd, "Save the data table periodically.")
	DeclFunc("TableAddAs", ZTableAddAs, "Save the data table periodically.")
	DeclFunc("TableAutoSave", ZTableAutoSave, "Save the data table periodically.")
	ZTables = ZTablesStruct{Data: make(map[string][]float64), Step: -1, AutoSavePeriod: 0.0, FlushInterval: 5 * time.Second}
}

var ZTables ZTablesStruct

// the Table is kept in RAM and used for the API
type ZTablesStruct struct {
	qs             []Quantity
	tables         []ZTable
	Data           map[string][]float64 `json:"data"`
	AutoSavePeriod float64              `json:"autoSavePeriod"`
	AutoSaveStart  float64              `json:"autoSaveStart"`
	Step           int                  `json:"step"`
	FlushInterval  time.Duration        `json:"flushInterval"`
}

type ZTable struct {
	Name   string
	buffer []byte
	io     httpfs.WriteCloseFlusher
}

func (ts *ZTablesStruct) WriteToBuffer() {
	buf := []float64{}
	// always save the current time
	buf = append(buf, Time)
	// for each quantity we append each component to the buffer
	for _, q := range ts.qs {
		buf = append(buf, AverageOf(q)...)
	}
	// size of buf should be same as size of []Ztable
	for i, b := range buf {
		ts.tables[i].buffer = append(ts.tables[i].buffer, zarr.Float64ToByte(b)...)
		ts.Data[ts.tables[i].Name] = append(ts.Data[ts.tables[i].Name], b)
		// ts.Tables[i].Data = append(ts.Tables[i].Data, b)
	}
}
func (ts *ZTablesStruct) Flush() {
	for i := range ts.tables {
		ts.tables[i].io.Write(ts.tables[i].buffer)
		ts.tables[i].buffer = []byte{}
		ts.tables[i].io.Flush()
		zarr.SaveFileTableZarray(OD()+"table/"+ts.tables[i].Name+"/.zarray", ts.Step)
	}
}
func (ts *ZTablesStruct) NeedSave() bool {
	return ts.AutoSavePeriod != 0 && (Time-ts.AutoSaveStart)-float64(ts.Step)*ts.AutoSavePeriod >= ts.AutoSavePeriod
}

func TableInit() {
	httpfs.Remove(OD() + "table")
	zarr.MakeZgroup("table", OD(), &zGroups)
	err := httpfs.Mkdir(OD() + "table/t")
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/t/0")
	util.FatalErr(err)
	ZTables.tables = append(ZTables.tables, ZTable{"t", []byte{}, f})
	go AutoFlush()

}

func AutoFlush() {
	for {
		ZTables.Flush()
		time.Sleep(ZTables.FlushInterval)
	}
}

func ZTableSave() {
	if len(ZTables.tables) == 0 {
		util.Fatal("Error: Add a variable to the table before saving")
	}
	ZTables.Step += 1
	ZTables.WriteToBuffer()
}

func CreateTable(name string) ZTable {
	err := httpfs.Mkdir(OD() + "table/" + name)
	util.FatalErr(err)
	f, err := httpfs.Create(OD() + "table/" + name + "/0")
	util.FatalErr(err)
	return ZTable{Name: name, buffer: []byte{}, io: f}
}

func ZTableAdd(q Quantity) {
	ZTableAddAs(q, NameOf(q))
}
func ZTableAddAs(q Quantity, name string) {
	if len(ZTables.tables) == 0 {
		TableInit()
	}
	for _, z := range ZTables.tables {
		if name == z.Name {
			return
		}
	}
	if ZTables.Step != -1 {
		util.Fatal("Add Table Quantity BEFORE you save the table for the first time")
	}
	ZTables.qs = append(ZTables.qs, q)
	if q.NComp() == 1 {
		ZTables.tables = append(ZTables.tables, CreateTable(name))
	} else {
		suffixes := []string{"x", "y", "z"}
		for comp := 0; comp < q.NComp(); comp++ {
			ZTables.tables = append(ZTables.tables, CreateTable(name+suffixes[comp]))
		}
	}
}

func ZTableAutoSave(period float64) {
	ZTables.AutoSaveStart = Time
	ZTables.AutoSavePeriod = period
}
