package systems

type System struct {
	SystemName string   `json:"systemname" bson:"systemname"`
	Users      []string `json:"users" bson:"users"`
	Team       string   `json:"team" bson:"team"`
	LastSeen   int      `json:"lastseen" bson:"lastseen"`
	CreatedAt  string   `json:"createdat" bson:"createdat"`
}
