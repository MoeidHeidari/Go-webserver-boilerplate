package services

import (
	"context"
	"main/lib"
	"main/models"
	"main/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// // WithTrx delegates transaction to repository database
// func (s TestService) WithTrx(trxHandle *gorm.DB) TestService {
// 	s.repository = s.repository.WithTrx(trxHandle)
// 	return s
// }

// GetOneTest gets one Test
func (s TestService) GetOneTest(id primitive.ObjectID) (result models.Test, err error) {
	filter := bson.D{{"_id", id}}
	err = s.repository.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		s.logger.Fatal(err)
	}
	return result, err
}

// GetAllTest get all the Test
func (s TestService) GetAllTest() (Tests []models.Test, err error) {
	filter := bson.D{{}}
	curr, err := s.repository.Collection.Find(context.TODO(), filter)
	if err != nil {
		s.logger.Fatal(err)
	}
	if err = curr.All(context.TODO(), &Tests); err != nil {
		panic(err)
	}

	return Tests, err
}

// GetAllTestField get all the Field
func (s TestService) GetAllTestField(field_name string) (Tests []string, err error) {
	filter := bson.D{{}}
	opts := options.Find().SetProjection(bson.D{{field_name, 0}})
	curr, err := s.repository.Collection.Find(context.TODO(), filter, opts)

	if err != nil {
		s.logger.Fatal(err)
	}

	if err = curr.All(context.TODO(), &Tests); err != nil {
		panic(err)
	}

	return Tests, err
}

// CreateTest call to create the Test
func (s TestService) CreateTest(Test models.Test) error {

	_, err := s.repository.Collection.InsertOne(context.TODO(), Test)
	return err
}

// UpdateTest updates the Test
func (s TestService) UpdateTest(Test models.Test) error {
	id := Test.ID
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", Test}}
	_, err := s.repository.Collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// DeleteTest deletes the Test
func (s TestService) DeleteTest(id string) error {
	filter := bson.D{{"_id", id}}
	_, err := s.repository.Collection.DeleteOne(context.TODO(), filter)
	return err
}
