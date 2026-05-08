package routes

import (
	"splitwise/api/auth"
	"splitwise/api/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes(h *handlers.Handler) *mux.Router {
	//NewRouter this returns an instance of a router
	r := mux.NewRouter()
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/users", auth.WithJWTAuth(h.CheckUsers, h.Store)).Methods("GET")
	r.HandleFunc("/profile/{id}", auth.WithJWTAuth(h.CheckProfile, h.Store)).Methods("GET")
	r.HandleFunc("/profile", auth.WithJWTAuth(h.SetProfile, h.Store)).Methods("POST")
	r.HandleFunc("/friends", auth.WithJWTAuth(h.Friends, h.Store)).Methods("GET")
	r.HandleFunc("/friends/requests", auth.WithJWTAuth(h.Requests, h.Store)).Methods("GET")
	r.HandleFunc("/friends/requests/{id}", auth.WithJWTAuth(h.SendRequest, h.Store)).Methods("POST")
	r.HandleFunc("/groups", auth.WithJWTAuth(h.Groups, h.Store)).Methods("GET")
	r.HandleFunc("/groups", auth.WithJWTAuth(h.MakeGroup, h.Store)).Methods("POST")
	r.HandleFunc("/groups/{id}", auth.WithJWTAuth(h.GroupInfo, h.Store)).Methods("GET")
	r.HandleFunc("/groups/{id}", auth.WithJWTAuth(h.AddMember, h.Store)).Methods("POST")
	r.HandleFunc("/expenses", auth.WithJWTAuth(h.LogExpense, h.Store)).Methods("POST")
	r.HandleFunc("/expenses", auth.WithJWTAuth(h.GroupExpenses, h.Store)).Methods("GET")
	r.HandleFunc("/expenses/{id}", auth.WithJWTAuth(h.DeleteExpense, h.Store)).Methods("DELETE")
	r.HandleFunc("/expenses/{id}/split", auth.WithJWTAuth(h.Split, h.Store)).Methods("POST")
	r.HandleFunc("/balance", auth.WithJWTAuth(h.Balances, h.Store)).Methods("GET")
	r.HandleFunc("/groups/{id}/balance", auth.WithJWTAuth(h.GroupBalances, h.Store)).Methods("GET")
	r.HandleFunc("/expenses/{id}", auth.WithJWTAuth(h.Settle, h.Store)).Methods("PUT")
	return r
}
