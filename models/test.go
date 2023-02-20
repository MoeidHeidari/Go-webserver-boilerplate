package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coordinates struct {
	X int `bson:"x" json:"x"`
	Y int `bson:"y" json:"y"`
}

// Test model
type Node struct {
	ID            string      `bson:"id" json:"id"`
	CpuNumber     int         `bson:"cpuNumber" json:"cpu_number"`
	MemoryNumber  int         `bson:"memoryNumber" json:"memory_number"`
	StorageNumber int         `bson:"storageNumber" json:"storage_number"`
	Position      Coordinates `bson:"position" json:"position"`
	NodeName      string      `bson:"name" json:"name"`
	CardLabel     string      `bson:"cardLabel" json:"card_label"`
	LabelColor    string      `bson:"labelColor" json:"label_color"`
}

type Workspace struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Node   []Node             `bson:"nodes,omitempty" json:"nodes"`
	Name   string             `bson:"name" json:"name"`
	Edges  []Edge             `bson:"edges,omitempty" json:"edges"`
	Status bool               `bson:"status"`
}

type Workspaces struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Status string             `json:"status" bson:"status"`
}

type Edge struct {
	ID     string `bson:"id" json:"id"`
	Source string `bson:"source" json:"source"`
	Target string `bson:"target" json:"target"`
}

// TableName gives table name of model
func (t Node) TableName() string {
	return "test"
}
