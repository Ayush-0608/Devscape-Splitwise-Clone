package handlers

import (

	//"encoding/json"
	"database/sql"
	"fmt"
	"net/http"
	"splitwise/api/auth"
	"splitwise/models"
	"splitwise/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	Store models.UserStore
}

var secret = []byte("secret_key")

func NewHandler(s models.UserStore) *Handler {
	return &Handler{Store: s}
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
	_, err := h.Store.GetUserByMail(user.Email)
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

	err = h.Store.CreateUser(models.RegUser{
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
	u, err := h.Store.GetUserByMail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	//check password
	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	token, err := auth.CreateJWT(secret, u.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) CheckUsers(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	users, err := h.Store.GetUsers(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) CheckProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := h.Store.GetUserByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("%d", id))
		return
	}

	if userID != id {
		var ans bool
		ans, err = h.Store.CheckFriend(userID, id)

		if err != nil || !ans {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("user not valid or friend"))
			return
		}
	}

	user, err := h.Store.GetUserByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	profile := models.Profile{
		ID:        user.ID,
		Email:     user.Email,
		Fullname:  user.Fullname,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, profile)
}

func (h *Handler) SetProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var profile models.Profile
	if err := utils.ParseJSON(r, &profile); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.Store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID != userID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	h.Store.UpdateProfile(profile)

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) Friends(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	friends, err := h.Store.GetFriends(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, friends)
}

func (h *Handler) Requests(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	friends, err := h.Store.GetRequests(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, friends)
}

func (h *Handler) SendRequest(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if userID == id {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid relation"))
		return
	}

	ans, err := h.Store.CheckFriend(userID, id)

	if ans {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("relation already exists"))
		return
	}

	if err != sql.ErrNoRows {
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("relation already exists"))
		return

	}

	_, err = h.Store.CheckFriend(id, userID)

	if err == sql.ErrNoRows {
		err = h.Store.AddFriend(userID, id, false)
	} else if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	} else {
		err = h.Store.AddFriend(userID, id, true)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		err = h.Store.AcceptFriend(id, userID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) Groups(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	groups, err := h.Store.GetGroups(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, groups)
}

func (h *Handler) MakeGroup(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var group models.Group
	if err := utils.ParseJSON(r, &group); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(group); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.Store.AddGroup(&group, userID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Store.AddMember(group.ID, userID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) GroupInfo(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	ans, err := h.Store.CheckMember(userID, id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	details, err := h.Store.GetGroupDetails(id)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, details)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	ans, err := h.Store.CheckMember(userID, id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	var other models.SimpleID
	if err := utils.ParseJSON(r, &other); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ans, err = h.Store.CheckFriend(userID, other.ID)

	if !ans {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not friend"))
		return
	}

	ans, err = h.Store.CheckMember(other.ID, id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if ans {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("already exists"))
		return
	}

	err = h.Store.AddMember(id, other.ID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) LogExpense(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var expense models.NewExpense
	if err := utils.ParseJSON(r, &expense); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(expense); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	ans, err := h.Store.CheckMember(userID, expense.GroupID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	err = h.Store.AddExpense(expense)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) GroupExpenses(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var group models.SimpleID
	if err := utils.ParseJSON(r, &group); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ans, err := h.Store.CheckMember(userID, group.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	expenses, err := h.Store.GetExpenses(group.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, expenses)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	groupID, err := h.Store.GetExpenseGroup(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ans, err := h.Store.CheckMember(userID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	err = h.Store.RemoveExpense(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)

}

func (h *Handler) Split(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	groupID, err := h.Store.GetExpenseGroup(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ans, err := h.Store.CheckMember(userID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	var payload struct {
		SplitType string         `json:"split_type"`
		Splits    []models.Split `json:"splits"`
	}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusEarlyHints, err)
		return
	}

	if payload.SplitType != "Equal" && payload.SplitType != "Amount" && payload.SplitType != "Percentage" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid split type"))
		return
	}

	sum, count, err := h.Store.CheckValidity(payload.Splits, id, payload.SplitType)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	member_count, err := h.Store.MemberCount(groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.Store.AddSplits(payload.Splits, payload.SplitType, id, sum, member_count, count)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) Balances(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	users, err := h.Store.GetUsers(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	balances, err := h.Store.GetBalances(userID, users)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, *balances)
}

func (h *Handler) GroupBalances(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	groupID, _ := strconv.Atoi(params["id"])

	ans, err := h.Store.CheckMember(userID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	balances, err := h.Store.GetGroupBalances(userID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusFound, *balances)
}

func (h *Handler) Settle(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	groupID, err := h.Store.GetExpenseGroup(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ans, err := h.Store.CheckMember(userID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if !ans {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("unauthorized operation"))
		return
	}

	err = h.Store.MarkPaid(userID, id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
