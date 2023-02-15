package services

import (
	"context"
	"fmt"
	"main/lib"
	"main/models"
	"main/repository"
	"strconv"

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
		s.logger.Fatal(err)
	}
	return result, err
}

// GetAllTest get all the Test
func (s TestService) GetAllWorkspaces() (Tests []models.Workspace, err error) {
	filter := bson.D{{}}
	curr, err := s.repository.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = curr.All(context.TODO(), &Tests); err != nil {
		return nil, err
	}

	return Tests, err
}

// CreateTest call to create the Test
func (s TestService) CreateWorkspace(Workspace models.Workspace) (primitive.ObjectID, error) {
	Workspace.ID = primitive.NewObjectID()
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
	_, err := s.repository.Collection.DeleteOne(context.TODO(), filter)
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
