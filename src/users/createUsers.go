package users

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUsers godoc
// @Summary      Create a single user
// @Description  Create a single user
// @Tags         users
// @Produce      json
// @Success      200  string   successRes
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /users [post]
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		error := helpers.ErrorMsg("error reading request body")
		w.Write(error)
		helpers.EndpointError("error reading request body", r)
		return
	}

	var postUser User
	json.Unmarshal(reqBody, &postUser)

	switch {
	case strings.TrimSpace(postUser.Username) == "":
		error := helpers.ErrorMsg("no username in request")
		w.Write(error)
		helpers.EndpointError("no username in request", r)
		return
	case strings.TrimSpace(postUser.Fullname) == "":
		error := helpers.ErrorMsg("no fullname in request")
		w.Write(error)
		helpers.EndpointError("no fullname in request", r)
		return
	case strings.TrimSpace(postUser.Email) == "":
		error := helpers.ErrorMsg("no email in request")
		w.Write(error)
		helpers.EndpointError("no email in request", r)
		return
	}

	var passHashed, hashErr = helpers.HashPassword(postUser.Password)
	if hashErr != nil {
		error := helpers.ErrorMsg("error hashing password")
		w.Write(error)
		helpers.EndpointError("error hashing password", r)
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

	var userCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "username", Value: postUser.Username}}).Decode(&userCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("username already exists")
		w.Write(error)
		helpers.EndpointError("username already exists", r)
		return
	}

	createdAt := time.Now()

	doc := bson.D{{Key: "username", Value: strings.TrimSpace(postUser.Username)}, {Key: "fullname", Value: strings.TrimSpace(postUser.Fullname)}, {Key: "email", Value: strings.TrimSpace(postUser.Email)}, {Key: "password", Value: passHashed}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	documentID := fmt.Sprintf("%v", result.InsertedID)
	success := helpers.SuccessMsg(postUser.Username + " successfully added at " + documentID)
	w.Write(success)
	helpers.SuccessLog(r)
	return
}
