package services

import (
	"auth/services/example"
	"auth/services/health"
	"auth/services/root"
	"auth/services/user"
	"auth/services/app"
	"auth/services/authorize"
)

// Services holds all available services
type Services struct {
	Root    	*root.Service
	Health  	*health.Service
	User    	*user.Service
	Example 	*example.Service
	App 		*app.Service
	Authorize 	*authorize.Service
}
