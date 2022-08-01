package um

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID        int       `json:"-"`
	UID       uuid.UUID `json:"uid"`
	Identity  string    `json:"identity"`
	EncPasswd string    `json:"-"`
	Roles     []Role    `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Role the roles of the system
type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
