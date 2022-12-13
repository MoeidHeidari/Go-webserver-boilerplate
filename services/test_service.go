package services

import (
	"github.com/dipeshdulal/clean-gin/lib"
	"gorm.io/gorm"
	"main/models"
	"main/repository"
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

// WithTrx delegates transaction to repository database
func (s TestService) WithTrx(trxHandle *gorm.DB) TestService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneTest gets one Test
func (s TestService) GetOneTest(id uint) (Test models.Test, err error) {
	return Test, s.repository.Find(&Test, id).Error
}

// GetAllTest get all the Test
func (s TestService) GetAllTest() (Tests []models.Test, err error) {
	return Tests, s.repository.Find(&Tests).Error
}

// CreateTest call to create the Test
func (s TestService) CreateTest(Test models.Test) error {
	return s.repository.Create(&Test).Error
}

// UpdateTest updates the Test
func (s TestService) UpdateTest(Test models.Test) error {
	return s.repository.Save(&Test).Error
}

// DeleteTest deletes the Test
func (s TestService) DeleteTest(id uint) error {
	return s.repository.Delete(&models.Test{}, id).Error
}
