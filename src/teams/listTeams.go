package teams

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
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
// @Router       /teams/list [get]
func ListTeams(w http.ResponseWriter, r *http.Request) {
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

	cursor, err := coll.Find(context.TODO(), bson.D{}, options.Find().SetLimit(10))
	if err != nil {
		panic(err)
	}

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	helpers.SuccessLog(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
