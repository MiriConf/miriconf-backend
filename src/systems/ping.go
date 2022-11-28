package systems

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SystemID struct {
	SystemID string `json:"username"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

	headerToken := r.Header.Get("Authorization")
	if headerToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := helpers.ValidateToken(headerToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	systemData := strings.Split(headerToken, ".")

	tokenDec, err := base64.RawStdEncoding.DecodeString(systemData[1])
	if err != nil {
		panic(err)
	}

	var systemId SystemID
	json.Unmarshal(tokenDec, &systemId)

	sysID, err := primitive.ObjectIDFromHex(systemId.SystemID)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("miriconf").Collection("systems")

	var nameCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: sysID}}).Decode(&nameCheck)
	if err != nil {
		error := helpers.ErrorMsg("system does not exist or cannot be found")
		w.Write(error)
		helpers.EndpointError("system does not exist or cannot be found", r)
		fmt.Println(err)
		return
	}

	lastSeen := time.Now().Unix()

	doc := bson.D{{"$set", bson.D{{Key: "lastseen", Value: lastSeen}}}}
	filter := bson.D{{Key: "_id", Value: sysID}}

	result, err := coll.UpdateOne(context.TODO(), filter, doc)
	if err != nil {
		panic(err)
	}

	documentID := fmt.Sprintf("%v", result.UpsertedID)
	success := helpers.SuccessMsg(documentID + "successfully checked in")
	w.Write(success)
	helpers.SuccessLog(r)
	return
}
