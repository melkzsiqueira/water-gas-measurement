package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/melkzsiqueira/water-gas-measurement/configs"
	_ "github.com/melkzsiqueira/water-gas-measurement/docs"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Water and Gas Measurement
// @version         1.0
// @description     Water and Gas Measurement API with auhtentication
// @termsOfService	http://swagger.io/terms/

// @contact.name	Melkz Siqueira
// @contact.url		https://linkedin.com/in/melkzsiqueira
// @contact.email	melkz.siqueira@gmail.com

// @license.name	Apache-2.0 license
// @license.url		https://github.com/melkzsiqueira/water-gas-measurement?tab=Apache-2.0-1-ov-file#

// @securityDefinitions.apikey	ApiKeyAuth
// @in 							header
// @name 						Authorization
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
	r.Route("/"+config.APIVersion, func(r chi.Router) {
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
		r.Route("/docs", func(r chi.Router) {
			r.Get("/*", httpSwagger.Handler(httpSwagger.URL(config.SwaggerURL)))
		})
	})

	err = http.ListenAndServe(":"+config.WebServerPort, r)
	if err != nil {
		panic(err)
	}
}
