package models

import (
	"time"
)

type UserStore interface {
	GetUserByMail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user RegUser) error
	GetUsers(userID int) (*[]PublicUser, error)
	UpdateProfile(profile Profile) error
	CheckFriend(userID int, ID int) (bool, error)
	GetFriends(userID int) (*[]Friend, error)
	GetRequests(userID int) (*[]Friend, error)
	AddFriend(id1 int, id2 int, ans bool) error
	AcceptFriend(id1 int, id2 int) error
	GetGroups(userID int) (*[]Group, error)
	AddGroup(group *Group, userID int) error
	AddMember(groupID int, id int) error
	CheckMember(userID int, id int) (bool, error)
	GetGroupDetails(id int) (*GroupDetails, error)
	AddExpense(e NewExpense) error
	GetExpenses(groupID int) (*[]Expense, error)
	GetExpenseGroup(id int) (int, error)
	RemoveExpense(id int) error
	CheckValidity(splits []Split, id int, splitType string) (int, int, error)
	MemberCount(groupID int) (int, error)
	AddSplits(splits []Split, splitType string, id int, sum int, member_count int, count int) error
	GetBalances(userID int, set *[]PublicUser) (*Balance, error)
	GetGroupBalances(userID int, groupID int) (*Balance, error)
	MarkPaid(userID int, id int) error
}

type RegUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Profile struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type PublicUser struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type Friend struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Result   bool   `json:"result"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type GroupDetails struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	CreatedBy int          `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	Members   []PublicUser `json:"members"`
}

type SimpleID struct {
	ID int `json:"id"`
}

type Expense struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	PaidBy      int       `json:"paid_by"`
	SplitType   string    `json:"split_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type NewExpense struct {
	GroupID     int    `json:"group_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Amount      int    `json:"amount" validate:"required,gt=0"`
	PaidBy      int    `json:"paid_by" validate:"required"`
	SplitType   string `json:"split_type"`
}

type Split struct {
	UserID int `json:"user_id"`
	Value  int `json:"value"`
}

type Slot struct {
	SplitID   int    `json:"split_id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	ExpenseID int    `json:"expense_id"`
	Value     int    `json:"value"`
}

type Balance struct {
	OwedTo []Slot `json:"owedto"`
	OwedBy []Slot `json:"owedby"`
}
