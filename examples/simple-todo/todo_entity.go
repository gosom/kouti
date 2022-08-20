package main

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"id"`
	Content     string    `json:"content"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
