package models

// User represents a single user
type User struct {
	Login    string
	Password string
}

// LoginForm is a struct representing the form to login.
type LoginForm struct {
	Login    string `form:"login" binding:"required"`
	Password string `form:"password" binding:"required"`
}
