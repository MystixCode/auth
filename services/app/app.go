package app

type App struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	AppName       string `json:"app_name"`
	AppURL        string `json:"app_url"`
	RedirectURL   string `json:"redirect_url"`
	ClientType    string `json:"client_type"` //public or confidential. maybe create a table for it
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

type AppInput struct {
	AppName  string `json:"app_name"`
	AppURL        string `json:"app_url"`
	//RedirectURL   string `json:"redirect_url"`
	//ClientType    string `json:"client_type"`
}

// type PasswordInput struct {
// 	ActivePassword string `json:"active_password" validate:"required,password"`
// 	NewPassword    string `json:"new_password" validate:"required,password"`
// 	RepeatPassword string `json:"repeat_password" validate:"required,password"`
// }
