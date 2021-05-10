package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"tradoc/message"
	"tradoc/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/phuslu/log"
)

func (dp *DBPsql) FindTranslator(loginBody models.LoginBody) (models.Translator, *models.Log) {
	translator := models.Translator{}
	pwd := &sql.NullString{}
	if err := dp.db.QueryRow(context.Background(), `
		SELECT id, pseudonim, senhal, suspension FROM traductors
		WHERE pseudonim = $1 LIMIT 1;
	`, loginBody.Username).Scan(&translator.ID, &translator.Username, pwd, &translator.Suspension); err != nil {
		if err == pgx.ErrNoRows {
			return models.Translator{}, &models.Log{
				Level:          log.WarnLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrNotFound.Code,
				Message:        message.ResponseNotRegistered,
			}
		}
		return models.Translator{}, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}

	if translator.Suspension {
		return models.Translator{}, &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrSuspendedAccount,
			HttpStatusCode: fiber.ErrForbidden.Code,
			Message:        message.ResponseAccountSuspended,
		}
	}
	if pwd.String == "" {
		return models.Translator{}, &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrAccountNotValidated,
			HttpStatusCode: fiber.ErrConflict.Code,
			Message:        message.ResponseAccountNotValidated,
		}
	}
	translator.HashedPassword = pwd.String

	return translator, nil
}

func (dp *DBPsql) UpdateRefreshToken(translatorID int, refreshToken string) *models.Log {
	if _, err := dp.db.Exec(context.Background(), `
		UPDATE traductors SET refresh_token = $1 WHERE id = $2
	`, refreshToken, translatorID); err != nil {
		return &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	return nil
}

func (dp *DBPsql) FetchDialectsByTranslator(translatorID int) ([]string, *models.Log) {
	rows, err := dp.db.Query(context.Background(), "SELECT * FROM dialectes_auth WHERE traductor_id = $1 LIMIT 1", translatorID)
	if err != nil {
		return []string{}, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}

	// récupére tous les dialectes
	dialectsByAbbr := []string{}
	for i, field := range rows.FieldDescriptions() {
		// n'ajoute pas traducteur_id
		if i >= len(rows.FieldDescriptions())-1 {
			break
		}
		// split le dialect et le sous-dialect afin d'avoir les abbréviations
		dial_split := strings.Split(string(field.Name), "_")
		dialectsByAbbr = append(dialectsByAbbr, dial_split[0][0:3]+"_"+dial_split[1][0:3])
	}

	// compare les valeurs avec les dialectes récupérés plus haut
	translatorDialectsByAbbr := []string{}
	for rows.Next() {
		rowValues, err := rows.Values()
		if err != nil {
			return []string{}, &models.Log{
				Level:          log.ErrorLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrServiceUnavailable.Code,
				Message:        message.ResponseErrServer,
			}
		}
		for i, value := range rowValues {
			// s'arrête avant la colonne traductor_id
			if i > len(rowValues)-1 {
				break
			}
			// compare les 2 tableaux qui sont triés dans le même ordre
			// ["auvernhat_estandard", "gascon_estandard", "lengadocian_estandard", ...]
			// [true, false, false, ...]
			if value == true {
				translatorDialectsByAbbr = append(translatorDialectsByAbbr, dialectsByAbbr[i])
			}
		}
	}
	rows.Close()

	return translatorDialectsByAbbr, nil
}

func (dp *DBPsql) FindAllDialect() ([]string, *models.Log) {
	rows, err := dp.db.Query(context.Background(), "SELECT * FROM dialectes_auth WHERE false;")
	if err != nil {
		return []string{}, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	defer rows.Close()
	// récupére tous les dialectes
	dialects := []string{}
	for i, field := range rows.FieldDescriptions() {
		// n'ajoute pas traducteur_id
		if i >= len(rows.FieldDescriptions())-1 {
			break
		}
		dialects = append(dialects, string(field.Name))
	}
	return dialects, nil
}

func (dp *DBPsql) FetchTotalSentenceTranslatedByDialect(dialectTableName string) (int, *models.Log) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s;", dialectTableName)
	var total int
	if err := dp.db.QueryRow(context.Background(), query).Scan(&total); err != nil {
		return total, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	return total, nil
}

func (dp *DBPsql) FetchAllSentencesByDialect(dialectTableName string, translatorID int, ruleOutTexteIDs []int) ([]models.Texte, *models.Log) {
	withQuery := ruleOutTextesIdswithQueryGenerator(ruleOutTexteIDs)
	query := fmt.Sprintf(`
		%s, occitan_table AS (
			SELECT DISTINCT id FROM %s
		), merge_ids AS (
			SELECT * FROM rule_out_IDs
			UNION
			SELECT * FROM occitan_table
		)
		SELECT id, frasa
		FROM textes_oc
		WHERE
			id NOT IN (SELECT * FROM merge_ids)
			AND
			(SELECT %s FROM dialectes_auth WHERE traductor_id = $1)
		OFFSET floor(random() * (SELECT COUNT(*) FROM textes_oc))
		LIMIT 10;
	`, withQuery, dialectTableName, dialectTableName)

	rows, err := dp.db.Query(context.Background(), query, translatorID)
	if err != nil {
		return []models.Texte{}, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}

	textes := []models.Texte{}
	for rows.Next() {
		texte := models.Texte{}
		err := rows.Scan(&texte.ID, &texte.Sentence)
		if err != nil {
			return []models.Texte{}, &models.Log{
				Level:          log.ErrorLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrServiceUnavailable.Code,
				Message:        message.ResponseErrServer,
			}
		}
		textes = append(textes, texte)
	}
	rows.Close()

	if len(textes) == 0 {
		return []models.Texte{}, &models.Log{
			Level:          log.WarnLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrNotFound.Code,
			Message:        message.ResponseNoSentenceToTranslate,
		}
	}
	return textes, nil
}

func ruleOutTextesIdswithQueryGenerator(ruleOutTexteIDs []int) string {
	if len(ruleOutTexteIDs) == 0 {
		query := "WITH rule_out_IDs AS (VALUES (0))"
		return query
	}
	query := "WITH rule_out_IDs AS (VALUES "
	for i, id := range ruleOutTexteIDs {
		query += fmt.Sprintf("(%d)", id)
		if i < len(ruleOutTexteIDs)-1 {
			query += ","
		}
	}
	query += ")"

	return query
}

func (dp *DBPsql) AddTranslation(translation models.Translation, dialectTableName string, translatorID int) *models.Log {
	var query string
	if translation.Default.English == nil {
		query = fmt.Sprintf("INSERT INTO %s(frasa_%s, frasa_fr, traductor_id, texte_oc_id) VALUES ($1,$2,$3,$4)", dialectTableName, translation.Abbr)
		if _, err := dp.db.Exec(context.Background(), query, translation.Default.Occitan, translation.Default.French, translatorID, translation.ID); err != nil {
			return &models.Log{
				Level:          log.ErrorLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrServiceUnavailable.Code,
				Message:        message.ResponseErrServer,
			}
		}
	} else {
		query = fmt.Sprintf("INSERT INTO %s(frasa_%s, frasa_fr, frasa_an, traductor_id, texte_oc_id) VALUES ($1,$2,$3,$4,$5)", dialectTableName, translation.Abbr)
		if _, err := dp.db.Exec(context.Background(), query, translation.Default.Occitan, translation.Default.French, *translation.Default.English, translatorID, translation.ID); err != nil {
			return &models.Log{
				Level:          log.ErrorLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrServiceUnavailable.Code,
				Message:        message.ResponseErrServer,
			}
		}
	}
	if translation.Feminize != nil {
		if translation.Feminize.English == nil {
			query = fmt.Sprintf("INSERT INTO %s(frasa_%s, frasa_fr, traductor_id, texte_oc_id) VALUES ($1,$2,$3,$4)", dialectTableName, translation.Abbr)
			if _, err := dp.db.Exec(context.Background(), query, translation.Feminize.Occitan, translation.Feminize.French, translatorID, translation.ID); err != nil {
				return &models.Log{
					Level:          log.ErrorLevel,
					Error:          err,
					HttpStatusCode: fiber.ErrServiceUnavailable.Code,
					Message:        message.ResponseErrServer,
				}
			}
		} else {
			query = fmt.Sprintf("INSERT INTO %s(frasa_%s, frasa_fr, frasa_an, traductor_id, texte_oc_id) VALUES ($1,$2,$3,$4,$5)", dialectTableName, translation.Abbr)
			if _, err := dp.db.Exec(context.Background(), query, translation.Feminize.Occitan, translation.Feminize.French, *translation.Feminize.English, translatorID, translation.ID); err != nil {
				return &models.Log{
					Level:          log.ErrorLevel,
					Error:          err,
					HttpStatusCode: fiber.ErrServiceUnavailable.Code,
					Message:        message.ResponseErrServer,
				}
			}
		}
	}
	return nil
}

func (dp *DBPsql) IsRefreshTokenExists(userID int, refreshToken string) (models.Translator, *models.Log) {
	translator := models.Translator{}
	if err := dp.db.QueryRow(context.Background(), `
		SELECT id, pseudonim, refresh_token FROM traductors WHERE id = $1 LIMIT 1;
	`, userID).Scan(&translator.ID, &translator.Username, &translator.RefreshToken); err != nil {
		if err == pgx.ErrNoRows {
			return translator, &models.Log{
				Level:          log.WarnLevel,
				Error:          err,
				HttpStatusCode: fiber.ErrUnauthorized.Code,
				Message:        message.ResponseAuthFailed,
			}
		}
		return translator, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	if translator.RefreshToken != refreshToken {
		return translator, &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrInvalidRefreshToken,
			HttpStatusCode: fiber.ErrUnauthorized.Code,
			Message:        message.ResponseAuthFailed,
		}
	}
	return translator, nil
}

func (dp *DBPsql) FirstRegister(username string, code int) (int, *models.Log) {
	var translatorID int
	var activationCode int
	password := &sql.NullString{}
	creatAt := &sql.NullTime{}
	if err := dp.db.QueryRow(context.Background(), `
		SELECT id, senhal, code_activation, creat_lo FROM traductors
		WHERE pseudonim = $1 AND code_activation = $2;
	`, username, code).Scan(&translatorID, password, &activationCode, creatAt); err != nil {
		if err == pgx.ErrNoRows {
			return 0, &models.Log{
				Level:          log.WarnLevel,
				Error:          message.ErrTranslatorAccountNotFound,
				HttpStatusCode: fiber.ErrForbidden.Code,
				Message:        message.ResponseTranslatorAccountNotFound,
			}
		}
		return 0, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	if password.Valid && creatAt.Valid {
		return 0, &models.Log{
			Level:          log.WarnLevel,
			Error:          message.ErrAccountAlreadyValidated,
			HttpStatusCode: fiber.ErrForbidden.Code,
			Message:        message.ResponseAccountAlreadyValidated,
		}
	}
	return translatorID, nil
}

func (dp *DBPsql) CreateNewTranslator(translatorID int, username string, hashedPWd string) *models.Log {
	if _, err := dp.db.Exec(context.Background(), `
		UPDATE traductors SET senhal = $1, creat_lo = $2 WHERE id = $3 AND pseudonim = $4
	`, hashedPWd, time.Now(), translatorID, username); err != nil {
		return &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	return nil
}

func (dp *DBPsql) FetchTranslatedFiles(dialectTableName string) (*string, *string, *models.Log) {
	var filepathFr string
	filepathEn := &sql.NullString{}
	if err := dp.db.QueryRow(context.Background(), `
		SELECT filename_fr, filename_en
		FROM translation_files
		WHERE dialect_name = $1
		ORDER BY creat_lo DESC
		LIMIT 1
	`, dialectTableName).Scan(&filepathFr, filepathEn); err != nil {
		return nil, nil, &models.Log{
			Level:          log.ErrorLevel,
			Error:          err,
			HttpStatusCode: fiber.ErrServiceUnavailable.Code,
			Message:        message.ResponseErrServer,
		}
	}
	if filepathEn.Valid {
		return &filepathFr, &filepathEn.String, nil
	}
	return &filepathFr, nil, nil
}

func (dp *DBPsql) FetchLastGeneratedFile() *time.Time {
	lastUpdate := &sql.NullTime{}
	if err := dp.db.QueryRow(context.Background(), `
		SELECT creat_lo
		FROM translation_files
		ORDER BY creat_lo DESC
		LIMIT 1
	`).Scan(lastUpdate); err != nil {
		return nil
	}
	if lastUpdate.Valid {
		return &lastUpdate.Time
	}
	return nil
}
