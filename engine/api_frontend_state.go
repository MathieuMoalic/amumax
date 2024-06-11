package engine

type FrontendState struct {
	DisplayQuantity  string `json:"displayQuantity"`
	DisplayComponent int    `json:"displayComponent"`
	DisplayLayer     int    `json:"displayLayer"`
	ParameterRegion  int    `json:"parameterRegion"`
}

var Frontend FrontendState
var Params map[string]Field

func init() {
	Frontend = FrontendState{
		DisplayQuantity:  "m",
		DisplayComponent: 0,
		DisplayLayer:     0,
		ParameterRegion:  0,
	}
}

type Field struct {
	Name        string           `json:"name"`
	Value       func(int) string `json:"value"`
	Description string           `json:"description"`
}

func AddParameter(name string, value interface{}, doc string) {
	if Params == nil {
		Params = make(map[string]Field)
	}
	if v, ok := value.(*RegionwiseScalar); ok {
		Params[name] = Field{
			name,
			v.GetRegionToString,
			doc,
		}
	}
	if v, ok := value.(*RegionwiseVector); ok {
		Params[name] = Field{
			name,
			v.GetRegionToString,
			doc,
		}
	}
	if v, ok := value.(*inputValue); ok {
		Params[name] = Field{
			name,
			v.GetRegionToString,
			doc,
		}
	}
}
