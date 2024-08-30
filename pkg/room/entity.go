package room

// RoomEntity represents the structure of a room entity.
type RoomEntity struct {
	ID          string   `json:"_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Creator     string   `json:"creator,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
	Members     []string `json:"members,omitempty"`
	Hashtags    []string `json:"hashtags,omitempty"`
	Messages    []string `json:"messages,omitempty"`
}

// Member of a room
type MemberEntity struct {
	ID string `json:"ID,omitempty" `
}

// Hashtag of a room
type HashtagEntity struct {
	Hashtag string `json:"hashtag,omitempty" `
}
