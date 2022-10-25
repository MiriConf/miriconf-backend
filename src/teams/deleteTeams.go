package teams

import (
	"context"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteTeams godoc
// @Summary      Delete a single team
// @Description  Delete a single team
// @Tags         teams
// @Produce      json
// @Success      200  string   successRes
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /teams [delete]
func DeleteTeams(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	w.Header().Set("Content-Type", "application/json")

	status, teamID := helpers.GetRequestID("team", r, w)
	if status == 1 {
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

	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: teamID}})
	if err != nil {
		panic(err)
	}

	if result.DeletedCount >= 1 {
		success := helpers.SuccessMsg("team deleted")
		w.Write(success)
		helpers.SuccessLog(r)
		return
	} else {
		error := helpers.ErrorMsg("no team matching id requested")
		w.Write(error)
		helpers.EndpointError("no team matching id requested", r)
		return
	}
}
