package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListTeams godoc
// @Summary      List all teams
// @Description  list all teams
// @Tags         teams
// @Produce      json
// @Success      200  {array}   teams.Team
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /teams/list [get]
func ListTeams(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

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
	filter := bson.D{{"id", bson.D{{"$lte", 500}}}}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		if output != nil {
			fmt.Println("successful call on /api/v1/teams/list")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
