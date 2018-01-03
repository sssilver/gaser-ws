package main

import (
	"time"

	"github.com/google/uuid"
)

type Frame struct {
	eventNum int
	gameID   uuid.UUID `json:"game_id",omitempty`
	time     time.Time

	data []byte `json:omitempty`
}

type OutFrame = Frame
type InFrame = Frame
