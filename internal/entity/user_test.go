package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("John Doe", "j@j.com", "123456")

	assert.Nil(t, err)
	assert.NotNil(t, u)
	assert.NotEmpty(t, u.ID)
	assert.NotEmpty(t, u.Password)
	assert.Equal(t, "John Doe", u.Name)
	assert.Equal(t, "j@j.com", u.Email)
}

func TestUserValidatePassword(t *testing.T) {
	u, err := NewUser("John Doe", "j@j.com", "123456")

	assert.Nil(t, err)
	assert.True(t, u.ValidatePassword("123456"))
	assert.False(t, u.ValidatePassword("654321"))
	assert.NotEqual(t, "123456", u.Password)
}
