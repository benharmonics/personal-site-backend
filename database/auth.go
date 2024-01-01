package database

import (
	"encoding/base64"
	"fmt"

	"github.com/benharmonics/backend/database/models"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(user *models.User, password string) error {
	if user == nil {
		return fmt.Errorf("cannot authenticate nil user")
	}
	passwordHash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
}
