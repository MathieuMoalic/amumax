package api

import (
	"github.com/MathieuMoalic/amumax/engine"
)

var SelectedRegion int

func init() {
	SelectedRegion = 0
}

type Field struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Parameters struct {
	Regions        []int   `json:"regions"`
	Fields         []Field `json:"fields"`
	SelectedRegion int     `json:"selectedRegion"`
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
		}
		fields = append(fields, field)
	}
	return fields
}
