package teams

type Team struct {
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
	CreatedAt  string `json:"createdat" bson:"createdat"`
}
