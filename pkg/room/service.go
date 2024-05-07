package room

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomService interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	CheckName(ctx context.Context, name string) error
	GetRoom(ctx context.Context, roomID primitive.ObjectID) (*Room, error)
	GetAllRooms(ctx context.Context) ([]Room, error)
	AddMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error)
	RemoveMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error)
	DeleteRoom(ctx context.Context, roomID primitive.ObjectID) error
}

type roomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) RoomService {
	return &roomService{repo: repo}

}

func (r *roomService) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	return r.repo.CreateRoom(ctx, room)
}

func (r *roomService) GetRoom(ctx context.Context, roomID primitive.ObjectID) (*Room, error) {
	return r.repo.GetRoom(ctx, roomID)
}

func (r *roomService) GetAllRooms(ctx context.Context) ([]Room, error) {
	return r.repo.GetAllRooms(ctx)
}
func (r *roomService) CheckName(ctx context.Context, name string) error {
	return r.repo.CheckName(ctx, name)
}

func (r *roomService) RemoveMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error) {
	return r.repo.RemoveMember(ctx, roomID, memberID)
}

func (r *roomService) AddMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error) {
	return r.repo.AddMember(ctx, roomID, memberID)
}

func (r *roomService) DeleteRoom(ctx context.Context, roomID primitive.ObjectID) error {
	return r.repo.Delete(ctx, roomID)
}
