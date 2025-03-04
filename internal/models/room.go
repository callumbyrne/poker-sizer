package models

import "time"

type RoomState string

const (
	Voting   RoomState = "voting"
	Revealed RoomState = "revealed"
	Reset    RoomState = "reset"
)

type Room struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	Users       map[string]*User
	Votes       map[string]*Vote
	State       RoomState
	ActiveIssue string
}
