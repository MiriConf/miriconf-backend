package systems

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EditSystems(w http.ResponseWriter, r *http.Request) {
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

	status, systemID := helpers.GetRequestID("system", r, w)
	if status == 1 {
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error := helpers.ErrorMsg("error reading request body")
		w.Write(error)
		helpers.EndpointError("error reading request body", r)
		return
	}

	var putSystem System
	json.Unmarshal(reqBody, &putSystem)

	if strings.TrimSpace(putSystem.SystemName) == "" {
		error := helpers.ErrorMsg("no system name in request")
		w.Write(error)
		helpers.EndpointError("no system name in request", r)
		return
	}

	if len(putSystem.Users) <= 0 {
		error := helpers.ErrorMsg("no users in request")
		w.Write(error)
		helpers.EndpointError("no users in request", r)
		return
	}

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
	err = coll.FindOne(context.TODO(), bson.D{{Key: "systemname", Value: putSystem.SystemName}}).Decode(&nameCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("system with this name already exists")
		w.Write(error)
		helpers.EndpointError("system with this name already exists", r)
		return
	}

	createdAt := time.Now()
	lastSeen := time.Now().Unix()

	doc := bson.D{{"$set", bson.D{{Key: "systemname", Value: strings.TrimSpace(putSystem.SystemName)}, {Key: "users", Value: []string(putSystem.Users)}, {Key: "lastseen", Value: lastSeen}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}}}
	filter := bson.D{{Key: "_id", Value: systemID}}

	result, err := coll.UpdateOne(context.TODO(), filter, doc)
	if err != nil {
		panic(err)
	}

	documentID := fmt.Sprintf("%v", result.UpsertedID)
	success := helpers.SuccessMsg(putSystem.SystemName + " successfully updated at " + documentID)
	w.Write(success)
	helpers.SuccessLog(r)
	return
}
