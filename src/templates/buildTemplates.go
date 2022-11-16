package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MiriConf/miriconf-backend/helpers"
	"github.com/MiriConf/miriconf-backend/teams"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BuildTemplate(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	hostName := os.Getenv("MIRICONF_HOSTNAME")
	if hostName == "" {
		log.Fatal("miriconf hostname is not specified, set with MIRICONF_HOSTNAME environment variable")
	}
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

	var result teams.Team
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

	directory := "/mnt/data/" + result.Name
	_, err = os.Stat(directory)
	if os.IsNotExist(err) {
		panic("directory does not exist... build failed")
	}

	openRepo, err := git.PlainOpen(directory)
	if err != nil {
		panic(err)
	}

	repoTree, err := openRepo.Worktree()
	if err != nil {
		panic(err)
	}

	fmt.Printf("pulling latest commits from %v\n", result.SourceRepo)
	err = repoTree.Pull(&git.PullOptions{Auth: &githttp.BasicAuth{Username: "null", Password: result.SourcePAT}, RemoteName: "origin", Progress: os.Stdout})
	if err != nil && err == git.NoErrAlreadyUpToDate {
		fmt.Println("no changes to pull")
	} else if err != nil {
		panic(err)
	}

	repoTree, err = openRepo.Worktree()
	if err != nil {
		panic(err)
	}

	filename := filepath.Join(directory, "newfile.txt")
	err = ioutil.WriteFile(filename, []byte("testing testing testing"), 0644)
	if err != nil {
		panic(err)
	}

	_, err = repoTree.Add("newfile.txt")
	if err != nil {
		panic(err)
	}

	repoStatus, err := repoTree.Status()
	if err != nil {
		panic(err)
	}

	fmt.Println(repoStatus)

	newCommit, err := repoTree.Commit("test commit", &git.CommitOptions{Author: &object.Signature{Name: "miriconf-bot", Email: "miriconf-bot@" + hostName, When: time.Now()}})
	if err != nil {
		panic(err)
	}

	obj, err := openRepo.CommitObject(newCommit)
	if err != nil {
		panic(err)
	}

	helpers.SuccessLog(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("build commited successfully with message: " + obj.Message)
}

func ParseTemplate() {
	data := struct {
		Username string
		GitUser  string
		GitEmail string
	}{
		Username: "gbolmida",
		GitUser:  "test-user",
		GitEmail: "gbolmida@georgebolmida.com",
	}

	result, err := os.Create("/templates/git.nix")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	template, err := template.ParseFiles("/templates/nix.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	err = template.Execute(result, data)
	if err != nil {
		log.Fatal(err)
	}

	result, err = os.Open("/templates/git.nix")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, result)
}
