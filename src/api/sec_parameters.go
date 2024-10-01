package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/util"
	"github.com/labstack/echo/v4"
)

var SelectedRegion int

func init() {
	SelectedRegion = 0
}

type Field struct {
	Name        string `msgpack:"name"`
	Value       string `msgpack:"value"`
	Description string `msgpack:"description"`
	Changed     bool   `msgpack:"changed"`
}

func (f *Field) IsDefault(value string) bool {
	return true
}

type Parameters struct {
	Regions        []int   `msgpack:"regions"`
	Fields         []Field `msgpack:"fields"`
	SelectedRegion int     `msgpack:"selectedRegion"`
}

func newParameters() *Parameters {
	return &Parameters{
		Regions:        engine.Regions.GetExistingIndices(),
		Fields:         getFields(),
		SelectedRegion: SelectedRegion,
	}
}

func getFields() []Field {
	fields := make([]Field, 0)
	for _, param := range engine.Params {
		field := Field{
			Name:        param.Name,
			Value:       param.Value(SelectedRegion),
			Description: param.Description,
			Changed:     engine.QuantityChanged[param.Name],
		}
		fields = append(fields, field)
	}
	return fields
}

func postSelectParameterRegion(c echo.Context) error {
	type Request struct {
		SelectedRegion int `msgpack:"selectedRegion"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		util.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	SelectedRegion = req.SelectedRegion
	preview.Refresh = true
	broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
