package models

import (
	"fmt"
	"log"
	"messenger/database"
	"time"
)

// User represents the user entity
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"` // Omit password when serializing to JSON
	AuthToken string    `json:"auth_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetUserByUsername(username string) (*User, error) {
	currentUser := &User{}
	db := database.GetDB()
	qry := fmt.Sprintf("SELECT id, password, username, auth_token, created_at, updated_at "+
		"FROM users WHERE username = '%s'", username)
	err := db.QueryRow(qry).Scan(&currentUser.ID, &currentUser.Password, &currentUser.Username,
		&currentUser.AuthToken, &currentUser.CreatedAt, &currentUser.UpdatedAt)
	if err != nil {
		log.Println(fmt.Sprintf("Authentication failed for user: %s qry: %s error: %v\n", username, qry, err))
		return nil, err
	}
	return currentUser, nil
}
