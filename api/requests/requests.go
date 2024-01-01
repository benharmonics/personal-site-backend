package requests

import (
	"github.com/benharmonics/backend/utils/validation"
	"github.com/go-playground/validator/v10"
)

type (
	NewBlogPost struct {
		Title    string  `validate:"required,max=200" json:"title"`
		Subtitle *string `validate:"omitempty,max=200" json:"subtitle"`
		Author   string  `validate:"required,max=200" json:"author"`
		Content  string  `validate:"required,max=5000" json:"content"`
	}

	NewUser struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required,bcrypt" json:"password"`
	}

	Login struct{ NewUser }
)

func canBeEncrypted(fl validator.FieldLevel) bool {
	return len([]byte(fl.Field().String())) <= 72
}

func (req NewBlogPost) Validate() error {
	return validation.ValidateStruct(req, nil)
}

func (req NewUser) Validate() error {
	val := validator.New()
	val.RegisterValidation("bcrypt", canBeEncrypted)
	return validation.ValidateStruct(req, val)
}
