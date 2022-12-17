package models

// Test model
type Test struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Number uint8  `json:"number"`
}

// TableName gives table name of model
func (t Test) TableName() string {
	return "test"
}
