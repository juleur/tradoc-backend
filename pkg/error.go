package pkg

const (
	ErrDefault                   string = "Une erreur est survenue sur le serveur"
	ErrSecretQuestionsNotFound   string = "impossible de trouver les questions secretes"
	ErrBadCredentials            string = "mauvais crédentials"
	ErrEmailUsed                 string = "Cette email est déjà utilisé"
	ErrPseudoUsed                string = "Ce pseudonyme est déjà utilisé"
	ErrRefreshTokenNotFound      string = "Une erreur est survenue lors de la réauthentification de votre compte"
	ErrResetPasswordTokenExpired string = "Le token de remise à 0 du mot de passe a expiré"
	ErrTranslatorNotFound        string = "Impossible de trouver ce traducteur"
	ErrDialectNotFound           string = "Impossible de trouver ce dialect"
	ErrDialectsNotFound          string = "Impossible de trouver les dialectes"
)

type DBError struct {
	Code    int
	Message string
	Wrapped error
}

func (e DBError) Error() string {
	return e.Wrapped.Error()
}
