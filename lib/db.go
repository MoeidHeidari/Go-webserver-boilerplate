package lib

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database modal
type Database struct {
	Collection *mongo.Collection
}

// NewDatabase creates a new database instance
func NewDatabase(env Env, logger Logger) Database {

	// username := env.DBUsername
	// password := env.DBPassword
	//host := env.DBHost
	//port := env.DBPort
	dbname := env.DBName
	dbcollection := env.DBCollection
	dbUrl := env.DbUrl

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(dbUrl).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(err)
	}

	datab := client.Database(dbname)
	collection := datab.Collection(dbcollection)

	if err == mongo.ErrNoDocuments {
		logger.Info("No document was found")
	}
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Database connection established")
	return Database{
		Collection: collection,
	}
}
