package room

import (
	"context"
)

type RoomService interface {
	CreateRoom(ctx context.Context, room *RoomEntity) (*RoomEntity, error)
	CheckName(ctx context.Context, name string) error
	GetRoom(ctx context.Context, roomID string) (*RoomEntity, error)
	GetUserRooms(ctx context.Context, userID string) ([]RoomEntity, error)
	GetAllRooms(ctx context.Context) ([]RoomEntity, error)
	GetRoomsCreatedByAdmin(ctx context.Context, adminID string) ([]RoomEntity, error)
	AddMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error)
	RemoveMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error)
	AddHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error)
	RemoveHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error)
	DeleteRoom(ctx context.Context, roomID string) error
}

type roomService struct {
	repo RoomRepository
}

func NewRoomService(repo RoomRepository) RoomService {
	return &roomService{repo: repo}

}

func (r *roomService) CreateRoom(ctx context.Context, room *RoomEntity) (*RoomEntity, error) {
	return r.repo.CreateRoom(ctx, room)
}

func (r *roomService) GetRoom(ctx context.Context, roomID string) (*RoomEntity, error) {
	return r.repo.GetRoom(ctx, roomID)
}
func (r *roomService) GetUserRooms(ctx context.Context, userID string) ([]RoomEntity, error) {
	return r.repo.GetUserRooms(ctx, userID)
}
func (r *roomService) GetAllRooms(ctx context.Context) ([]RoomEntity, error) {
	return r.repo.GetAllRooms(ctx)
}

func (r *roomService) GetRoomsCreatedByAdmin(ctx context.Context, adminID string) ([]RoomEntity, error) {
	return r.repo.GetRoomsCreatedByAdmin(ctx, adminID)
}
func (r *roomService) CheckName(ctx context.Context, name string) error {
	return r.repo.CheckName(ctx, name)
}

func (r *roomService) RemoveMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error) {
	return r.repo.RemoveMember(ctx, roomID, memberID)
}

func (r *roomService) AddMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error) {
	return r.repo.AddMember(ctx, roomID, memberID)
}

func (r *roomService) AddHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error) {
	return r.repo.AddHashtag(ctx, roomID, hashtag)
}
func (r *roomService) RemoveHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error) {
	return r.repo.RemoveHashtag(ctx, roomID, hashtag)
}
func (r *roomService) DeleteRoom(ctx context.Context, roomID string) error {
	return r.repo.Delete(ctx, roomID)
}
