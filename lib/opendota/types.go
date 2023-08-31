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

type Match struct {
	Id           int  `json:"match_id"`
	Slot         int  `json:"player_slot"`
	RadiantWin   bool `json:"radiant_win"`
	Duration     int  `json:"duration"`
	GameMode     int  `json:"game_mode"`
	LobbyType    int  `json:"lobby_type"`
	HeroId       int  `json:"hero_id"`
	StartTime    int  `json:"start_time"`
	Version      *int `json:"version"`
	Kills        int  `json:"kills"`
	Deaths       int  `json:"deaths"`
	Assists      int  `json:"assists"`
	SkillBracket *int `json:"skill"`
	RankAverage  *int `json:"average_rank"`
	LeaverStatus int  `json:"leaver_status"`
	PartySize    int  `json:"party_size"`
}
