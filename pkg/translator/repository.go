package translator

import (
	"btradoc/entities"
	"btradoc/helpers"
	"btradoc/pkg"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	GetTranslatorByUsername(username string) (*entities.Translator, error)
	InsertNewTranslator(newTranslator entities.NewTranslator) error
	InsertNewRefreshToken(translatorID string) (string, error)
	GetRefreshToken(refreshToken string) (*entities.Translator, error)
	RemoveRefreshToken(refreshToken string) error
	GetSecretQuestionsByToken(token string) (*entities.TranslatorSecretQuestions, error)
	GetSecretQuestions() ([]string, error)
	CreateTokenResetPassword(email string) (*entities.TranslatorResetPassword, error)
	UpdatePassword(translatorID string, hashedPassword string) error
}

type repository struct {
	MongoDB *mongo.Database
}

func NewRepo(mongoDB *mongo.Database) Repository {
	return &repository{
		MongoDB: mongoDB,
	}
}

func (r *repository) GetTranslatorByUsername(username string) (*entities.Translator, error) {
	translatorsColl := r.MongoDB.Collection("Translators")

	opts := options.FindOne().SetSort(bson.D{{Key: "username", Value: 1}})
	var result bson.M
	err := translatorsColl.FindOne(context.Background(), bson.D{{Key: "username", Value: username}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrBadCredentials,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ID := result["_id"].(primitive.ObjectID)

	translator := entities.Translator{
		ID:        ID.Hex(),
		Email:     result["email"].(string),
		Username:  result["username"].(string),
		Hpwd:      result["hpwd"].(string),
		Confirmed: result["confirmed"].(bool),
		Suspended: result["suspended"].(bool),
	}

	permissions := result["permissions"].(primitive.A)
	for _, permission := range permissions {
		perm := permission.(string)
		translator.Permissions = append(translator.Permissions, perm)
	}

	return &translator, nil
}

func (r *repository) InsertNewTranslator(newTranslator entities.NewTranslator) error {
	translatorsColl := r.MongoDB.Collection("Translators")

	doc := bson.D{
		{Key: "username", Value: newTranslator.Username},
		{Key: "email", Value: newTranslator.Email},
		{Key: "hpwd", Value: newTranslator.Hpwd},
		{Key: "confirmed", Value: false},
		{Key: "suspended", Value: false},
		{Key: "secretQuestions", Value: []bson.D{
			{{Key: "question", Value: newTranslator.SecretQuestions[0].Question}, {Key: "response", Value: newTranslator.SecretQuestions[0].Response}},
			{{Key: "question", Value: newTranslator.SecretQuestions[1].Question}, {Key: "response", Value: newTranslator.SecretQuestions[1].Response}},
		}},
		{Key: "permissions", Value: []string{}},
		{Key: "createdAt", Value: time.Now()},
	}

	if _, err := translatorsColl.InsertOne(context.Background(), doc); err != nil {
		if ok := mongo.IsDuplicateKeyError(err); ok {
			if strings.Contains(err.Error(), "username_unique") {
				return &pkg.DBError{
					Code:    409,
					Message: pkg.ErrPseudoUsed,
					Wrapped: err,
				}
			} else if strings.Contains(err.Error(), "email_unique") {
				return &pkg.DBError{
					Code:    409,
					Message: pkg.ErrEmailUsed,
					Wrapped: err,
				}
			}
		}
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}

func (r *repository) InsertNewRefreshToken(translatorID string) (string, error) {
	// retry if refreshToken generated is already set
	for i := 0; i < 50; i++ {
		refreshToken := helpers.GenerateID(12)

		id, err := primitive.ObjectIDFromHex(translatorID)
		if err != nil {
			return "", &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		translatorsColl := r.MongoDB.Collection("Translators")
		filter := bson.M{"_id": id}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "refreshToken", Value: refreshToken}}}}
		_, err = translatorsColl.UpdateOne(context.Background(), filter, update)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}

			return "", &pkg.DBError{
				Code:    500,
				Message: pkg.ErrDefault,
				Wrapped: err,
			}
		}

		return refreshToken, nil
	}

	return "", &pkg.DBError{
		Code:    500,
		Message: pkg.ErrDefault,
		Wrapped: errors.New("cannot update refresh token"),
	}
}

func (r *repository) GetRefreshToken(refreshToken string) (*entities.Translator, error) {
	translatorsColl := r.MongoDB.Collection("Translators")

	opts := options.FindOne().SetSort(bson.D{{Key: "refreshToken", Value: 1}})
	var result bson.M
	err := translatorsColl.FindOne(context.Background(), bson.D{{Key: "refreshToken", Value: refreshToken}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrRefreshTokenNotFound,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ID := result["_id"].(primitive.ObjectID)

	translator := entities.Translator{
		ID: ID.Hex(),
	}

	permissions := result["permissions"].(primitive.A)
	for _, permission := range permissions {
		perm := permission.(string)
		translator.Permissions = append(translator.Permissions, perm)
	}

	return &translator, nil
}

func (r *repository) RemoveRefreshToken(refreshToken string) error {
	translatorsColl := r.MongoDB.Collection("Translators")

	if _, err := translatorsColl.UpdateOne(context.Background(), bson.M{"refreshToken": refreshToken}, bson.M{"$set": bson.M{"refreshToken": nil}}); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}

func (r *repository) GetSecretQuestionsByToken(token string) (*entities.TranslatorSecretQuestions, error) {
	temporaryTokensColl := r.MongoDB.Collection("TemporaryTokens")

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "token", Value: token}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "Translators"}, {Key: "localField", Value: "translator"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "translator"}}}}
	project1Stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "issuedAt", Value: "$issuedAt"}, {Key: "secretQuestions", Value: "$translator.secretQuestions"}}}}
	project2Stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "issuedAt", Value: "$issuedAt"}, {Key: "secretQuestions", Value: bson.D{{Key: "$first", Value: "$secretQuestions"}}}}}}

	ctx := context.Background()
	cursor, err := temporaryTokensColl.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, project1Stage, project2Stage})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	} else if len(result) == 0 {
		return nil, &pkg.DBError{
			Code:    404,
			Message: pkg.ErrSecretQuestionsNotFound,
			Wrapped: err,
		}
	}

	ID := result[0]["_id"].(primitive.ObjectID)
	timestamp := result[0]["issuedAt"].(primitive.DateTime)

	if time.Now().After(timestamp.Time().Add(12 * time.Hour)) {
		return nil, &pkg.DBError{
			Code:    403,
			Message: pkg.ErrResetPasswordTokenExpired,
			Wrapped: errors.New("reset token password is expired"),
		}
	}

	mapSecretQuestions := make(map[string]string, len(result[0]["secretQuestions"].(primitive.A)))
	for _, mapElems := range result[0]["secretQuestions"].(primitive.A) {
		maps := mapElems.(primitive.M)
		question := maps["question"].(string)
		response := maps["response"].(string)
		mapSecretQuestions[question] = response
	}

	translatorSecretQuestions := entities.TranslatorSecretQuestions{
		ID:              ID,
		SecretQuestions: mapSecretQuestions,
	}
	return &translatorSecretQuestions, nil
}

func (r *repository) GetSecretQuestions() ([]string, error) {
	secretQuestionsColl := r.MongoDB.Collection("SecretQuestions")
	ctx := context.Background()

	cursor, err := secretQuestionsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}
	defer cursor.Close(ctx)

	var secretQuestionsDocs []bson.M
	if err = cursor.All(ctx, &secretQuestionsDocs); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	var secretQuestions []string
	for _, value := range secretQuestionsDocs {
		sq := value["question"].(string)

		secretQuestions = append(secretQuestions, sq)
	}

	return secretQuestions, nil
}

func (r *repository) CreateTokenResetPassword(email string) (*entities.TranslatorResetPassword, error) {
	translatorColl := r.MongoDB.Collection("Translators")
	temporaryTokensColl := r.MongoDB.Collection("TemporaryTokens")

	opts := options.FindOne().SetSort(bson.D{{Key: "email", Value: 1}})
	var result bson.M
	err := translatorColl.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &pkg.DBError{
				Code:    404,
				Message: pkg.ErrTranslatorNotFound,
				Wrapped: err,
			}
		}
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	ID := result["_id"].(primitive.ObjectID)
	usernameRes := result["username"].(string)

	id, err := primitive.ObjectIDFromHex(ID.Hex())
	if err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	token := helpers.GenerateID(20)
	doc := bson.D{
		{Key: "type", Value: "PASSWORD_RESET"},
		{Key: "token", Value: token},
		{Key: "translator", Value: id},
		{Key: "issuedAt", Value: time.Now()},
	}
	if _, err = temporaryTokensColl.InsertOne(context.Background(), doc); err != nil {
		return nil, &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	translatorResetPassword := entities.TranslatorResetPassword{
		Email:    email,
		Username: usernameRes,
		Token:    token,
	}

	return &translatorResetPassword, nil
}

func (r *repository) UpdatePassword(translatorID string, hashedPassword string) error {
	translatorColl := r.MongoDB.Collection("Translators")

	translatorObjectID, err := primitive.ObjectIDFromHex(translatorID)
	if err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	if _, err = translatorColl.UpdateOne(context.Background(), bson.M{"_id": translatorObjectID}, bson.D{{Key: "$set", Value: bson.D{{Key: "hpwd", Value: hashedPassword}}}}); err != nil {
		return &pkg.DBError{
			Code:    500,
			Message: pkg.ErrDefault,
			Wrapped: err,
		}
	}

	return nil
}
