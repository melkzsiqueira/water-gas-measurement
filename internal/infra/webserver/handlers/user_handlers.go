package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/goccy/go-json"
	"github.com/melkzsiqueira/water-gas-measurement/internal/dto"
	"github.com/melkzsiqueira/water-gas-measurement/internal/entity"
	"github.com/melkzsiqueira/water-gas-measurement/internal/infra/database"
)

type UserHandler struct {
	UserDB database.UserInterface
}

type Error struct {
	Message string `json:"message"`
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

// Get token 	godoc
// @Summary    	Get a user token
// @Description	Get a user token
// @Tags       	users
// @Accept     	json
// @Produce    	json
// @Param      	request				body		dto.GetTokenInput	true	"user credentials"
// @Success    	200					{object}	dto.GetTokenOutput
// @Failure    	400					{object}	Error
// @Failure    	401					{object}  	Error
// @Router     	/users/token		[post]
func (h *UserHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("token").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("token_expires_in").(int)

	var login dto.GetTokenInput
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	u, err := h.UserDB.FindByEmail(login.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		error := Error{Message: "email or password invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}
	if !u.ValidatePassword(login.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := Error{Message: "email or password invalid"}
		json.NewEncoder(w).Encode(error)
		return
	}
	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	response := &dto.GetTokenOutput{AccessToken: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Create user 	godoc
// @Summary     Create user
// @Description	Create user
// @Tags		users
// @Accept      json
// @Produce     json
// @Param       request	body      dto.CreateUserInput	true	"user request"
// @Success     201		{object}  entity.User
// @Failure     400		{object}  Error
// @Failure     500		{object}  Error
// @Router      /users	[post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
