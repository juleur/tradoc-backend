package entities

type Dataset struct {
	ID       string `json:"_id" bson:"_id"`
	Sentence string `json:"sentence" bson:"sentence"`
}
