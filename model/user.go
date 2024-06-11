package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

type Room struct {
	GameMaster   User   `json:"game_master"`
	Players      []User `json:"players"`
	TotalPlayers int    `json:"total_players"`
}
