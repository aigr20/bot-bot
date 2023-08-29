package opendota

type TotalField struct {
	FieldName string  `json:"field"`
	Matches   int     `json:"n"`
	Sum       float64 `json:"sum"`
}

type TotalResponse = []TotalField
