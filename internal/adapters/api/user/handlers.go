package user

import (
	"auth-api/internal/adapters/api"
	userDomain "auth-api/internal/domain/user"
	customError "auth-api/internal/error"
	"auth-api/internal/midlleware"
	"auth-api/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const (
	createUserURL   = "/register"
	loginUserURL    = "/login"
	userSettingsURL = "/settings"
	GET             = "GET "
	POST            = "POST "
	PUT             = "PUT "
	PATCH           = "PATCH "
	DELETE          = "DELETE "
)

type handler struct {
	userService userDomain.ServiceUser
}

func (h *handler) Register(router *http.ServeMux) {
	router.Handle(POST+createUserURL, midlleware.TimeoutMiddleware(http.HandlerFunc(h.CreateUser)))
	router.Handle(PUT+userSettingsURL, midlleware.TimeoutMiddleware(midlleware.AuthMiddleware(http.HandlerFunc(h.UpdateUser))))
	router.Handle(POST+loginUserURL, midlleware.TimeoutMiddleware(http.HandlerFunc(h.LoginUser)))
}

func NewHandler(service userDomain.ServiceUser) api.Handler {
	return &handler{userService: service}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating user..")
	var dtoUser = &userDomain.CreateUserDTO{}
	if err := json.NewDecoder(r.Body).Decode(dtoUser); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		var syntaxError *json.SyntaxError
		if errors.As(err, &unmarshalTypeError) {
			log.Println(err.Error())
			http.Error(w, "Invalid request data type", http.StatusBadRequest)
			return
		} else if errors.As(err, &syntaxError) || errors.Is(err, io.ErrUnexpectedEOF) {
			log.Println(err.Error())
			http.Error(w, "Invalid JSON syntax", http.StatusBadRequest)
			return
		} else if errors.Is(err, io.EOF) {
			log.Println(err.Error())
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
		} else {
			log.Println(err.Error())
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}
	}
	if err := h.userService.CreateUser(r.Context(), dtoUser); err != nil {
		if errors.Is(err, customError.CreateUserBadInputError) {
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		} else if errors.Is(err, customError.BusyUpdateEmailError) {
			http.Error(w, "Email is busy", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			log.Panic(err.Error())
			return
		}
	}
	utils.RenderJSON(w, http.StatusCreated, "User has been created")
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var dtoUser = &userDomain.UpdateUserDTO{}
	if err := json.NewDecoder(r.Body).Decode(dtoUser); err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	u, err := h.userService.UpdateUser(r.Context(), dtoUser)
	if err != nil {
		if errors.Is(err, customError.NothingToUpdateError) {
			http.Error(w, "No fields have been changed", http.StatusBadRequest)
			return
		} else if errors.Is(err, customError.NotFoundError) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		} else if errors.Is(err, customError.UpdateUserBadInputError) {
			http.Error(w, "Invalid password", http.StatusBadRequest)
		}
		return
	}
	utils.RenderJSON(w, http.StatusOK, u)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var dtoUser = &userDomain.CreateUserDTO{}
	if json.NewDecoder(r.Body).Decode(dtoUser) != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token, err := h.userService.Login(r.Context(), dtoUser)
	if err != nil {
		if errors.Is(err, customError.LoginError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			log.Println(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	utils.SetCookie(w, token.Token)
	utils.RenderJSON(w, http.StatusOK, token)
}
