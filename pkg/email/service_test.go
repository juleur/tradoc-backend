package email

import (
	"btradoc/entities"
	"btradoc/storage/mongodb"
	"testing"
)

func TestEmailService(t *testing.T) {
	db := mongodb.NewMongoClient()
	emailService := NewService(db)
	emailService.Mailer(nil) // smtp disable for now

	translTest := &entities.TranslatorResetPassword{
		Email:    "test@test.com",
		Username: "test",
		Token:    "4a5z11vf24e5",
	}

	for i := 10_000; i > 0; i-- {
		emailService.SendResetPasswordLink(translTest)
	}
}
