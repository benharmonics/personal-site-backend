package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/benharmonics/personal-site-backend/logging"
)

// Database represents our connection to MongoDB
type Database struct {
	client    *mongo.Client
	host      string
	port      int
	username  *string
	password  *string
	encrypted bool
}

// NewDatabase creates a new MongoDB connection.
func NewDatabase(opts ...Option) (*Database, error) {
	db := &Database{host: "localhost", port: 27017}
	for _, f := range opts {
		f(db)
	}
	if err := db.setMongoDBClient(); err != nil {
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

func (db *Database) setMongoDBClient() error {
	// proto
	uri := "mongodb"
	if db.encrypted {
		uri += "+srv"
	}
	uri += "://"
	logging.Info("Connecting to MongoDB at", uri+db.host)
	// credentials
	if db.username != nil {
		uri += *db.username
		if db.password != nil {
			uri += fmt.Sprintf(":%s", *db.password)
		}
		uri += "@"
	}
	// host & port
	uri += db.host
	if !db.encrypted { // Port is not allowed with +srv connections
		uri += fmt.Sprintf(":%d", db.port)
	}
	logging.Debug("MongoDB URI:", uri)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wc := &writeconcern.WriteConcern{W: writeconcern.Majority()}
	opts := options.Client().ApplyURI(uri).SetRetryWrites(true).SetWriteConcern(wc)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	db.client = client
	return nil
}
