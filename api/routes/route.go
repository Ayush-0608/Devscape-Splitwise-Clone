package routes

import (
	"splitwise/api/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(h *handlers.Handler) *mux.Router {
	//NewRouter this returns an instance of a router
	r := mux.NewRouter()
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/users", h.CheckUsers).Methods("GET")
	return r
}
