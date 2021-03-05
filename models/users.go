package models

// Users struct
type Users struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	FullName string `json:"fullname" db:"fullname"`
}
