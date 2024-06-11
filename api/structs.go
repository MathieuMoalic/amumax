package api

import (
	"github.com/MathieuMoalic/amumax/engine"
)

type Header struct {
	Path     string  `json:"path"`
	Progress float64 `json:"progress"`
	Status   string  `json:"status"`
}
type Solver struct {
	Type       string  `json:"type"`
	Steps      int     `json:"steps"`
	Time       float64 `json:"time"`
	Dt         float64 `json:"dt"`
	ErrPerStep float64 `json:"errPerStep"`
	MaxTorque  float64 `json:"maxTorque"`
	Fixdt      float64 `json:"fixdt"`
	Mindt      float64 `json:"mindt"`
	Maxdt      float64 `json:"maxdt"`
	Maxerr     float64 `json:"maxerr"`
}
type Console struct {
	Hist string `json:"hist"`
}
type Mesh struct {
	Dx   float64 `json:"dx"`
	Dy   float64 `json:"dy"`
	Dz   float64 `json:"dz"`
	Nx   int     `json:"Nx"`
	Ny   int     `json:"Ny"`
	Nz   int     `json:"Nz"`
	Tx   float64 `json:"Tx"`
	Ty   float64 `json:"Ty"`
	Tz   float64 `json:"Tz"`
	PBCx int     `json:"PBCx"`
	PBCy int     `json:"PBCy"`
	PBCz int     `json:"PBCz"`
}

type TablePlotData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type TablePlot struct {
	AutoSaveInterval float64         `json:"autoSaveInterval"`
	Columns          []string        `json:"columns"`
	XColumn          string          `json:"xColumn"`
	YColumn          string          `json:"yColumn"`
	Data             []TablePlotData `json:"data"`
}

func newTablePlot() *TablePlot {
	tablePlotData := make([]TablePlotData, len(engine.Table.Data[engine.Tableplot.X]))
	for i := range tablePlotData {
		tablePlotData[i] = TablePlotData{
			X: engine.Table.Data[engine.Tableplot.X][i],
			Y: engine.Table.Data[engine.Tableplot.Y][i],
		}
	}
	data := TablePlot{
		AutoSaveInterval: engine.Table.AutoSavePeriod,
		Columns:          engine.Table.GetTableNames(),
		XColumn:          engine.Tableplot.X,
		YColumn:          engine.Tableplot.Y,
		Data:             tablePlotData,
	}
	return &data
}

type EngineState struct {
	Header    Header     `json:"header"`
	Solver    Solver     `json:"solver"`
	Console   Console    `json:"console"`
	Mesh      Mesh       `json:"mesh"`
	Params    Parameters `json:"parameters"`
	TablePlot TablePlot  `json:"tablePlot"`
}

func NewEngineState() *EngineState {
	status := ""
	if engine.Pause {
		status = "paused"
	} else {
		status = "running"

	}
	engineState := EngineState{
		Header: Header{
			Path:   engine.OD(),
			Status: status,
		},
		Solver: Solver{
			Type:       engine.Solvernames[engine.Solvertype],
			Steps:      engine.NSteps,
			Time:       engine.Time,
			Dt:         engine.Dt_si,
			ErrPerStep: engine.LastErr,
			MaxTorque:  engine.LastTorque,
			Fixdt:      engine.FixDt,
			Mindt:      engine.MinDt,
			Maxdt:      engine.MaxDt,
			Maxerr:     engine.MaxErr,
		},
		Console: Console{
			Hist: engine.Hist,
		},
		Mesh: Mesh{
			Dx:   engine.Dx,
			Dy:   engine.Dy,
			Dz:   engine.Dz,
			Nx:   engine.Nx,
			Ny:   engine.Ny,
			Nz:   engine.Nz,
			Tx:   engine.Tx,
			Ty:   engine.Ty,
			Tz:   engine.Tz,
			PBCx: engine.PBCx,
			PBCy: engine.PBCy,
			PBCz: engine.PBCz,
		},
		Params:    *newParameters(),
		TablePlot: *newTablePlot(),
	}
	return &engineState
}
