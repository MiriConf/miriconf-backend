package teams

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/MiriConf/miriconf-backend/helpers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTeams godoc
// @Summary      Display data for a single team
// @Description  Display data for a single team
// @Tags         teams
// @Produce      json
// @Success      200  {array}   teams.Team
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /teams/get/{id} [get]
func GetTeams(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

	requestParams := mux.Vars(r)
	requestIDRaw := requestParams["id"]
	reqID, _ := strconv.ParseInt(requestIDRaw, 10, 64)

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

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"id", reqID}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			testerr := `{"error":"no teams matching id requested"}`
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(testerr)

			helpers.EndpointError("no teams matching id", r)
			return
		}
		panic(err)
	}

	if result != nil {
		helpers.SuccessLog(r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
