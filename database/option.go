package database

import (
	"context"
	"fmt"
	"time"

	"github.com/benharmonics/personal-site-backend/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	Option func(*ConnectOption)

	ConnectOption struct {
		host      string
		port      int
		username  *string
		password  *string
		encrypted bool
	}
)

func WithEncryptedConnection() Option {
	return func(db *ConnectOption) { db.encrypted = true }
}

func WithHost(host string) Option {
	return func(db *ConnectOption) { db.host = host }
}

func WithPort(port int) Option {
	return func(db *ConnectOption) { db.port = port }
}

func WithoutPort() Option {
	return func(db *ConnectOption) { db.port = 0 }
}

func WithCredentials(username, password string) Option {
	return func(db *ConnectOption) {
		db.username = &username
		db.password = &password
	}
}

func (db *Database) connect(opt ConnectOption) error {
	logging.Info("Connecting to MongoDB at", opt.host)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wc := &writeconcern.WriteConcern{W: "majority"}
	mongoOpts := options.Client().ApplyURI(mongodbURI(opt)).SetRetryWrites(true).SetWriteConcern(wc)
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		return err
	}
	db.client = client
	return nil
}

func mongodbURI(opt ConnectOption) string {
	// proto
	uri := "mongodb"
	if opt.encrypted {
		uri += "+srv"
	}
	uri += "://"
	// credentials
	if opt.username != nil {
		uri += *opt.username
		if opt.password != nil {
			uri += fmt.Sprintf(":%s", *opt.password)
		}
		uri += "@"
	}
	// host & port
	uri += opt.host
	if !opt.encrypted { // Port is not allowed with +srv connections
		uri += fmt.Sprintf(":%d", opt.port)
	}
	logging.Debug("MongoDB URI:", uri)
	return uri
}
