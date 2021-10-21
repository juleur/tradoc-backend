package mongodb

type collectionDB int

const (
	Datasets collectionDB = iota
	Occitan
	OnGoingTranslations
	SecretQuestions
	TemporaryTokens
	Translations
	TranslationsFiles
	Translators
)

func (c collectionDB) ColletionName() string {
	switch c {
	case Datasets:
		return "Datasets"
	case Occitan:
		return "Occitan"
	case OnGoingTranslations:
		return "OnGoingTranslations"
	case SecretQuestions:
		return "SecretQuestions"
	case TemporaryTokens:
		return "TemporaryTokens"
	case Translations:
		return "Translations"
	case TranslationsFiles:
		return "TranslationsFiles"
	case Translators:
		return "Translators"
	}
	return "unknown"
}
