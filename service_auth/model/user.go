package model

// User ...
type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}
