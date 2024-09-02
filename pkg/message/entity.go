package message

// MessageEntity struct
type MessageEntity struct {
	ID        string `json:"_id,omitempty"`
	RoomID    string `json:"roomId,omitempty"`
	Username  string `json:"username,omitempty"`
	UserID    string `json:"userId,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
