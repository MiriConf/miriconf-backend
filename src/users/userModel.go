package users

type User struct {
	Username  string `json:"username" bson:"username"`
	Fullname  string `json:"fullname" bson:"fullname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	CreatedAt string `json:"createdat" bson:"createdat"`
}
