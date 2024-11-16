package entity

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/melkzsiqueira/water-gas-measurement/pkg/entity"
)

type Measurement struct {
	ID        entity.ID `json:"id"`
	Value     int       `json:"value"`
	Image     string    `json:"image"`
	Type      string    `json:"type"`
	Confirmed bool      `json:"confirmed"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrValueIsRequired = errors.New("value is required")
	ErrInvalidValue    = errors.New("invalid value")
	ErrImageIsRequired = errors.New("image is required")
	ErrInvalidImage    = errors.New("invalid image")
	ErrTypeIsRequired  = errors.New("type is required")
	ErrInvalidType     = errors.New("invalid type")
	ErrUserIsRequired  = errors.New("user is required")
	ErrInvalidUser     = errors.New("invalid user")
)

func NewMeasurement(value int, image string, measurementType string, user string) (*Measurement, error) {
	measurement := &Measurement{
		ID:        entity.NewID(),
		Value:     value,
		Image:     image,
		Type:      measurementType,
		Confirmed: false,
		User:      user,
		CreatedAt: time.Now(),
	}

	err := measurement.Validate()

	if err != nil {
		return nil, err
	}

	return measurement, nil
}

func (m *Measurement) Validate() error {
	if m.Value == 0 {
		return ErrValueIsRequired
	}

	if m.Value < 0 {
		return ErrInvalidValue
	}

	if m.Image == "" {
		return ErrImageIsRequired
	}

	if _, err := base64.StdEncoding.DecodeString(m.Image); err != nil {
		return ErrInvalidImage
	}

	if m.Type == "" {
		return ErrTypeIsRequired
	}

	if m.Type != "1" && m.Type != "2" {
		return ErrInvalidType
	}

	if m.User == "" {
		return ErrUserIsRequired
	}

	if _, err := entity.ParseID(m.User); err != nil {
		return ErrInvalidUser
	}

	return nil
}
