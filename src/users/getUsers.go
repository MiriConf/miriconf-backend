package users

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

// GetUsers godoc
// @Summary      Display data for a single user
// @Description  Display data for a single user
// @Tags         teams
// @Produce      json
// @Success      200  {array}   users.User
// @Failure      200  {array}   helpers.Error
// @Router       /user/get/{_id} [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
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

	status, userID := helpers.GetRequestID("user", r, w)
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

	coll := client.Database("miriconf").Collection("users")

	var result GetUser
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: userID}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err != nil {
				error := helpers.ErrorMsg("no user matching id requested")
				w.Write(error)
				helpers.EndpointError("no user matching id requested", r)
				return
			}
		}
		log.Fatal(err)
	}

	helpers.SuccessLog(r)
	json.NewEncoder(w).Encode(result)
	return
}
