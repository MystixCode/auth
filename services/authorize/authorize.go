package authorize
type AuthorizationCode struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `json:"code"`
	UserID      uint   `json:"user_id"`  // Foreign key referencing the User table
	AppID       uint   `json:"app_id"`   // Foreign key referencing the App table
	RedirectURL	string `json:"redirect_uri"`
	Expiry      int64  `json:"expiry"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type AuthorizationCodeInput struct {
	Code        string `json:"code"`
	UserID      uint   `json:"user_id"`  // Foreign key referencing the User table
	AppID       uint   `json:"app_id"`   // Foreign key referencing the App table
	RedirectURL string `json:"redirect_uri"`
	Expiry      int64  `json:"expiry"`
}






// type Authorize struct {
// 	ID          			uint   `gorm:"primaryKey"`
// 	ClientID        		string	`json:"client_id"`
// 	RedirectURL      		string  `json:"redirect_url"`  // Foreign key referencing the User table
// 	ResponseType       		string  `json:"response_type"`   // Foreign key referencing the App table
// 	Scope 					string	`json:"scope"`
// 	State		 			string  `json:"state"`
// 	CodeChallenge    		string  `json:"code_challenge"`
// 	CodeChallengeMethod		string  `json:"code_challenge_method"`
// }

type AuthorizeInput struct {
	ClientID        		string	`json:"client_id"`
	RedirectURL      		string  `json:"redirect_url"`  // Foreign key referencing the User table
	ResponseType       		string  `json:"response_type"`   // Foreign key referencing the App table
	Scope 					string	`json:"scope"`
	State		 			string  `json:"state"`
	CodeChallenge    		string  `json:"code_challenge"`
	CodeChallengeMethod		string  `json:"code_challenge_method"`
}
