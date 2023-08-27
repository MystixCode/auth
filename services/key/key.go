package key

type Key struct {
	ID         	int		`json:"id" gorm:"primaryKey"`
	AppID		int		`json:"app_id" gorm:"not null"`
	CreatedAt	int64	`json:"created_at" gorm:"not null"`
}

type KeyInput struct {
	AppID int		`json:"app_id" validate:"required,number"`
}
