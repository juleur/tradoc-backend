package controllers

const (
	ErrDefault                    string = "Une erreur est survenue sur le serveur"
	ErrAccountNotConfirmed        string = "Votre compte a besoin d'être activer afin de pouvoir vous connecter"
	ErrAccountSuspended           string = "Votre compte est suspendu, contactez l'administrateur"
	ErrDialectNotProvided         string = "Aucun dialect n'a été trouvé"
	ErrTooMuchTranslationsFetched string = "tu as déja traduit plus de 300 phrases pour aujourd'hui, reviens demain"
	ErrNoMoreDataset              string = "Plus aucune phrases à traduire"
	ErrPasswordTooShort           string = "Ton mot de passe doit contenir au moins 10 caractères"
	ErrSecretQuestions            string = "Tu n'as pas saissi les 2 secret questions"
	ErrMailerNotAllowed           string = "Tu as déjà récemment essayé de changer ton mot de passe, ressaye ultérieurement"
	ErrSecretQuestionsNoMatch     string = "Une erreur est survenue avec les questions secrètes"
	ErrBadCredentials             string = "Mauvais credentials"
	ErrNoPermDialect              string = "Aucune dialect ne t'a été attribué, contacte l'administrateur"
	ErrBadFullDialectFormat       string = "Une erreur avec le dialect saissi"
)
