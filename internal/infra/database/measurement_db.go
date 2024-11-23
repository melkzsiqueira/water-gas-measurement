package database

import (
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"gorm.io/gorm"
)

type Measurement struct {
	DB *gorm.DB
}

func NewMeasurement(db *gorm.DB) *Measurement {
	return &Measurement{
		DB: db,
	}
}

func (m *Measurement) Create(measurement *entity.Measurement) (*entity.Measurement, error) {
	err := m.DB.Create(measurement).Error
	return measurement, err
}

func (m *Measurement) FindAll(page, limit int, sort string) ([]entity.Measurement, error) {
	var measurements []entity.Measurement
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	err := m.DB.Order("created_at " + sort).Offset((page - 1) * limit).Limit(limit).Find(&measurements).Error
	return measurements, err
}

func (m *Measurement) FindById(id string) (*entity.Measurement, error) {
	var measurement entity.Measurement
	err := m.DB.First(&measurement, "id = ?", id).Error
	return &measurement, err
}

func (m *Measurement) Update(measurement *entity.Measurement) error {
	_, err := m.FindById(measurement.ID.String())
	if err != nil {
		return err
	}
	return m.DB.Save(measurement).Error
}

func (m *Measurement) Delete(id string) error {
	measurement, err := m.FindById(id)
	if err != nil {
		return err
	}
	return m.DB.Delete(measurement).Error
}
