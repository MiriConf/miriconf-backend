package teams

type Team struct {
	Name       string   `json:"name" bson:"name"`
	Department string   `json:"department" bson:"department"`
	SourceRepo string   `json:"source_repo" bson:"source_repo"`
	SourcePAT  string   `json:"source_pat" bson:"source_pat"`
	Apps       []string `json:"apps" bson:"apps"`
	CreatedAt  string   `json:"createdat" bson:"createdat"`
}

type GetTeam struct {
	Name       string   `json:"name" bson:"name"`
	Department string   `json:"department" bson:"department"`
	SourceRepo string   `json:"source_repo" bson:"source_repo"`
	Apps       []string `json:"apps" bson:"apps"`
	CreatedAt  string   `json:"createdat" bson:"createdat"`
}
