package token

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
	"tradoc/message"
	"tradoc/models"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

const regex = `^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`

var rxJWT = regexp.MustCompile(regex)

func IsItAJwtToken(jwt []byte) bool {
	return rxJWT.Match(jwt)
}

func GetJwtPayload(SECRET_KEY []byte, jwToken []byte) (int, bool) {
	pl := models.JWTPayload{}

	signature := jwt.NewHS256(SECRET_KEY)
	expValidator := jwt.ExpirationTimeValidator(time.Now())
	validatePayload := jwt.ValidatePayload(&pl.Payload, expValidator)
	if _, err := jwt.Verify(jwToken, signature, &pl, validatePayload); err != nil {
		if err == jwt.ErrExpValidation {
			return 0, true
		}
		return 0, false
	}
	return pl.TranslatorID, false
}

func GenerateTokens(SECRET_KEY []byte, translator models.Translator) (*models.Tokens, *models.Log) {
	pl := models.JWTPayload{
		Payload: jwt.Payload{
			Issuer:         "https://backend.trad-oc.xyz",
			ExpirationTime: jwt.NumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:       jwt.NumericDate(time.Now()),
		},
		TranslatorID: translator.ID,
		Username:     translator.Username,
		Dialects:     translator.Dialects,
	}
	jwtToken, err := jwt.Sign(&pl, jwt.NewHS256(SECRET_KEY))
	if err != nil {
		return nil, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrUnauthorized.Code,
			Message:        message.ResponseAuthFailed,
		}
	}
	tokens := models.Tokens{}
	tokens.JWT = string(jwtToken)
	tokens.RefreshToken = HexKeyGenerator(32)
	return &tokens, nil
}

func HexKeyGenerator(nb int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const letterBytes = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.Grow(nb)
	for ; nb > 0; nb-- {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes)-1)])
	}
	return sb.String()
}

func IsAlphanumeric(token []byte) bool {
	for _, c := range token {
		if (c >= 48 && c <= 57) || (c >= 65 && c <= 90) || (c >= 97 && c <= 122) {
			continue
		} else {
			return false
		}
	}
	return true
}
