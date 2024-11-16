package database

import "github.com/melkzsiqueira/water-gas-measurement/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(id string) (*entity.User, error)
}
