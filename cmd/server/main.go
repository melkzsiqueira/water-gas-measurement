package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
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

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("token", config.TokenAuth))
	r.Use(middleware.WithValue("token_expires_in", config.JWTExpiresIn))
	r.Route("/v1", func(r chi.Router) {
		r.Route("/measurements", func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/", measurementHandler.CreateMeasurement)
			r.Get("/", measurementHandler.GetMeasurements)
			r.Get("/{id}", measurementHandler.GetMeasurement)
			r.Put("/{id}", measurementHandler.UpdateMeasurement)
			r.Delete("/{id}", measurementHandler.DeleteMeasurement)
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Post("/token", userHandler.GetToken)
		})
	})

	err = http.ListenAndServe(":"+config.WebServerPort, r)
	if err != nil {
		panic(err)
	}
}
