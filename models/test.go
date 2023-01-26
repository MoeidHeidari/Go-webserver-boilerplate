package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test model
type Test struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Passed    bool               `bson:"passed"`
	Number    uint8              `bson:"number"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// TableName gives table name of model
func (t Test) TableName() string {
	return "test"
}
