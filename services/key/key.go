package key

type Key struct {
	ID         	int		`json:"id" gorm:"primaryKey"`
	AppID		int		`json:"app_id" gorm:"not null"`
	Alg			string	`json:"alg" gorm:"not null"`
	CreatedAt	int64	`json:"created_at" gorm:"not null"`
}

type KeyInput struct {
	AppID int		`json:"app_id" validate:"required,number"`
	Alg		string	`json:"alg" validate:"alphanumeric"`
}
