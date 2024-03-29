package systems

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateSystems(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	jwtKey := []byte(os.Getenv("JWT_KEY"))

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

	var systemInfo System
	json.Unmarshal(reqBody, &systemInfo)

	if strings.TrimSpace(systemInfo.SystemName) == "" {
		error := helpers.ErrorMsg("no system name in request")
		w.Write(error)
		helpers.EndpointError("no system name in request", r)
		return
	}

	if len(systemInfo.Users) <= 0 {
		error := helpers.ErrorMsg("no users in request")
		w.Write(error)
		helpers.EndpointError("no users in request", r)
		return
	}

	if strings.TrimSpace(systemInfo.Team) == "" {
		error := helpers.ErrorMsg("no team in request")
		w.Write(error)
		helpers.EndpointError("no team in request", r)
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

	coll := client.Database("miriconf").Collection("systems")

	var nameCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "systemname", Value: systemInfo.SystemName}}).Decode(&nameCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("system with this name already exists")
		w.Write(error)
		helpers.EndpointError("system with this name already exists", r)
		return
	}

	createdAt := time.Now()

	doc := bson.D{{Key: "systemname", Value: strings.TrimSpace(systemInfo.SystemName)}, {Key: "users", Value: []string(systemInfo.Users)}, {Key: "team", Value: systemInfo.Team}, {Key: "lastseen", Value: -1}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	objectID := result.InsertedID.(primitive.ObjectID).Hex()

	var tokenClaim = helpers.Token{
		Username: objectID,
		Hostname: os.Getenv("MIRICONF_HOSTNAME"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(17532 * time.Hour).Unix(),
		},
	}

	agentToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	tokenString, err := agentToken.SignedString(jwtKey)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tokenString)
	helpers.SuccessLog(r)
	return
}
