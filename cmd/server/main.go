package main

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/melkzsiqueira/water-gas-measurement/configs"
	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.DBDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Measurement{}, &entity.User{})

	measurementDB := database.NewMeasurement(db)
	measurementHandler := NewMeasurementHandler(measurementDB)

	http.HandleFunc("/v1/measurements", measurementHandler.CreateMeasurement)
	err = http.ListenAndServe(":"+config.WebServerPort, nil)
	if err != nil {
		panic(err)
	}
}

type MeasurementHandler struct {
	MeasurementDB database.MeasurementInterface
}

func NewMeasurementHandler(db database.MeasurementInterface) *MeasurementHandler {
	return &MeasurementHandler{
		MeasurementDB: db,
	}
}

func (h *MeasurementHandler) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	var measurement dto.CreateMeasurementInput
	err := json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, err := entity.NewMeasurement(measurement.Value, measurement.Image, measurement.Type, measurement.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdMeasurement, err := h.MeasurementDB.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(createdMeasurement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
