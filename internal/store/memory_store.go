package store

import (
	"errors"
	"sync"

	"github.com/callumbyrne/poker-sizer/internal/models"
)

type MemoryStore struct {
	rooms map[string]*models.Room
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rooms: make(map[string]*models.Room),
	}
}

func (s *MemoryStore) SaveRoom(room *models.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rooms[room.ID] = room
	return nil
}

func (s *MemoryStore) GetRoom(id string) (*models.Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (s *MemoryStore) DeleteRoom(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.rooms, id)
	return nil
}
