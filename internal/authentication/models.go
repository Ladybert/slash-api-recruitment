package authentication

import "github.com/go-playground/validator/v10"

type User struct {
    ID       string `json:"id"`
    Username string `json:"username" binding:"required,min=3,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=100"`
}

var validate = validator.New()
