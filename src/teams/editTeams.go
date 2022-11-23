package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EditTeams(w http.ResponseWriter, r *http.Request) {
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

	status, teamID := helpers.GetRequestID("team", r, w)
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

	var putTeam Team
	json.Unmarshal(reqBody, &putTeam)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("miriconf").Collection("teams")

	createdAt := time.Now()

	doc := bson.D{{"$set", bson.D{{Key: "apps", Value: []string(putTeam.Apps)}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}}}
	filter := bson.D{{Key: "_id", Value: teamID}}

	result, err := coll.UpdateOne(context.TODO(), filter, doc)
	if err != nil {
		panic(err)
	}

	documentID := fmt.Sprintf("%v", result.UpsertedID)
	success := helpers.SuccessMsg(putTeam.Name + " successfully updated at " + documentID)
	w.Write(success)
	helpers.SuccessLog(r)
	return
}
