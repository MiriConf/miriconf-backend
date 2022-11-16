package teams

type Team struct {
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
	SourceRepo string `json:"source_repo" bson:"source_repo"`
	SourcePAT  string `json:"source_pat" bson:"source_pat"`
	CreatedAt  string `json:"createdat" bson:"createdat"`
}
