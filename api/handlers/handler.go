package handlers

import (
	"database/sql"
	//"encoding/json"
	"net/http"
	//"splitwise/models"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
}
