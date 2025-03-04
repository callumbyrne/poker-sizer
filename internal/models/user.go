package models

import "time"

type User struct {
	ID       string
	Name     string
	IsAdmin  bool
	RoomID   string
	JoinedAt time.Time
}
