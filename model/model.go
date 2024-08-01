package model

type User struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Lang            string          `json:"lang"`
	Status          string          `json:"status"`
	PhotoScoreIndex int             `json:"photo_score_index"`
	Score           map[int]float32 `json:"score"`
	Photo           map[int]string  `json:"photo"`
}

type RoomUsers struct {
	GameMasterID string          `json:"game_master_id"`
	TotalPlayers int             `json:"total_players"`
	Users        map[string]User `json:"users"`
}

type Room struct {
	GameMasterID     string          `json:"game_master_id"`
	TotalPlayers     int             `json:"total_players"`
	GameRounds       int             `json:"game_rounds"`
	CurrentRound     int             `json:"current_round"`
	GameStatus       string          `json:"game_status"`
	PhotoDescription []string        `json:"photo_description"`
	Users            map[string]User `json:"users"`
}

type GameStatus struct {
	GameStatus   string `json:"game_status"`
	CurrentRound int    `json:"current_round"`
}

const (
	PreGame         = "pre-game"
	GameMasterPhoto = "game-master-photo"
	PlayerPhoto     = "player-photo"
	Result          = "result"
)
