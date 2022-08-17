package users

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosom/kouti/um"
	"github.com/gosom/kouti/web"
	"github.com/rs/zerolog"
)

type UserQueryParams struct {
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Roles     []um.Role `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserCreate struct {
	Email    string `json:"email" validate:"required,email" example:"aris.paparis@example.com"`
	Password string `json:"password" validate:"required,password" example:"Ar9Sp7891!!#"`
}

func NewUserHandler(log zerolog.Logger, srv *UserSrv) web.ResourceHandler[UserCreate, UserQueryParams, User] {
	h := web.ResourceHandler[UserCreate, UserQueryParams, User]{}
	h.Logger = log
	h.Srv = srv
	return h
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email" example:"aris.paparis@example.com"`
	Password string `json:"password" validate:"required" example:"Ar9Sp7891!!#"`
}

type L struct {
	AccessToken string `json:"accessToken"`
}

func NewAuthHandler(log zerolog.Logger, srv *UserSrv) web.AuthHandler[UserLogin, L] {
	h := web.AuthHandler[UserLogin, L]{}
	h.Logger = log
	h.Srv = srv
	return h
}
