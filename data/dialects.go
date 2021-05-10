package data

import "tradoc/models"

// l'ordre est important
var DIALECTS = []models.Dialect{
	{
		Nom: "Auvernhat",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "auv_est",
			},
			{
				Nom:                      "Brivadés",
				TotalSentencesTranslated: 0,
				Abbr:                     "auv_bri",
			},
			{
				Nom:                      "Septentrional",
				TotalSentencesTranslated: 0,
				Abbr:                     "auv_sep",
			},
		},
	},
	{
		Nom: "Gascon",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "gas_est",
			},
			{
				Nom:                      "Aranés",
				TotalSentencesTranslated: 0,
				Abbr:                     "gas_ara",
			},
			{
				Nom:                      "Bearnés",
				TotalSentencesTranslated: 0,
				Abbr:                     "gas_bea",
			},
		},
	},
	{
		Nom: "Lengadocian",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "len_est",
			},
			{
				Nom:                      "Agenés",
				TotalSentencesTranslated: 0,
				Abbr:                     "len_age",
			},
			{
				Nom:                      "Besierenc",
				TotalSentencesTranslated: 0,
				Abbr:                     "len_bes",
			},
			{
				Nom:                      "Carcassés",
				TotalSentencesTranslated: 0,
				Abbr:                     "len_car",
			},
			{
				Nom:                      "Roergat",
				TotalSentencesTranslated: 0,
				Abbr:                     "len_roe",
			},
		},
	},
	{
		Nom: "Lemosin",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "lem_est",
			},
			{
				Nom:                      "Marchés",
				TotalSentencesTranslated: 0,
				Abbr:                     "lem_mar",
			},
			{
				Nom:                      "Peiregordin",
				TotalSentencesTranslated: 0,
				Abbr:                     "lem_pei",
			},
		},
	},
	{
		Nom: "Provençau",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "pro_est",
			},
			{
				Nom:                      "Maritime",
				TotalSentencesTranslated: 0,
				Abbr:                     "pro_mar",
			},
			{
				Nom:                      "Niçard",
				TotalSentencesTranslated: 0,
				Abbr:                     "pro_nic",
			},
			{
				Nom:                      "Rodanenc",
				TotalSentencesTranslated: 0,
				Abbr:                     "pro_rod",
			},
		},
	},
	{
		Nom: "Vivaroaupenc",
		Subdialects: []models.Subdialect{
			{
				Nom:                      "Estandard",
				TotalSentencesTranslated: 0,
				Abbr:                     "viv_est",
			},
			{
				Nom:                      "Aupenc",
				TotalSentencesTranslated: 0,
				Abbr:                     "viv_aup",
			},
			{
				Nom:                      "Gavòt",
				TotalSentencesTranslated: 0,
				Abbr:                     "viv_gav",
			},
			{
				Nom:                      "Vivarodaufinenc",
				TotalSentencesTranslated: 0,
				Abbr:                     "viv_viv",
			},
		},
	},
}
