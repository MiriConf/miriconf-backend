package users

import (
	"context"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteUsers godoc
// @Summary      Delete a single user
// @Description  Delete a single user
// @Tags         teams
// @Produce      json
// @Success      200  string   successRes
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /users [delete]
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
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

	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: userID}})
	if err != nil {
		panic(err)
	}

	if result.DeletedCount >= 1 {
		success := helpers.SuccessMsg("user deleted")
		w.Write(success)
		helpers.SuccessLog(r)
		return
	} else {
		error := helpers.ErrorMsg("no user matching id requested")
		w.Write(error)
		helpers.EndpointError("no user matching id requested", r)
		return
	}
}
