package teams

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EditTeams(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

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

	if strings.TrimSpace(putTeam.Name) == "" {
		error := helpers.ErrorMsg("no team name in request")
		w.Write(error)
		helpers.EndpointError("no team name in request", r)
		return
	}

	if strings.TrimSpace(putTeam.Department) == "" {
		error := helpers.ErrorMsg("no department name in request")
		w.Write(error)
		helpers.EndpointError("no department name in request", r)
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

	coll := client.Database("miriconf").Collection("teams")

	var nameCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: putTeam.Name}}).Decode(&nameCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("team with this name already exists")
		w.Write(error)
		helpers.EndpointError("team with this name already exists", r)
		return
	}

	createdAt := time.Now()

	doc := bson.D{{"$set", bson.D{{Key: "name", Value: strings.TrimSpace(putTeam.Name)}, {Key: "department", Value: strings.TrimSpace(putTeam.Department)}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}}}
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
