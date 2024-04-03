package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/us190190/messenger/contracts"
	"github.com/us190190/messenger/database"
	"github.com/us190190/messenger/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func UserRegisterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body into contracts.UpdateUserRequest struct
	var createUserReqBody contracts.CreateUserRequest
	err := json.NewDecoder(request.Body).Decode(&createUserReqBody)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(responseWriter, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	timeNow := time.Now()

	randomToken := make([]byte, 32)
	_, err = rand.Read(randomToken)
	if err != nil {
		http.Error(responseWriter, "Failed to create user", http.StatusInternalServerError)
		return
	}
	authToken := base64.URLEncoding.EncodeToString(randomToken)

	// Insert user into database
	db := database.GetDB()
	_, err = db.Exec("INSERT INTO users (username, password, auth_token, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", createUserReqBody.Username, hashedPassword, authToken, timeNow, timeNow)
	if err != nil {
		http.Error(responseWriter, "Failed to create user", http.StatusInternalServerError)
		return
	}

	userResponse := contracts.UserResponse{AuthToken: authToken, Username: createUserReqBody.Username, CreatedAt: timeNow, UpdatedAt: timeNow}

	// Convert struct to JSON
	jsonData, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(responseWriter, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}

func UserAuthenticationHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON request body into contracts.UpdateUserRequest struct
	var authUserReqBody contracts.CreateUserRequest
	err := json.NewDecoder(request.Body).Decode(&authUserReqBody)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Retrieve user from database
	var storedPassword string
	var currentUser models.User
	db := database.GetDB()
	err = db.QueryRow("SELECT password, username, auth_token, created_at, updated_at FROM users WHERE username = ?", authUserReqBody.Username).Scan(&storedPassword, &currentUser.Username, &currentUser.AuthToken, &currentUser.CreatedAt, &currentUser.UpdatedAt)
	if err != nil {
		http.Error(responseWriter, "Username or password incorrect", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(authUserReqBody.Password))
	if err != nil {
		http.Error(responseWriter, "Username or password incorrect", http.StatusUnauthorized)
		return
	}

	userResponse := contracts.UserResponse{Username: authUserReqBody.Username, AuthToken: currentUser.AuthToken, CreatedAt: currentUser.CreatedAt, UpdatedAt: currentUser.UpdatedAt}

	// Convert struct to JSON
	jsonData, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(responseWriter, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}

func UserUpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authToken := strings.Split(request.Header.Get("Authorization"), "Bearer ")[1]

	var currentUser models.User

	// Validate user token from database
	db := database.GetDB()
	err := db.QueryRow("SELECT id, created_at FROM users WHERE auth_token = ?", authToken).Scan(&currentUser.ID, &currentUser.CreatedAt)
	if err != nil {
		http.Error(responseWriter, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// Decode JSON request body into contracts.UpdateUserRequest struct
	var updateUserReqBody contracts.UpdateUserRequest
	err = json.NewDecoder(request.Body).Decode(&updateUserReqBody)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	updatedTimeNow := time.Now()

	// Update username into database
	_, err = db.Exec("UPDATE users SET username = ?, updated_at = ? where id = ?", updateUserReqBody.NewUsername, updatedTimeNow, currentUser.ID)
	if err != nil {
		http.Error(responseWriter, "Failed to update user", http.StatusInternalServerError)
		return
	}

	userResponse := contracts.UserResponse{AuthToken: authToken, Username: updateUserReqBody.NewUsername, CreatedAt: currentUser.CreatedAt, UpdatedAt: updatedTimeNow}

	// Convert struct to JSON
	jsonData, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(responseWriter, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}

func UserRemoveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodDelete {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authToken := strings.Split(request.Header.Get("Authorization"), "Bearer ")[1]

	var currentUser models.User

	// Validate user token from database
	db := database.GetDB()
	err := db.QueryRow("SELECT id, created_at FROM users WHERE auth_token = ?", authToken).Scan(&currentUser.ID, &currentUser.CreatedAt)
	if err != nil {
		http.Error(responseWriter, "Invalid token!", http.StatusUnauthorized)
		return
	}

	// Decode JSON request body into contracts.UpdateUserRequest struct
	var updateUserReqBody contracts.UpdateUserRequest
	err = json.NewDecoder(request.Body).Decode(&updateUserReqBody)
	if err != nil {
		http.Error(responseWriter, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Remove user from database
	_, err = db.Exec("DELETE FROM users where id = ?", currentUser.ID)
	if err != nil {
		http.Error(responseWriter, "Failed to remove user", http.StatusInternalServerError)
		return
	}

	commonResponse := contracts.CommonResponse{Message: "User removed successfully!"}

	// Convert struct to JSON
	jsonData, err := json.Marshal(commonResponse)
	if err != nil {
		http.Error(responseWriter, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}

func UserAllHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := database.GetDB()
	rows, err := db.Query("SELECT id, username, auth_token, created_at, updated_at FROM users")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// log  fatal error
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var curUser models.User
		err := rows.Scan(&curUser.ID, &curUser.Username, &curUser.AuthToken, &curUser.CreatedAt, &curUser.UpdatedAt)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, curUser)
	}
	if err := rows.Err(); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}

func GroupAllHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := database.GetDB()
	rows, err := db.Query("SELECT id, group_name FROM `groups`")
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// log  fatal error
		}
	}(rows)

	var groups []models.Group
	for rows.Next() {
		var curGrp models.Group
		err := rows.Scan(&curGrp.ID, &curGrp.GroupName)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		groups = append(groups, curGrp)
	}
	if err := rows.Err(); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(groups)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	responseWriter.Header().Set("Content-Type", "application/json")

	// Write JSON response
	_, err = responseWriter.Write(jsonData)
	if err != nil {
		return
	}
}
