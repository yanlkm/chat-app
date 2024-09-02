package message

import (
	"context"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error)
	GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error)
	GetMessage(ctx context.Context, messageID string) (*MessageEntity, error)
	DeleteMessage(ctx context.Context, messageID string) error
}

type messageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (m *messageService) CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error) {
	return m.repo.CreateMessage(ctx, message)
}

func (m *messageService) GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error) {
	return m.repo.GetMessages(ctx, roomID)
}
func (m *messageService) GetMessage(ctx context.Context, messageID string) (*MessageEntity, error) {
	return m.repo.GetMessage(ctx, messageID)
}

func (m *messageService) DeleteMessage(ctx context.Context, messageID string) error {
	return m.repo.DeleteMessage(ctx, messageID)
}
