package db

import (
	"testing"
	"tradoc/models"
	"tradoc/utils"
)

func TestFindTranslator(t *testing.T) {
	db := OpenDB()
	dp := DBPsql{db: db}
	loginBody := models.LoginBody{
		Username: "oooooo",
		Password: "test",
	}
	translator, err := dp.FindTranslator(loginBody)
	if err != nil {
		t.Error(err.Message)
	}
	t.Logf("%+v\n", translator)
}

func TestDialectsFetcher(t *testing.T) {
	db := OpenDB()
	dp := DBPsql{db: db}
	translaterID := 1
	dialects, err := dp.FetchDialectsByTranslator(translaterID)
	if err != nil {
		t.Error(err.Message)
	}
	t.Log(dialects)
}

func TestAddTranslation(t *testing.T) {
	db := OpenDB()
	dp := DBPsql{db: db}
	dialectsBrut, err := dp.FindAllDialect()
	if err != nil {
		t.Error(err)
	}
	english := "english english english english"
	// translation := models.Translation{
	// 	ID:   8,
	// 	Abbr: "len_est",
	// 	Default: models.Default{
	// 		Occitan: "La navèra crubada qu’es donc arribada",
	// 		French:  "languedocien Estandard languedocien Estandard",
	// 	},
	// }
	// translation := models.Translation{
	// 	ID:   11,
	// 	Abbr: "len_est",
	// 	Default: models.Default{
	// 		Occitan: "Era Val d'Aran que reclame investiments en carretères entà melhorar-ne era seguretat",
	// 		French:  "languedocien Estandard languedocien Estandard",
	// 		English: &english,
	// 	},
	// }
	// translation := models.Translation{
	// 	ID:   17,
	// 	Abbr: "gas_ara",
	// 	Default: models.Default{
	// 		Occitan: "Lo rampèl serviguèt pas gaire, pasmens: doas minutas après, Alex Costa balhèt la victòria als parisencs",
	// 		French:  "gascon Arane gascon Arane gascon Arane gascon Arane",
	// 	},
	// 	Feminize: &models.Feminize{
	// 		Occitan: "Lo rampèl serviguèt feminize pas gaire, pasmens: doas minutas après, Alex Costa balhèt la victòria als parisencs",
	// 		French:  "gascon Arane gascon feminize",
	// 	},
	// }
	translation := models.Translation{
		ID:   15,
		Abbr: "gas_ara",
		Default: models.Default{
			Occitan: "Montpelhièr seguís sul meteis ritmeDins la correguda al títol Montpelhièr e París, tots dos vencedors aquela dimenjada, daissan pas l'adversari alenar",
			French:  "gascon Arane gascon Arane gascon Arane gascon Arane",
			English: &english,
		},
		Feminize: &models.Feminize{
			Occitan: "Montpelhièr seguís sul meteis ritmeDins la correguda al títol Montpelhièr e París, tots dos vencedors aquela dimenjada, daissan pas l'adversari alenar",
			French:  "gascon Arane gascon feminize",
			English: &english,
		},
	}
	dialectTableName, err := utils.FindDialectByAbbre(dialectsBrut, translation.Abbr)
	if err != nil {
		t.Error(err)
	}
	userID := 1
	err = dp.AddTranslation(translation, dialectTableName, userID)
	if err != nil {
		t.Error(err)
	}
}
