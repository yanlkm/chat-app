package message

import (
	"context"
)

// MessageService defines the methods that a message service should implement
type MessageService interface {
	CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error)
	GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error)
	GetMessage(ctx context.Context, messageID string) (*MessageEntity, error)
	DeleteMessage(ctx context.Context, messageID string) error
}

// messageService is a struct that embeds the message repository
type messageService struct {
	repo MessageRepository
}

// NewMessageService creates a new instance of MessageService
func NewMessageService(repo MessageRepository) MessageService {
	return &messageService{repo: repo}
}

// CreateMessage creates a new message
func (m *messageService) CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error) {
	return m.repo.CreateMessage(ctx, message)
}

// GetMessages retrieves all messages from a room
func (m *messageService) GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error) {
	return m.repo.GetMessages(ctx, roomID)
}

// GetMessage retrieves a message by its ID
func (m *messageService) GetMessage(ctx context.Context, messageID string) (*MessageEntity, error) {
	return m.repo.GetMessage(ctx, messageID)
}

// DeleteMessage deletes a message by its ID
func (m *messageService) DeleteMessage(ctx context.Context, messageID string) error {
	return m.repo.DeleteMessage(ctx, messageID)
}
