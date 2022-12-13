package models

// Test model
type Test struct {
	ID       uint    `json:"id"`
	TestName string  `json:"test_name"`
	Passed   *string `json:"passed"`
	Number   uint8   `json:"number"`
}

// TableName gives table name of model
func (t Test) TableName() string {
	return "test"
}
