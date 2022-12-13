package repository

import (
	"github.com/dipeshdulal/clean-gin/lib"
	"gorm.io/gorm"
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
func (r TestRepository) WithTrx(trxHandle *gorm.DB) TestRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = trxHandle
	return r
}
