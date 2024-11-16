package entity

import (
	"errors"

	"github.com/melkzsiqueira/water-gas-measurement/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

var (
	ErrNameIsRequired     = errors.New("name is required")
	ErrEmailIsRequired    = errors.New("email is required")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPasswordIsRequired = errors.New("password is required")
	ErrInvalidPassword    = errors.New("invalid password")
)

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	err = user.Validate()

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameIsRequired
	}

	if u.Email == "" {
		return ErrEmailIsRequired
	}

	if u.Password == "" {
		return ErrPasswordIsRequired
	}

	return nil
}
