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
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	var repoCheck bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "source_repo", Value: postTeam.SourceRepo}}).Decode(&repoCheck)
	if err != mongo.ErrNoDocuments {
		error := helpers.ErrorMsg("team with this repo already exists")
		w.Write(error)
		helpers.EndpointError("team with this repo already exists", r)
		return
	}

	directory := "/mnt/data/" + postTeam.Name
	_, err = os.Stat(directory)
	if os.IsNotExist(err) {
		fmt.Printf("%v does not exist yet, attempting to create it now...\n", directory)
		err := os.Mkdir(directory, 0755)
		if err != nil {
			panic(err)
		}
		fmt.Println("created directory at " + directory + " for team " + postTeam.Name)
	}

	gitClone, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "null",
			Password: postTeam.SourcePAT,
		},
		URL:               postTeam.SourceRepo + ".git",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})
	if err != nil {
		fmt.Printf("failed to clone repo at %v... rolling back...\n", postTeam.SourceRepo)
		err := os.RemoveAll(directory)
		if err != nil {
			panic(err)
		}
		cloneError := "could not clone repo at " + postTeam.SourceRepo + " with the provided credentials"
		cloneErr := helpers.ErrorMsg(cloneError)
		w.Write(cloneErr)
		return
	}

	ref, err := gitClone.Head()
	if err != nil {
		panic(err)
	}

	commit, err := gitClone.CommitObject(ref.Hash())
	if err != nil {
		panic(err)
	}

	fmt.Println(commit)

	createdAt := time.Now()

	doc := bson.D{{Key: "name", Value: strings.TrimSpace(postTeam.Name)}, {Key: "department", Value: strings.TrimSpace(postTeam.Department)}, {Key: "source_repo", Value: strings.TrimSpace(postTeam.SourceRepo)}, {Key: "source_pat", Value: strings.TrimSpace(postTeam.SourcePAT)}, {Key: "createdat", Value: createdAt.Format("01-02-2006 15:04:05")}}

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
