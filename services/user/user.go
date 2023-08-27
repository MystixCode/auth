package user

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserName  string `json:"user_name" gorm:"not null"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hash      string `json:"hash" gorm:"not null"`
	CreatedAt int64  `json:"created_at" gorm:"not null"`
	UpdatedAt int64  `json:"updated_at"`
}

type UserInput struct {
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,alphanum"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,alphanum"`
	Hash      string `json:"hash" validate:"required,base64"`
}

type LoginInput struct {
	ClientID  string `json:"client_id" validate:"required"`
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Hash      string `json:"hash" validate:"required,base64"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MyCustomClaims struct {
	Username string `json:"username"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

// type PasswordInput struct {
// 	ActivePassword string `json:"active_password" validate:"required,password"`
// 	NewPassword    string `json:"new_password" validate:"required,password"`
// 	RepeatPassword string `json:"repeat_password" validate:"required,password"`
// }
