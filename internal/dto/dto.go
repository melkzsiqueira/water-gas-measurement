package dto

import (
	"time"

	"github.com/melkzsiqueira/water-gas-measurement/pkg/entity"
)

type CreateMeasurementInput struct {
	ID        entity.ID `json:"id"`
	Value     int       `json:"value"`
	Image     string    `json:"image"`
	Type      string    `json:"type"`
	Confirmed bool      `json:"confirmed"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
