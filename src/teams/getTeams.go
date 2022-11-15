package teams

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTeams(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	w.Header().Set("Content-Type", "application/json")

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
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: teamID}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err != nil {
				error := helpers.ErrorMsg("no team matching id requested")
				w.Write(error)
				helpers.EndpointError("no team matching id requested", r)
				return
			}
		}
		log.Fatal(err)
	}

	helpers.SuccessLog(r)
	json.NewEncoder(w).Encode(result)
	return
}
