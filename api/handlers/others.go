package handlers

import (
	"database/sql"
	"fmt"
	"splitwise/models"
)

type Store struct {
	DB *sql.DB
}

func ScanRowToUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ScanRowToPublicUser(rows *sql.Rows) (*models.PublicUser, error) {
	user := new(models.PublicUser)
	err := rows.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByMail(email string) (*models.User, error) {
	rows, err := s.DB.Query("SELECT id, fullname, email, phone, password, created_at FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	for rows.Next() {
		user, err = ScanRowToUser(rows)
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
		user, err = ScanRowToUser(rows)
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

func (s *Store) GetUsers() (*[]models.PublicUser, error) {
	var users []models.PublicUser
	rows, err := s.DB.Query("SELECT id, fullname, email, phone, password, created_at FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user *models.PublicUser
		user, err = ScanRowToPublicUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return &users, nil
}
