package users

type User struct {
	Username  string `json:"username" bson:"username"`
	Fullname  string `json:"fullname" bson:"fullname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	CreatedAt string `json:"createdat" bson:"createdat"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUser struct {
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdat"`
}
