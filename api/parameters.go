package api

import (
	"github.com/MathieuMoalic/amumax/engine"
)

type Field struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Parameters struct {
	Regions []int   `json:"regions"`
	Fields  []Field `json:"fields"`
}

func newParameters() *Parameters {
	regionsIndices := engine.Regions.GetExistingIndices()
	parameters := Parameters{
		Regions: regionsIndices,
	}
	fields := make([]Field, 0)
	for _, param := range engine.Params {
		field := Field{
			Name:        param.Name,
			Value:       param.Value(engine.Frontend.ParameterRegion),
			Description: param.Description,
		}
		fields = append(fields, field)
	}
	parameters.Fields = fields
	return &parameters
}
