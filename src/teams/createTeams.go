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
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTeams godoc
// @Summary      Create a single team
// @Description  Create a single team
// @Tags         teams
// @Produce      json
// @Success      200  string   successRes
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /teams [post]
func CreateTeams(w http.ResponseWriter, r *http.Request) {
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

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error := helpers.ErrorMsg("error reading request body")
		w.Write(error)
		helpers.EndpointError("error reading request body", r)
		return
	}

	var postTeam Team
	json.Unmarshal(reqBody, &postTeam)

	if strings.TrimSpace(postTeam.Name) == "" {
		error := helpers.ErrorMsg("no team name in request")
		w.Write(error)
		helpers.EndpointError("no team name in request", r)
		return
	}

	if strings.TrimSpace(postTeam.Department) == "" {
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
	err = coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: postTeam.Name}}).Decode(&nameCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("team with this name already exists")
		w.Write(error)
		helpers.EndpointError("team with this name already exists", r)
		return
	}

	createdAt := time.Now()

	doc := bson.D{{Key: "name", Value: strings.TrimSpace(postTeam.Name)}, {Key: "department", Value: strings.TrimSpace(postTeam.Department)}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	documentID := fmt.Sprintf("%v", result.InsertedID)
	success := helpers.SuccessMsg(postTeam.Name + " successfully added at " + documentID)
	w.Write(success)
	helpers.SuccessLog(r)
	return
}
