package example

// import "go.mongodb.org/mongo-driver/bson/primitive"

type Example struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	ExampleName  string `json:"examplename"`
	ExampleValue string `json:"examplevalue"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type ExampleInput struct {
	ExampleName  string `json:"examplename"`
	ExampleValue string `json:"examplevalue"`
}
