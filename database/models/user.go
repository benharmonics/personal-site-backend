package models

import (
	"time"

	"github.com/benharmonics/backend/utils/auth"
	"github.com/benharmonics/backend/utils/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `validate:"required" bson:"_id" json:"id"`
	Email        string             `validate:"required,email" bson:"email" json:"email"`
	PasswordHash string             `validate:"required" bson:"passwordHash" json:"-"`
	LastLogin    time.Time          `validate:"required" bson:"lastLogin" json:"lastLogin"`
}

func NewUser(email, password string) (*User, error) {
	hash := auth.HashedPassword(password, auth.WithCPUThreads)
	user := &User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: hash,
		LastLogin:    time.Now(),
	}
	if err := validation.ValidateStruct(user, nil); err != nil {
		return nil, err
	}
	return user, nil
}
