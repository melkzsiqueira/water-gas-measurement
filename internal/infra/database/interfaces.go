package database

import "github.com/melkzsiqueira/water-gas-measurement/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(id string) (*entity.User, error)
}

type MeasurementInterface interface {
	Create(measurement *entity.Measurement) (*entity.Measurement, error)
	FindAll(page, limit int, sort string) ([]entity.Measurement, error)
	FindById(id string) (*entity.Measurement, error)
	Update(measurement *entity.Measurement) error
	Delete(id string) error
}
