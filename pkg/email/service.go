package email

import (
	"btradoc/entities"
	"fmt"

	"github.com/sirupsen/logrus"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	queue      chan entities.TranslatorResetPassword
	mongoDB    *mongo.Database
	SMTPClient *mail.SMTPClient
}

func NewService(mongoDB *mongo.Database) Service {
	// client := mail.NewSMTPClient()

	// client.Host = "smtp.gmail.com"
	// client.Port = 465
	// client.Username = "aaa@gmail.com"
	// client.Password = "asdfghh"
	// client.Encryption = mail.EncryptionSTARTTLS
	// client.ConnectTimeout = 10 * time.Second
	// client.SendTimeout = 10 * time.Second

	// //KeepAlive is not settted because by default is false

	// //Connect to client
	// smtpClient, err := client.Connect()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	emailService := Service{
		queue:      make(chan entities.TranslatorResetPassword),
		mongoDB:    mongoDB,
		SMTPClient: nil,
	}

	return emailService
}

func (es *Service) Mailer(logger *logrus.Logger) {
	go func() {
		for translResetPwd := range es.queue {
			// https://github.com/xhit/go-simple-mail/blob/master/example/example_test.go
			email := mail.NewMSG()

			email.SetFrom("From Example <from.email@example.com>").
				AddTo(translResetPwd.Email).
				SetSubject("Occitanofòn Traduccions: Reset Senhal")

			body := fmt.Sprintf(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /><title>Occitanofòn</title></head><body><p style="text-align:center;font-size:16px;font-weight:600">Adiu %s</p><p style="text-align:center;font-size:15px">Voici le lien permettant de procéder au reste du password</p> <p style="text-align:center;font-size:14px">https://occitanofòn.xyz/confirmacion/%s</p></body></html>`, translResetPwd.Username, translResetPwd.Token)

			//Get from each mail
			email.GetFrom()
			email.SetBody(mail.TextHTML, body)

			//Send with high priority
			email.SetPriority(mail.PriorityHigh)

			// always check error after send
			if email.Error != nil {
				logger.Error(email.Error)
				continue
			}

			if err := email.Send(es.SMTPClient); err != nil {
				logger.Error(err)
			}
		}
	}()
}

func (es *Service) SendResetPasswordLink(translResetPwd *entities.TranslatorResetPassword) {
	go func() {
		es.queue <- *translResetPwd
	}()
}
