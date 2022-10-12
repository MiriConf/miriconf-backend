package teams

type Team struct {
	MongoID    string            `json:"_id" bson:"_id"`
	ID         int               `json:"id" bson:"id"`
	Name       string            `json:"name" bson:"name"`
	Department string            `json:"department" bson:"department"`
	Tags       map[string]string `json:"tags" bson:"tags"`
	CreatedAt  string            `json:"createdat" bson:"createdat"`
}
