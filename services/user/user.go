package user

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Hash      string `json:"hash"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type UserInput struct {
	UserName  string `json:"user_name" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"omitempty,min=2,alphanum"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,alphanum"`
	Hash      string `json:"hash" validate:"required,alphanum"`
}

type LoginInput struct {
	Email     string `json:"email" validate:"required,email"`
	Hash      string `json:"hash" validate:"required,alphanum"`
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// type PasswordInput struct {
// 	ActivePassword string `json:"active_password" validate:"required,password"`
// 	NewPassword    string `json:"new_password" validate:"required,password"`
// 	RepeatPassword string `json:"repeat_password" validate:"required,password"`
// }
