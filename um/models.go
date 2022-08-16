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

func (o *User) HasRole(role string) bool {
	for i := range o.Roles {
		if o.Roles[i].Name == role {
			return true
		}
	}
	return false
}

func (o *User) HasAnyRole(roles ...string) bool {
	for i := range roles {
		if o.HasRole(roles[i]) {
			return true
		}
	}
	return false
}

// Role the roles of the system
type Role struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
