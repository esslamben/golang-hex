package routes

import (
	"encoding/json"
	"net/http"

	"github.com/esslamb/golang-hex/pkg/user"
	"github.com/esslamb/golang-hex/pkg/utils"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// IHandler lists the methods the handler package implements
type IHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserByUUID(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	UserService user.IService
}

// NewHandler returns new instance of handler with all methods
// attached
func NewHandler(u user.IService) IHandler {
	return &handler{
		u,
	}
}

// CreateUser unmarshal's request data into User type and
// passes along the information to pkg handler
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	b := &user.User{}
	utils.ParseBody(r, b)

	if ok, errors := utils.ValidateInputs(b); !ok {
		r, err := utils.CreateValidationResponse(errors)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(r)
		return
	}

	u, err := h.UserService.CreateUser(b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	res, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
}

// GetUserByUUID validates UUID passed and calls service function
// to find UUID.
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid, err := uuid.FromString(vars["userUUID"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	u, err := h.UserService.ReadUser(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := uuid.FromString(vars["userUUID"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	b := &user.User{}
	utils.ParseBody(r, b)

	if ok, errors := utils.ValidateInputs(b); !ok {
		r, err := utils.CreateValidationResponse(errors)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(r)
		return
	}

	err = h.UserService.UpdateUser(b, u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Updated"))
	return
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u, err := uuid.FromString(vars["userUUID"])
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Invalid UUID"))
	}

	err = h.UserService.DeleteUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Deleted"))
}
