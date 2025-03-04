package services

import (
	"errors"
	"time"

	"github.com/callumbyrne/poker-sizer/internal/models"
	"github.com/callumbyrne/poker-sizer/internal/store"
	"github.com/google/uuid"
)

type RoomService struct {
	store *store.MemoryStore
}

func NewRoomService(store *store.MemoryStore) *RoomService {
	return &RoomService{
		store: store,
	}
}

func (s *RoomService) CreateRoom(name string) (*models.Room, error) {
	if name == "" {
		return nil, errors.New("room name cannot be empty")
	}

	roomID := uuid.New().String()
	room := &models.Room{
		ID:        roomID,
		Name:      name,
		CreatedAt: time.Now(),
		Users:     make(map[string]*models.User),
		Votes:     make(map[string]*models.Vote),
		State:     models.Voting,
	}

	err := s.store.SaveRoom(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) GetRoom(id string) (*models.Room, error) {
	return s.store.GetRoom(id)
}

func (s *RoomService) AddUserToRoom(roomID, userName string) (*models.User, error) {
	room, err := s.store.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	userID := uuid.New().String()
	isAdmin := len(room.Users) == 0 // First user is admin

	user := &models.User{
		ID:       userID,
		Name:     userName,
		IsAdmin:  isAdmin,
		RoomID:   roomID,
		JoinedAt: time.Now(),
	}

	room.Users[userID] = user
	err = s.store.SaveRoom(room)
	if err != nil {
		return nil, err
	}

	return user, nil
}
