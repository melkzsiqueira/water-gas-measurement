package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
	entityPkg "github.com/melkzsiqueira/water-gas-measurement/pkg/entity"
)

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &measurement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := entity.NewMeasurement(measurement.Value, measurement.Image, measurement.Type, measurement.User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	createdMeasurement, err := h.MeasurementDB.Create(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMeasurement)
}

func (h *MeasurementHandler) GetMeasurements(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" || page == "0" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("page is invalid"))
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" || limit == "0" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("limit is invalid"))
		return
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "desc"
	}
	if sort != "asc" && sort != "desc" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("sort is invalid"))
		return
	}

	measurements, err := h.MeasurementDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(measurements)
}

func (h *MeasurementHandler) GetMeasurement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is required"))
		return
	}
	measurement, err := h.MeasurementDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(measurement)
}

func (h *MeasurementHandler) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is required"))
		return
	}
	var measurement entity.Measurement
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &measurement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	measurement.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is invalid"))
		return
	}
	_, err = h.MeasurementDB.FindById(measurement.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.MeasurementDB.Update(&measurement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(measurement)
}

func (h *MeasurementHandler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is required"))
		return
	}
	_, err := h.MeasurementDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.MeasurementDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
