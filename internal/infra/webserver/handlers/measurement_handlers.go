package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/gemini"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/storage"
	entityPkg "github.com/melkzsiqueira/water-gas-measurement/pkg/entity"
)

type MeasurementHandler struct {
	MeasurementDB      database.MeasurementInterface
	MeasurementStorage storage.MeasurementStorageInterface
	Gemini             gemini.GeminiInterface
}

func NewMeasurementHandler(db database.MeasurementInterface, storage storage.MeasurementStorageInterface, gemini gemini.GeminiInterface) *MeasurementHandler {
	return &MeasurementHandler{
		MeasurementDB:      db,
		MeasurementStorage: storage,
		Gemini:             gemini,
	}
}

// Create measurement	godoc
// @Summary      		Create measurement
// @Description  		Create measurement
// @Tags         		measurements
// @Accept       		json
// @Produce      		json
// @Param        		request				body		dto.CreateMeasurementInput	true	"measurement request"
// @Success      		201					{object}	entity.Measurement
// @Failure      		400         		{object}	Error
// @Failure      		500         		{object}	Error
// @Router       		/measurements		[post]
// @Security 			ApiKeyAuth
func (h *MeasurementHandler) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	var measurement dto.CreateMeasurementInput

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = json.Unmarshal(body, &measurement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	imgReq := dto.ProcessImageRequest{
		Data: measurement.Image.Data,
		Mime: measurement.Image.Mime,
	}
	imgResp, err := h.Gemini.ProcessImage(imgReq, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	measurement.Value, err = strconv.Atoi(imgResp.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	s, err := h.MeasurementStorage.UploadFile("data:"+measurement.Image.Mime+";base64,"+measurement.Image.Data, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	m, err := entity.NewMeasurement(
		measurement.Value,
		s.SecureURL,
		measurement.Type,
		measurement.User,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.MeasurementDB.Create(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

// List measurements	godoc
// @Summary      		List measurements
// @Description  		Get all measurements
// @Tags         		measurements
// @Accept       		json
// @Produce      		json
// @Param        		page      		query   	string  			false  "page number"
// @Param        		limit     		query   	string  			false  "records limit"
// @Success      		200       		{array} 	entity.Measurement
// @Failure      		400       		{object}	Error
// @Failure      		500       		{object}	Error
// @Router       		/measurements 	[get]
// @Security 			ApiKeyAuth
func (h *MeasurementHandler) GetMeasurements(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" || page == "0" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "page is invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" || limit == "0" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "limit is invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "desc"
	}
	if sort != "asc" && sort != "desc" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "sort is invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}

	m, err := h.MeasurementDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

// Get measurement	godoc
// @Summary      	Get a measurement
// @Description  	Get a measurement
// @Tags         	measurements
// @Accept       	json
// @Produce      	json
// @Param        	id   				path		string				true	"measurement ID"	Format(uuid)
// @Success      	200  				{object}	entity.Measurement
// @Failure      	400  				{object}  	Error
// @Failure      	404  				{object}  	Error
// @Router       	/measurements/{id}	[get]
// @Security 		ApiKeyAuth
func (h *MeasurementHandler) GetMeasurement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "id is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	m, err := h.MeasurementDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

// Update measurement	godoc
// @Summary     		Update a measurement
// @Description 		Update a measurement
// @Tags        		measurements
// @Accept      		json
// @Produce     		json
// @Param       		id        			path		string                  	true	"product ID"			Format(uuid)
// @Param       		request     		body      	dto.CreateMeasurementInput	true	"measurement request"
// @Success      		200  				{object}	entity.Measurement
// @Failure     		400	   				{object}	Error
// @Failure     		404	   				{object}	Error
// @Failure     		500       			{object}	Error
// @Router      		/measurements/{id} 	[put]
// @Security 			ApiKeyAuth
func (h *MeasurementHandler) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "id is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	var m entity.Measurement
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	m.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "id is invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}
	_, err = h.MeasurementDB.FindById(m.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.MeasurementDB.Update(&m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

// Delete measurement	godoc
// @Summary      		Delete a measurement
// @Description  		Delete a measurement
// @Tags         		measurements
// @Accept       		json
// @Produce      		json
// @Param        		id        				path      	string				true	"measurement ID"	Format(uuid)
// @Success      		200						{object}	entity.Measurement
// @Failure      		400						{object}	Error
// @Failure      		404						{object}	Error
// @Failure      		500       				{object}	Error
// @Router       		/measurements/{id}		[delete]
// @Security 			ApiKeyAuth
func (h *MeasurementHandler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "id is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	_, err := h.MeasurementDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.MeasurementDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get measurement image	godoc
// @Summary      			Get a measurement image
// @Description  			Get a measurement image
// @Tags         			measurements
// @Accept       			json
// @Produce      			json
// @Param        			id   						path		string		true	"measurement ID"	Format(uuid)
// @Success      			200  						{file}  	image
// @Failure      			400  						{object}  	Error
// @Failure      			404  						{object}  	Error
// @Router       			/measurements/{id}/image	[get]
// @Security 				ApiKeyAuth
func (h *MeasurementHandler) GetMeasurementImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: "id is required"}
		json.NewEncoder(w).Encode(error)
		return
	}
	m, err := h.MeasurementDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	image, err := base64.StdEncoding.DecodeString(m.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(image)
}
