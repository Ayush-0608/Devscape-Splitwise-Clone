package utils

import (
	"database/sql"
	"splitwise/models"
)

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

func ScanRowToFriend(rows *sql.Rows) (*models.Friend, error) {
	friend := new(models.Friend)
	err := rows.Scan(
		&friend.ID,
		&friend.Fullname,
		&friend.Email,
		&friend.Result,
	)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func ScanRowToGroup(rows *sql.Rows) (*models.Group, error) {
	group := new(models.Group)
	err := rows.Scan(
		&group.ID,
		&group.Name,
	)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func ScanRowToExpense(rows *sql.Rows) (*models.Expense, error) {
	expense := new(models.Expense)
	err := rows.Scan(
		&expense.ID,
		&expense.GroupID,
		&expense.Name,
		&expense.Description,
		&expense.Amount,
		&expense.PaidBy,
		&expense.SplitType,
		&expense.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return expense, nil
}
