package users

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Login(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error := helpers.ErrorMsg("error reading request body")
		w.Write(error)
		helpers.EndpointError("error reading request body", r)
		return
	}

	var postLogin LoginData
	json.Unmarshal(reqBody, &postLogin)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("miriconf").Collection("users")

	var result LoginData
	err = coll.FindOne(context.TODO(), bson.D{{Key: "username", Value: postLogin.Username}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err != nil {
				error := helpers.ErrorMsg("password is incorrect")
				w.Write(error)
				helpers.EndpointError("user does not exist", r)
				return
			}
		}
		log.Fatal(err)
	}

	if helpers.CheckPassword(result.Password, postLogin.Password) == false {
		error := helpers.ErrorMsg("password is incorrect")
		w.Write(error)
		helpers.EndpointError("password is incorrect", r)
		return
	}

	var tokenClaim = helpers.Token{
		Username: postLogin.Username,
		Hostname: os.Getenv("MIRICONF_HOSTNAME"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tokenString)
	helpers.SuccessLog(r)
	return
}
