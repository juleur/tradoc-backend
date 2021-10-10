package helpers

const (
	ErrDefault                string = "Une erreur est survenue sur le serveur"
	ErrOnlyOneSpaceByUsername string = "Pas mai que 1 espaci"
	ErrUsernameTooShort       string = "ton pseudo est trop court (3 caractères minimum)"
)

type HelperError struct {
	Code    int
	Message string
	Wrapped error
}

func (e HelperError) Error() string {
	return e.Wrapped.Error()
}
