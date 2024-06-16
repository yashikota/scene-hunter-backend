package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
	Status string `json:"status"`
	PhotoScoreIndex int `json:"photo_score_index"`
	Score map[int]float32 `json:"score"`
	Photo map[int]string `json:"photo"`
}

type Room struct {
	GameMasterID string `json:"game_master_id"`
	TotalPlayers int    `json:"total_players"`
	GameRounds   int    `json:"game_rounds"`
	Users        map[string]User `json:"users"`
}
