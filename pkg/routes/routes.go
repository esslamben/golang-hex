package routes

import (
	"github.com/gorilla/mux"
)

// CreateRouter creates a mux router and attaches routes
func CreateRouter(h IHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user", h.CreateUser).Methods("POST")
	r.HandleFunc("/user/{userUUID}", h.GetUserByUUID).Methods("GET")
	r.HandleFunc("/user/{userUUID}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{userUUID}", h.DeleteUser).Methods("DELETE")

	return r
}
