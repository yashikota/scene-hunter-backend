package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
	Status string `json:"status"`
	Score []float32 `json:"score"`
	Photo []string `json:"photo"`
}

type Room struct {
	GameMasterID string `json:"game_master_id"`
	TotalPlayers int    `json:"total_players"`
	GameRounds   int    `json:"game_rounds"`
	Users        map[string]User `json:"users"`
}
