package services

import (
	"context"
	"fmt"
	"main/lib"
	"main/models"
	"main/repository"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestService service layer
type TestService struct {
	logger     lib.Logger
	repository repository.TestRepository
}

// NewTestService creates a new Testservice
func NewTestService(logger lib.Logger, repository repository.TestRepository) TestService {
	return TestService{
		logger:     logger,
		repository: repository,
	}
}

// GetOneTest gets one Test
func (s TestService) GetOneWorkspace(id primitive.ObjectID) (result models.Workspace, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.repository.Database.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (s TestService) GetDeletedWorkspaces() ([]models.Workspaces, error) {
	filter := bson.D{{}}
	curr, err := s.repository.Trash.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var Tests []models.Workspace
	if err = curr.All(context.TODO(), &Tests); err != nil {
		return nil, err
	}
	deleted := make([]models.Workspaces, len(Tests))
	for i, test := range Tests {
		deleted[i].ID = test.ID
		deleted[i].Name = test.Name
		if test.Status {
			deleted[i].Status = "online"
		} else {
			deleted[i].Status = "offline"
		}
	}
	return deleted, err
}

// GetAllTest get all the Test
func (s TestService) GetAllWorkspaces() ([]models.Workspaces, error) {
	filter := bson.D{{}}
	curr, err := s.repository.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var Tests []models.Workspace
	if err = curr.All(context.TODO(), &Tests); err != nil {
		return nil, err
	}
	response := make([]models.Workspaces, len(Tests))
	for i, test := range Tests {
		response[i].ID = test.ID
		response[i].Name = test.Name
		if test.Status {
			response[i].Status = "online"
		} else {
			response[i].Status = "offline"
		}
	}
	return response, err
}

func (u TestService) ValidateCardLabel(cardLabel string) error {
	allowedNames := []string{
		"PV",
		"VM",
		"pod",
		"ingress",
		"service",
		"storage",
		"claim",
		"endpoints",
		"nodejs",
		"PVC",
		"rules",
	}
	for _, label := range allowedNames {
		if strings.EqualFold(cardLabel, label) {
			return nil
		}
	}
	return errors.New("invalid CardLabel")
}

func (u TestService) ValidateCpuNumber(n int) error {
	if 1 <= n && n <= 16 {
		return nil
	}
	return errors.New("CpuNumber must be between 1 and 16")
}

func (u TestService) ValidateMemoryNumber(n int) error {
	if 8 <= n && n <= 64 {
		return nil
	}
	return errors.New("MemoryNumber must be between 8 and 64")
}

func (u TestService) ValidateStorageNumber(n int) error {
	if 8 <= n && n <= 64 {
		return nil
	}
	return errors.New("StorageNumber must be between 8 and 64")
}

func (u TestService) ValidateNode(n models.Node) error {
	err := u.ValidateCardLabel(n.CardLabel)
	if err != nil {
		return err
	}
	err = u.ValidateCpuNumber(n.CpuNumber)
	if err != nil {
		return err
	}
	err = u.ValidateMemoryNumber(n.MemoryNumber)
	if err != nil {
		return err
	}
	err = u.ValidateStorageNumber(n.StorageNumber)
	if err != nil {
		return err
	}
	return nil
}

// CreateTest call to create the Test
func (s TestService) CreateWorkspace(Workspace models.Workspace) (primitive.ObjectID, error) {
	Workspace.ID = primitive.NewObjectID()
	Workspace.Status = true
	_, err := s.repository.Collection.InsertOne(context.TODO(), Workspace)
	return Workspace.ID, err
}

// UpdateTest updates the Test
func (s TestService) AddNode(Node models.Node, id primitive.ObjectID) error {
	var workspace models.Workspace
	filter := bson.D{{Key: "_id", Value: id}}
	err := s.repository.Collection.FindOne(context.TODO(), filter).Decode(&workspace)
	if err != nil {
		return err
	}
	Node.ID = fmt.Sprint(len(workspace.Node) + 1)
	int_id, _ := strconv.Atoi(Node.ID)
	for _, node := range workspace.Node {
		if Node.ID == node.ID {
			Node.ID = fmt.Sprint(int_id - 1)
		}
	}
	workspace.Node = append(workspace.Node, Node)
	update := bson.D{{Key: "$set", Value: workspace}}
	_, err = s.repository.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (s TestService) AddEdge(Edge models.Edge, id primitive.ObjectID) error {
	var workspace models.Workspace
	filter := bson.D{{Key: "_id", Value: id}}
	err := s.repository.Collection.FindOne(context.TODO(), filter).Decode(&workspace)
	if err != nil {
		return err
	}
	Edge.ID = fmt.Sprintf("e%s-%s", Edge.Source, Edge.Target)
	workspace.Edges = append(workspace.Edges, Edge)
	update := bson.D{{Key: "$set", Value: workspace}}
	_, err = s.repository.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// DeleteTest deletes the Test
func (s TestService) DeleteWorkspace(id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	var workspace models.Workspace
	err := s.repository.Collection.FindOne(context.TODO(), filter).Decode(&workspace)
	if err != nil {
		return err
	}
	_, err = s.repository.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	workspace.Status = false
	_, err = s.repository.Trash.InsertOne(context.TODO(), workspace)
	if err != nil {
		return err
	}
	return err
}

func (s TestService) DeleteNode(workspace_id primitive.ObjectID, id string) error {
	node_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	node_id -= 1
	filter := bson.D{{Key: "_id", Value: workspace_id}}
	var workspace models.Workspace
	err = s.repository.Collection.FindOne(context.TODO(), filter, nil).Decode(&workspace)
	if err != nil {
		return err
	}
	workspace.Node = append(workspace.Node[:node_id], workspace.Node[node_id+1:]...)

	update := bson.D{{Key: "$set", Value: workspace}}
	_, err = s.repository.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (s TestService) UpdateNode(workspace_id primitive.ObjectID, node models.Node) error {
	node_id, err := strconv.Atoi(node.ID)
	if err != nil {
		return err
	}
	node_id -= 1
	filter := bson.D{{Key: "_id", Value: workspace_id}}
	var workspace models.Workspace
	err = s.repository.Collection.FindOne(context.TODO(), filter, nil).Decode(&workspace)
	if err != nil {
		return err
	}
	workspace.Node[node_id] = node
	update := bson.D{{Key: "$set", Value: workspace}}
	_, err = s.repository.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}
