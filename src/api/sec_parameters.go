package api

import (
	"net/http"

	"github.com/MathieuMoalic/amumax/src/engine_old"
	"github.com/MathieuMoalic/amumax/src/log_old"
	"github.com/labstack/echo/v4"
)

type Field struct {
	Name        string `msgpack:"name"`
	Value       string `msgpack:"value"`
	Description string `msgpack:"description"`
	Changed     bool   `msgpack:"changed"`
}

func (f *Field) IsDefault(value string) bool {
	return true
}

type ParametersState struct {
	ws             *WebSocketManager
	Regions        []int   `msgpack:"regions"`
	Fields         []Field `msgpack:"fields"`
	SelectedRegion int     `msgpack:"selectedRegion"`
}

func initParameterAPI(e *echo.Group, ws *WebSocketManager) *ParametersState {
	parametersState := ParametersState{
		ws:             ws,
		Regions:        engine_old.Regions.GetExistingIndices(),
		SelectedRegion: 0,
	}
	parametersState.getFields()
	e.POST("/api/parameter/selected-region", parametersState.postSelectParameterRegion)
	return &parametersState
}
func (s *ParametersState) Update() {
	s.getFields()
}

func (s *ParametersState) getFields() {
	fields := make([]Field, 0)
	for _, param := range engine_old.Params {
		field := Field{
			Name:        param.Name,
			Value:       param.Value(s.SelectedRegion),
			Description: param.Description,
			Changed:     engine_old.QuantityChanged[param.Name],
		}
		fields = append(fields, field)
	}
	s.Fields = fields
}

func (s *ParametersState) postSelectParameterRegion(c echo.Context) error {
	type Request struct {
		SelectedRegion int `msgpack:"selectedRegion"`
	}
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log_old.Log.Err("%v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	s.SelectedRegion = req.SelectedRegion
	s.ws.engineState.Preview.Refresh = true
	s.ws.broadcastEngineState()
	return c.JSON(http.StatusOK, nil)
}
