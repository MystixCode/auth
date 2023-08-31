package services

import (
	"auth/services/example"
	"auth/services/health"
	"auth/services/root"
	"auth/services/auth"
)

// Services holds all available services
type Services struct {
	Root    	*root.Service
	Health  	*health.Service
	Auth    	*auth.Service
	Example 	*example.Service
}
