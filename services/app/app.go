package app

import (
	"auth/services/key"
)

type App struct {
	ID            	int    `json:"id" gorm:"primaryKey"`
	AppName       	string `json:"app_name"`
	AppURI        	string `json:"app_uri"`
	RedirectURI   	string `json:"redirect_uri"`
	ClientType    	string `json:"client_type"` //public or confidential. maybe create a table for it
	Alg				string `json:"alg"`
	ClientID		string `json:"client_id"`
	CreatedAt     	int64  `json:"created_at"`
	UpdatedAt     	int64  `json:"updated_at"`

	Keys []key.Key `json:"keys" gorm:"foreignKey:AppID;onDelete:CASCADE"` // Many keys belong to one app
}

type AppInput struct {
	AppName  		string `json:"app_name"`
	AppURI      	string `json:"app_uri"`
	Alg				string `json:"alg"`
	//RedirectURI   string `json:"redirect_uri"`
	//ClientType    string `json:"client_type"`
}

// type PasswordInput struct {
// 	ActivePassword string `json:"active_password" validate:"required,password"`
// 	NewPassword    string `json:"new_password" validate:"required,password"`
// 	RepeatPassword string `json:"repeat_password" validate:"required,password"`
// }
