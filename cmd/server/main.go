package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/melkzsiqueira/water-gas-measurement/configs"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/webserver/handlers"
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
	measurementHandler := handlers.NewMeasurementHandler(measurementDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/v1/measurements", measurementHandler.CreateMeasurement)
	r.Get("/v1/measurements", measurementHandler.GetMeasurements)
	r.Get("/v1/measurements/{id}", measurementHandler.GetMeasurement)
	r.Put("/v1/measurements/{id}", measurementHandler.UpdateMeasurement)
	r.Delete("/v1/measurements/{id}", measurementHandler.DeleteMeasurement)

	err = http.ListenAndServe(":"+config.WebServerPort, r)
	if err != nil {
		panic(err)
	}
}
