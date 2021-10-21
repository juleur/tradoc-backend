package helpers

const (
	ErrDefault                string = "une erreur est survenue sur le serveur"
	ErrOnlyOneSpaceByUsername string = "pas mai que 1 espaci"
	ErrUsernameTooShort       string = "ton pseudo est trop court (3 caractères minimum)"
	ErrEmailTooShort          string = "l'email est trop court"
	ErrEmailTooLong           string = "l'email est trop long"
	ErrEmailBadFormat         string = "ceci ne correspond pas au format d'une email"
	ErrEmailDomainNotExist    string = "ce domaine n'existe pas"
)

type HelperError struct {
	Code    int
	Message string
	Wrapped error
}

func (e HelperError) Error() string {
	return e.Wrapped.Error()
}
