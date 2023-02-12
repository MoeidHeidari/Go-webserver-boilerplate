package repository

import (
	"main/lib"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestRepository database structure
type TestRepository struct {
	lib.Database
	logger lib.Logger
}

// NewTestRepository creates a new Test repository
func NewTestRepository(db lib.Database, logger lib.Logger) TestRepository {
	return TestRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r TestRepository) WithTrx(trxHandle *mongo.Collection) TestRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	//r.Database.collection = trxHandle
	return r
}
