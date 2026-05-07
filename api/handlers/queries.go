package handlers

import (
	"database/sql"
	"fmt"
	"splitwise/models"
	"splitwise/utils"
)

type Store struct {
	DB *sql.DB
}

func (s *Store) GetUserByMail(email string) (*models.User, error) {
	rows, err := s.DB.Query("SELECT id, fullname, email, phone, password, created_at FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	for rows.Next() {
		user, err = utils.ScanRowToUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*models.User, error) {
	rows, err := s.DB.Query("SELECT id, fullname, email, phone, password, created_at FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	for rows.Next() {
		user, err = utils.ScanRowToUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user models.RegUser) error {
	_, err := s.DB.Exec("INSERT INTO users (fullname, email, phone, password) VALUES($1, $2, $3, $4)", user.Fullname, user.Email, user.Phone, user.Password)
	return err
}

func (s *Store) GetUsers(userID int) (*[]models.PublicUser, error) {
	var users []models.PublicUser
	rows, err := s.DB.Query("SELECT id, fullname, email FROM users WHERE id!=$1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user *models.PublicUser
		user, err = utils.ScanRowToPublicUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return &users, nil
}

func (s *Store) UpdateProfile(profile models.Profile) error {
	_, err := s.DB.Exec("UPDATE users SET fullname=$1, phone=$2 WHERE id=$3", profile.Fullname, profile.Phone, profile.ID)
	return err
}

func (s *Store) CheckFriend(userID int, ID int) (bool, error) {
	var ans bool
	err := s.DB.QueryRow("SELECT result FROM friends WHERE id1=$1 AND id2=$2", userID, ID).Scan(&ans)

	if err != nil {
		return false, err
	}
	return ans, nil
}

func (s *Store) GetFriends(userID int) (*[]models.Friend, error) {
	var friends []models.Friend

	rows, err := s.DB.Query("SELECT friends.id2, users.fullname, users.email, friends.result FROM friends JOIN users ON friends.id2=users.id WHERE friends.id1=$1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var friend *models.Friend
		friend, err = utils.ScanRowToFriend(rows)
		if err != nil {
			return nil, err
		}
		friends = append(friends, *friend)
	}
	return &friends, nil
}

func (s *Store) GetRequests(userID int) (*[]models.Friend, error) {
	var friends []models.Friend

	rows, err := s.DB.Query("SELECT friends.id1, users.fullname, users.email, friends.result FROM friends JOIN users ON friends.id1=users.id WHERE friends.id2=$1 AND friends.result=$2", userID, false)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var friend *models.Friend
		friend, err = utils.ScanRowToFriend(rows)
		if err != nil {
			return nil, err
		}
		friends = append(friends, *friend)
	}
	return &friends, nil
}

func (s *Store) AddFriend(id1 int, id2 int, ans bool) error {
	_, err := s.DB.Exec("INSERT INTO friends (id1, id2, result) VALUES($1, $2, $3)", id1, id2, ans)
	return err
}

func (s *Store) AcceptFriend(id1 int, id2 int) error {
	_, err := s.DB.Exec("UPDATE friends SET result=$1 WHERE id1=$2 AND id2=$3", true, id1, id2)
	return err
}

func (s *Store) GetGroups(userID int) (*[]models.Group, error) {
	var groups []models.Group

	rows, err := s.DB.Query("SELECT groups.id, groups.name FROM groups JOIN group_members ON groups.id=group_members.group_id WHERE group_members.user_id=$1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var group *models.Group
		group, err = utils.ScanRowToGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	return &groups, nil
}

func (s *Store) AddGroup(group *models.Group, userID int) error {
	err := s.DB.QueryRow("INSERT INTO groups (name, created_by) VALUES($1, $2) RETURNING id", group.Name, userID).Scan(&group.ID)
	return err
}

func (s *Store) AddMember(groupID int, id int) error {
	_, err := s.DB.Exec("INSERT INTO group_members (group_id, user_id) VALUES($1, $2)", groupID, id)
	return err
}

func (s *Store) CheckMember(userID int, id int) (bool, error) {
	var t int
	err := s.DB.QueryRow("SELECT group_id FROM group_members where group_id=$1 AND user_id=$2", id, userID).Scan(&t)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) GetGroupDetails(id int) (*models.GroupDetails, error) {
	details := new(models.GroupDetails)
	err := s.DB.QueryRow("SELECT id, name, created_by, created_at FROM groups WHERE id=$1", id).Scan(&details.ID, &details.Name, &details.CreatedBy, &details.CreatedAt)
	if err != nil {
		return nil, err
	}

	var members []int
	rows, err := s.DB.Query("SELECT group_members.user_id FROM group_members JOIN groups ON group_members.group_id=groups.id WHERE group_members.group_id=$1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var member *int
		err := rows.Scan(&member)
		if err != nil {
			return nil, err
		}
		members = append(members, *member)
	}

	details.Members = members
	return details, nil
}

func (s *Store) AddExpense(e models.NewExpense) error {
	_, err := s.DB.Exec("INSERT INTO expenses (group_id, name, description, amount, paid_by, split_type) VALUES($1, $2, $3, $4, $5, $6)", e.GroupID, e.Name, e.Description, e.Amount, e.PaidBy, e.SplitType)
	return err
}

func (s *Store) GetExpenses(groupID int) (*[]models.Expense, error) {
	var expenses []models.Expense
	rows, err := s.DB.Query("SELECT * FROM expenses WHERE group_id=$1", groupID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var expense *models.Expense
		expense, err := utils.ScanRowToExpense(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}
	return &expenses, nil
}
