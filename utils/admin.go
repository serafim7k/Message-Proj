package utils

import (
	"database/sql"
)

type User struct {
	ID       int
	Username string
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username); err == nil {
			users = append(users, u)
		}
	}
	return users, nil
}

func DeleteUser(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func IsAdmin(db *sql.DB, userID string) bool {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=$1", userID).Scan(&username)
	return err == nil && username == "admin"
}
