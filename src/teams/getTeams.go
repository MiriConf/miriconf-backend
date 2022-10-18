package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/middlewares"
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

	reqParams := mux.Vars(r)
	reqID := reqParams["id"]

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
			fmt.Println("no docs")
			// This error means your query did not match any documents.
			return
		}
		panic(err)
	}
	// end findOne

	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	if output != nil {
		middlewares.CallLog(r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)




	//var results bson.M
	//err = coll.FindOne(context.TODO(), bson.D{{"id", reqID}}).Decode(&results)
	//if err != nil {
	//	if err == mongo.ErrNoDocuments {
	//		// This error means your query did not match any documents.
	//		return
	//	}
	//	panic(err)
	//}
	////debug
	//fmt.Println(results)
	//if results != nil {
	//	middlewares.CallLog(r)
	//}
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(results)
}
