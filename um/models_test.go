package um

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserHasRoles(t *testing.T) {
	u := User{
		ID:        1,
		UID:       uuid.New(),
		Identity:  "username",
		EncPasswd: "",
		Roles: []Role{
			Role{ID: 1, Name: "admin"},
			Role{ID: 2, Name: "registered"},
		},
	}

	isAdmin := u.HasAnyRole("admin", "registered")
	if !isAdmin {
		t.Error("expected to find admin role")
		return
	}

	isMember := u.HasRole("registered")
	if !isMember {
		t.Error("expected to find registered role")
		return
	}

	isSuper := u.HasRole("super")
	if isSuper {
		t.Error("expected not to find super role")
		return
	}
}
