package handlers

import (

	//"encoding/json"
	"fmt"
	"net/http"
	"splitwise/api/auth"
	"splitwise/models"
	"splitwise/utils"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(s models.UserStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	//get JSON
	var user models.RegUser
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//check user
	_, err := h.store.GetUserByMail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
		return
	}

	//create user
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(models.RegUser{
		Fullname: user.Fullname,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	//get JSON
	var user models.LoginUser
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate payload
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//check user
	u, err := h.store.GetUserByMail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	//check password
	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte("secret_key")
	token, err := auth.CreateJWT(secret, u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) CheckUsers(w http.ResponseWriter, r *http.Request) {
	u, err := h.store.GetUsers()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, u)
}
