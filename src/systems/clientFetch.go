package systems

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MiriConf/miriconf-backend/helpers"
	"github.com/MiriConf/miriconf-backend/teams"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClientFetch(w http.ResponseWriter, r *http.Request) {
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

	systemData := strings.Split(headerToken, ".")

	tokenDec, err := base64.RawStdEncoding.DecodeString(systemData[1])
	if err != nil {
		panic(err)
	}

	var systemId SystemID
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

	coll := client.Database("miriconf").Collection("systems")

	var nameCheck System
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: sysID}}).Decode(&nameCheck)
	if err != nil {
		error := helpers.ErrorMsg("system does not exist or cannot be found")
		w.Write(error)
		helpers.EndpointError("system does not exist or cannot be found", r)
		fmt.Println(err)
		return
	}

	coll = client.Database("miriconf").Collection("teams")

	teamID, err := primitive.ObjectIDFromHex(nameCheck.Team)
	if err != nil {
		error := helpers.ErrorMsg("invalid team id requested")
		w.Write(error)
		helpers.EndpointError("invalid team id requested", r)
		return
	}

	var teamCheck teams.Team
	err = coll.FindOne(context.TODO(), bson.D{{Key: "_id", Value: teamID}}).Decode(&teamCheck)
	if err != nil {
		error := helpers.ErrorMsg("team does not exist or cannot be found")
		w.Write(error)
		helpers.EndpointError("team does not exist or cannot be found", r)
		fmt.Println(err)
		return
	}

	tmpDir := uuid.New()
	directory := "/tmp/" + tmpDir.String()
	_, err = os.Stat(directory)
	if os.IsNotExist(err) {
		fmt.Printf("%v does not exist, attempting to create it now...\n", directory)
		err := os.Mkdir(directory, 0755)
		if err != nil {
			panic(err)
		}
		fmt.Println("created directory at " + directory)
	}

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &githttp.BasicAuth{
			Username: "null",
			Password: teamCheck.SourcePAT,
		},
		URL:               teamCheck.SourceRepo + ".git",
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})
	if err != nil {
		fmt.Printf("failed to clone repo at %v... rolling back...\n", teamCheck.SourceRepo)
		err := os.RemoveAll(directory)
		if err != nil {
			panic(err)
		}
		cloneError := "could not clone repo at " + teamCheck.SourceRepo + " with the provided credentials"
		cloneErr := helpers.ErrorMsg(cloneError)
		w.Write(cloneErr)
		return
	}

	resultZip := directory + tmpDir.String() + ".zip"

	zipFile, err := os.Create(resultZip)
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(zipFile)

	dirWalk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk("tmp/"+tmpDir.String()+"/", dirWalk)
	if err != nil {
		panic(err)
	}

	zipWriter.Close()
	zipFile.Close()

	http.ServeFile(w, r, directory+tmpDir.String()+".zip")
}
