package mailer

import (
	"btradoc/entities"
	"log"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const RETRY_DELAY time.Duration = 15 * time.Minute

func NewService(mongoDB *mongo.Database, logrus *logrus.Logger) Service {
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

	antiSpammer := &antiSpammer{
		lastTrySeen: make(map[string]time.Time),
		mu:          &sync.RWMutex{},
	}
	go antiSpammer.removeTryAfterSomeTime()

	mailerService := Service{
		mongoDB:     mongoDB,
		logger:      logrus,
		queue:       make(chan entities.TranslatorResetPassword),
		antiSpammer: antiSpammer,
		// SMTPClient:  nil,
	}
	go mailerService.sender()

	return mailerService
}

type Service struct {
	mongoDB     *mongo.Database
	logger      *logrus.Logger
	queue       chan entities.TranslatorResetPassword
	antiSpammer *antiSpammer
	// SMTPClient  *mail.SMTPClient
}

func (es *Service) sender() {
	for translResetPwd := range es.queue {
		log.Println(translResetPwd)
		// https://github.com/xhit/go-simple-mail/blob/master/example/example_test.go
		// email := mail.NewMSG()

		// email.SetFrom("From Example <from.email@example.com>").
		// 	AddTo(translResetPwd.Email).
		// 	SetSubject("Occitanofòn Traduccions: Reset Senhal")

		// body := fmt.Sprintf(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /><title>Occitanofòn</title></head><body><p style="text-align:center;font-size:16px;font-weight:600">Adiu %s</p><p style="text-align:center;font-size:15px">Voici le lien permettant de procéder au reste du password</p> <p style="text-align:center;font-size:14px">https://occitanofòn.xyz/confirmacion/%s</p></body></html>`, translResetPwd.Username, translResetPwd.Token)

		// //Get from each mail
		// email.GetFrom()
		// email.SetBody(mail.TextHTML, body)

		// //Send with high priority
		// email.SetPriority(mail.PriorityHigh)

		// // always check error after send
		// if email.Error != nil {
		// 	es.logger.Error(email.Error)
		// 	continue
		// }

		// if err := email.Send(es.SMTPClient); err != nil {
		// 	es.logger.Error(err)
		// }
	}
}

func (es *Service) SendResetPasswordLink(translResetPwd *entities.TranslatorResetPassword) {
	go func() {
		// notifies that this translator has already tried to reset his password
		es.antiSpammer.notify(translResetPwd.Email)
		es.queue <- *translResetPwd
	}()
}

func (es *Service) IsAllowed(email string) bool {
	es.antiSpammer.mu.RLock()
	defer es.antiSpammer.mu.RUnlock()

	_, has := es.antiSpammer.lastTrySeen[email]

	return !has
}

type antiSpammer struct {
	lastTrySeen map[string]time.Time
	mu          *sync.RWMutex
}

func (as *antiSpammer) removeTryAfterSomeTime() {
	for {
		<-time.After(RETRY_DELAY)

		for email, lastTry := range as.lastTrySeen {
			now := time.Now()
			if now.After(lastTry.Add(RETRY_DELAY)) {
				as.mu.Lock()
				delete(as.lastTrySeen, email)
				as.mu.Unlock()
			}
		}
	}
}

func (as *antiSpammer) notify(email string) {
	as.mu.Lock()
	as.lastTrySeen[email] = time.Now()
	as.mu.Unlock()
}
