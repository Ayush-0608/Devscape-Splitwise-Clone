package utils

import (
	"database/sql"
	"fmt"
	"splitwise/models"
)

type Store struct {
	DB *sql.DB
}

func (s *Store) GetUserByMail(email string) (*models.User, error) {
	rows, err := s.DB.Query("SELECT id, fullname, email, phone FROM users WHERE email=?", email)
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

func ScanRowToUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.Phone,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
