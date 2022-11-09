package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/helpers"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Templates(w http.ResponseWriter, r *http.Request) {
	//mongoURI := os.Getenv("MONGO_URI")

	ParseTemplate()

	endpoint := "minio-svc:9000"
	accessKeyID := "TgcVtoIaKttPFqUC"
	secretAccessKey := "LE6jWe3EGJLrv4apKDZkUZu3W7Eob37f"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := minioClient.ListObjects(ctx, "test", minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}

	object, err := os.Open("/templates/git.nix")
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	n, err := minioClient.PutObject(ctx, "test", "git.nix", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		panic(err)
	}
	fmt.Println(n)

	helpers.SuccessLog(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("aaaaaa")
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
