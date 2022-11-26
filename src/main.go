package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/apps"
	"github.com/MiriConf/miriconf-backend/systems"
	"github.com/MiriConf/miriconf-backend/teams"
	"github.com/MiriConf/miriconf-backend/templates"
	"github.com/MiriConf/miriconf-backend/users"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	hostName := os.Getenv("MIRICONF_HOSTNAME")
	if hostName == "" {
		log.Fatal("miriconf hostname is not specified, set with MIRICONF_HOSTNAME environment variable")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("mongo URI is not specified, set with MONGO_URI environment variable")
	}

	jwtKey := []byte(os.Getenv("JWT_KEY"))
	if len(jwtKey) == 0 {
		log.Fatal("JWT key is not specified, set with JWT_KEY environment variable")
	}

	fmt.Println("miriconf-backend ready for requests...")

	r := mux.NewRouter()
	r.Host("backend-svc")

	// API Endpoints

	// Teams
	r.HandleFunc("/api/v1/teams/list", teams.ListTeams).Methods("GET")
	r.HandleFunc("/api/v1/teams/get/{id}", teams.GetTeams).Methods("GET")
	r.HandleFunc("/api/v1/teams", teams.CreateTeams).Methods("POST")
	r.HandleFunc("/api/v1/teams/{id}", teams.EditTeams).Methods("PUT")
	r.HandleFunc("/api/v1/teams/{id}", teams.DeleteTeams).Methods("DELETE")
	// Users
	r.HandleFunc("/api/v1/login", users.Login).Methods("POST")
	r.HandleFunc("/api/v1/users/list", users.ListUsers).Methods("GET")
	r.HandleFunc("/api/v1/users/get/{id}", users.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users", users.CreateUsers).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", users.EditUsers).Methods("PUT")
	r.HandleFunc("/api/v1/users/{id}", users.DeleteUsers).Methods("DELETE")
	// Systems
	r.HandleFunc("/api/v1/systems/list", systems.ListSystems).Methods("GET")
	r.HandleFunc("/api/v1/systems/ping", systems.Ping).Methods("GET")
	r.HandleFunc("/api/v1/systems/fetch", systems.ClientFetch).Methods("GET")
	r.HandleFunc("/api/v1/systems/get/{id}", systems.GetSystems).Methods("GET")
	r.HandleFunc("/api/v1/systems", systems.CreateSystems).Methods("POST")
	r.HandleFunc("/api/v1/systems/{id}", systems.EditSystems).Methods("PUT")
	r.HandleFunc("/api/v1/systems/{id}", systems.DeleteSystems).Methods("DELETE")
	// Templates
	r.HandleFunc("/api/v1/template/build/{id}", templates.BuildTemplate).Methods("GET")
	r.HandleFunc("/api/v1/template/publish/{id}", templates.PublishTemplate).Methods("GET")
	// Applications
	r.HandleFunc("/api/v1/apps/list", apps.ListApps).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept", "Accept-Encoding"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
