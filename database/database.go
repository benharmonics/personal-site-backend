package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/benharmonics/personal-site-backend/logging"
)

// Database represents our connection to MongoDB
type Database struct{ client *mongo.Client }

// NewDatabase creates a new MongoDB connection.
func NewDatabase(opts ...Option) (*Database, error) {
	connectOpt := ConnectOption{host: "localhost", port: 27017}
	for _, f := range opts {
		f(&connectOpt)
	}
	db := &Database{}
	if err := db.connect(connectOpt); err != nil {
		return nil, err
	}
	return db, nil
}

// Disconnect attempts to disconnect from MongoDB.
// It's best practice to close open connections when you're done using them.
func (db *Database) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := db.client.Disconnect(ctx); err != nil {
		logging.Error("Failed to disconnect from database:", err)
	}
}
