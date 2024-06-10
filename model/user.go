package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

type Room struct {
	RoomID int    `json:"room_id,omitempty"`
	Message string `json:"message"`
}
