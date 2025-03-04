package models

import "time"

type Vote struct {
	UserID    string
	RoomID    string
	Value     string
	CreatedAt time.Time
}
