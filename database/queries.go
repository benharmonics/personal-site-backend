package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/benharmonics/personal-site-backend/database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	database = "personal-server"

	blogsCollection = "blogs"
	usersCollection = "users"
)

func (db *Database) InsertBlog(blog models.BlogPost) error {
	coll := db.client.Database(database).Collection(blogsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := coll.InsertOne(ctx, blog)
	return err
}

func (db *Database) FindBlog(filter any, opts *options.FindOneOptions) (*models.BlogPost, error) {
	post := &models.BlogPost{}
	coll := db.client.Database(database).Collection(blogsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := coll.FindOne(ctx, filter, opts).Decode(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (db *Database) FindBlogs(filter any, opts *options.FindOptions) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	coll := db.client.Database(database).Collection(blogsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := coll.Find(ctx, filter, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []models.BlogPost{}, nil
	} else if err != nil {
		return nil, err
	}
	if err := res.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (db *Database) InsertUser(user *models.User) error {
	if user == nil {
		return fmt.Errorf("cannot insert nil user")
	}
	coll := db.client.Database(database).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := coll.FindOne(ctx, bson.M{"email": user.Email}).Err(); err == nil {
		return fmt.Errorf("user already exists")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	_, err := coll.InsertOne(ctx, user)
	return err
}

func (db *Database) FindUser(email string) (*models.User, error) {
	user := &models.User{}
	coll := db.client.Database(database).Collection(usersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := coll.FindOne(ctx, bson.M{"email": email}).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}
