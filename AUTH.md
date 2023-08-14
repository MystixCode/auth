``` golang
package model

type App struct {
	ID          uint   `gorm:"primaryKey"`
	AppName     string `json:"app_name"`
	AppURL      string `json:"app_url"`
	RedirectURL string `json:"redirect_url"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`

	// Associations
	AccessTokens       []AccessToken       `json:"-" gorm:"foreignKey:AppID"`
	RefreshTokens      []RefreshToken      `json:"-" gorm:"foreignKey:AppID"`
	RevokedTokens      []RevokedToken      `json:"-" gorm:"foreignKey:AppID"`
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	// Associations
	AccessTokens  []AccessToken  `json:"-" gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `json:"-" gorm:"foreignKey:UserID"`
	RevokedTokens []RevokedToken `json:"-" gorm:"foreignKey:UserID"`
	Roles         []Role         `json:"roles" gorm:"many2many:user_roles;"`
}

type AccessToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	AppID     uint   `json:"app_id"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	AppID     uint   `json:"app_id"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RevokedToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `json:"token"`
	UserID    uint   `json:"user_id"`
	AppID     uint   `json:"app_id"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
}

type Scope struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
}

type ResourceServer struct {
    ID         uint   `gorm:"primaryKey"`
    Name       string `json:"name"`
    PubKeyPath string `json:"pub_key_path"` // Path to the public key file
    PrivateKey string `json:"private_key"`  // Path to the private key file
    Scopes     []Scope `json:"scopes" gorm:"many2many:resource_server_scopes;"`
}

type ResourceServerScope struct {
	ResourceServerID uint `gorm:"primaryKey"`
	ScopeID          uint `gorm:"primaryKey"`
}

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	RoleName  string `json:"role_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	// Associations
	Users []User `json:"users" gorm:"many2many:user_roles;"`
}

type RoleScope struct {
	RoleID  uint `gorm:"primaryKey"`
	ScopeID uint `gorm:"primaryKey"`
}


// App (like web_client or game_client)
v1.HandleFunc("/apps", a.App.Create).Methods(http.MethodPost)
v1.HandleFunc("/apps", a.App.GetAll).Methods(http.MethodGet)
v1.HandleFunc("/apps/{id}", a.App.GetById).Methods(http.MethodGet)
v1.HandleFunc("/apps/{id}", a.App.Update).Methods(http.MethodPut)
v1.HandleFunc("/apps/{id}", a.App.Delete).Methods(http.MethodDelete)

// User
v1.HandleFunc("/users", a.User.Create).Methods(http.MethodPost)
v1.HandleFunc("/users", a.User.GetAll).Methods(http.MethodGet)
v1.HandleFunc("/users/{id}", a.User.GetById).Methods(http.MethodGet)
v1.HandleFunc("/users/{id}", a.User.Update).Methods(http.MethodPut)
v1.HandleFunc("/users/{id}", a.User.Delete).Methods(http.MethodDelete)
v1.HandleFunc("/users/login", a.User.Login).Methods(http.MethodPost)

// Role
v1.HandleFunc("/roles", a.Role.Create).Methods(http.MethodPost)
v1.HandleFunc("/roles", a.Role.GetAll).Methods(http.MethodGet)
v1.HandleFunc("/roles/{id}", a.Role.GetById).Methods(http.MethodGet)
v1.HandleFunc("/roles/{id}", a.Role.Update).Methods(http.MethodPut)
v1.HandleFunc("/roles/{id}", a.Role.Delete).Methods(http.MethodDelete)

// Scope
v1.HandleFunc("/scopes", a.Scope.Create).Methods(http.MethodPost)
v1.HandleFunc("/scopes", a.Scope.GetAll).Methods(http.MethodGet)
v1.HandleFunc("/scopes/{id}", a.Scope.GetById).Methods(http.MethodGet)
v1.HandleFunc("/scopes/{id}", a.Scope.Update).Methods(http.MethodPut)
v1.HandleFunc("/scopes/{id}", a.Scope.Delete).Methods(http.MethodDelete)

// Server (resource_server. like api or game_server)
v1.HandleFunc("/servers", a.Server.Create).Methods(http.MethodPost)
v1.HandleFunc("/servers", a.Server.GetAll).Methods(http.MethodGet)
v1.HandleFunc("/servers/{id}", a.Server.GetById).Methods(http.MethodGet)
v1.HandleFunc("/servers/{id}", a.Server.Update).Methods(http.MethodPut)
v1.HandleFunc("/servers/{id}", a.Server.Delete).Methods(http.MethodDelete)

// ...

```
