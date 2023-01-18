package models

import (
	"time"
)

// Test model
type Test struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Passed    bool      `json:"passed"`
	Number    uint8     `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName gives table name of model
func (t Test) TableName() string {
	return "test"
}
