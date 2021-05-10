package message

import "errors"

var (
	ErrAccountNotValidated             = errors.New("This account hasn't been validated yet")
	ErrAddOngoingDialectTranslating    = errors.New("AddOngoingDialectTranslating")
	ErrAddTranslation                  = errors.New("Impossible to add new translaiton")
	ErrBadBodyContent                  = errors.New("Body content has bad format")
	ErrComparingPasswords              = errors.New("Imposible to compare hashing password and password")
	ErrDeleteOngoingDialectTranslating = errors.New("DeleteOngoingDialectTranslating")
	ErrHashedPasswordFailed            = errors.New("Hashing password failed")
	ErrJWTExpired                      = errors.New("JWT Expired")
	ErrJWTPayload                      = errors.New("JWT Payload error")
	ErrNoPasswords                     = errors.New("No hashing password or password")
	ErrNoDialectAbbrFound              = errors.New("Impossible to found a dialect")
	ErrNoDialectAbbrProvided           = errors.New("No dialect abbreviation provided")
	ErrNoSentencesFound                = errors.New("No sentence to translate in database")
	ErrNullTimestamp                   = errors.New("Timestamp is null")
	ErrPasswordNoMatch                 = errors.New("Typed password didn't match")
	ErrRefreshTokenBadFormat           = errors.New("Refresh Token has bad format")
	ErrInvalidRefreshToken             = errors.New("Invalid Refresh Token")
	ErrTranslationLimitRate            = errors.New("Already more than 300 sentences translated")
	ErrTranslatorAccountNotFound       = errors.New("Impossible to find this translator")
	ErrAccountAlreadyValidated         = errors.New("Already validated")
	ErrIncorrectActivationCode         = errors.New("Incorrect activation code")
	ErrWrongBearerJWT                  = errors.New("Wrong JWT Bearer")
	ErrSuspendedAccount                = errors.New("account is suspended")
)
