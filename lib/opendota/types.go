package opendota

type HeroAlias struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AliasMap = map[string]HeroAlias

type TotalField struct {
	FieldName string  `json:"field"`
	Matches   int     `json:"n"`
	Sum       float64 `json:"sum"`
}

type TotalResponse = []TotalField

type WinRateResponse struct {
	Wins   int `json:"win"`
	Losses int `json:"lose"`
}
