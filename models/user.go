package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  string `json:"isAdmin,omitempty"`
}
