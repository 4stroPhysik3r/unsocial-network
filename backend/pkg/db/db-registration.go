package db

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID          int    `json:"userID"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Avatar          string `json:"avatar"`
	Nickname        string `json:"nickname"`
	AboutMe         string `json:"aboutMe"`
	ProfilePublic   bool   `json:"profilePublic,omitempty"`
	CreatedAt       string `json:"createdAt"`
}

func RegisterUser(userData User) error {

	// Check if the email is already taken
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	err := DB.QueryRow(query, userData.Email).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking username existence: %w", err)
	}
	if count > 0 {
		return errors.New("email is already taken")
	}

	// Hash the password
	passwordHash, err := HashPassword(userData.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	query = `INSERT INTO users (email, password, firstname, lastname, date_of_birth, avatar, nickname, about_me, profile_public) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Insert the new user into the database
	_, err = DB.Exec(query, userData.Email, passwordHash, userData.FirstName, userData.LastName,
		userData.DateOfBirth, userData.Avatar, userData.Nickname, userData.AboutMe, userData.ProfilePublic)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %w", err)
	}

	return nil
}

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
