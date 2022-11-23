package templates

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MiriConf/miriconf-backend/helpers"
	"github.com/MiriConf/miriconf-backend/systems"
	"github.com/MiriConf/miriconf-backend/teams"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PublishTemplate(w http.ResponseWriter, r *http.Request) {
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

	systemData := strings.Split(headerToken, ".")

	tokenDec, err := base64.StdEncoding.DecodeString(systemData[1])
	if err != nil {
		panic(err)
	}

	var systemId systems.SystemID
	json.Unmarshal(tokenDec, &systemId)

	sysID, err := primitive.ObjectIDFromHex(systemId.SystemID)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collSystems := client.Database("miriconf").Collection("systems")

	var systemResult systems.System
	err = collSystems.FindOne(context.TODO(), bson.D{{Key: "_id", Value: sysID}}).Decode(&systemResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err != nil {
				error := helpers.ErrorMsg("no system matching id requested")
				w.Write(error)
				helpers.EndpointError("no system matching id requested", r)
				return
			}
		}
		log.Fatal(err)
	}

	coll := client.Database("miriconf").Collection("teams")

	var result teams.Team
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: systemResult.Team}}).Decode(&result)
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

	fmt.Println(result)

	//directory := "/mnt/data/" + result.Name + "-published"
	//_, err = os.Stat(directory)
	//if os.IsNotExist(err) {
	//	fmt.Printf("%v does not exist yet, attempting to create it now...\n", directory)
	//	err := os.Mkdir(directory, 0755)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("created directory at " + directory + " for team " + result.Name)
	//}
	//
	//gitClone, err := git.PlainClone(directory, false, &git.CloneOptions{
	//	Auth: &githttp.BasicAuth{
	//		Username: "null",
	//		Password: postTeam.SourcePAT,
	//	},
	//	URL:               postTeam.SourceRepo + ".git",
	//	RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	//	Progress:          os.Stdout,
	//})
	//if err != nil {
	//	fmt.Printf("failed to clone repo at %v... rolling back...\n", postTeam.SourceRepo)
	//	err := os.RemoveAll(directory)
	//	if err != nil {
	//		panic(err)
	//	}
	//	cloneError := "could not clone repo at " + postTeam.SourceRepo + " with the provided credentials"
	//	cloneErr := helpers.ErrorMsg(cloneError)
	//	w.Write(cloneErr)
	//	return
	//}
	//
	//ref, err := gitClone.Head()
	//if err != nil {
	//	panic(err)
	//}
	//
	//commit, err := gitClone.CommitObject(ref.Hash())
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(commit)
	//
	//doc := bson.D{{"$set", bson.D{{Key: "commit_sha", Value: }}}}
	//filter := bson.D{{Key: "_id", Value: teamID}}
	//
	//_, err = coll.UpdateOne(context.TODO(), filter, doc)
	//if err != nil {
	//	panic(err)
	//}

	helpers.SuccessLog(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("changes pushed successfully")
}
