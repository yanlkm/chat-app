package message

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *Message) (*Message, error)
	GetMessages(ctx context.Context, roomID primitive.ObjectID) ([]*Message, error)
	GetMessage(ctx context.Context, messageID primitive.ObjectID) (*Message, error)
	DeleteMessage(ctx context.Context, messageID primitive.ObjectID) error
}

type messageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (m *messageService) CreateMessage(ctx context.Context, message *Message) (*Message, error) {
	return m.repo.CreateMessage(ctx, message)
}

func (m *messageService) GetMessages(ctx context.Context, roomID primitive.ObjectID) ([]*Message, error) {
	return m.repo.GetMessages(ctx, roomID)
}
func (m *messageService) GetMessage(ctx context.Context, messageID primitive.ObjectID) (*Message, error) {
	return m.repo.GetMessage(ctx, messageID)
}

func (m *messageService) DeleteMessage(ctx context.Context, messageID primitive.ObjectID) error {
	return m.repo.DeleteMessage(ctx, messageID)
}
