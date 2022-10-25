package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MiriConf/miriconf-backend/teams"

	_ "github.com/MiriConf/miriconf-backend/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           MiriConf Backend API
// @version         1.0
// @description     The backend API for MiriConf.
// @contact.name    MiriConf
// @contact.url     https://github.com/MiriConf/miriconf-backend
// @contact.email   bolmidgk@mail.uc.edu
// @license.name    GPL3
// @license.url     https://www.gnu.org/licenses/gpl-3.0.en.html

// @host      localhost:8081
// @BasePath  /api/v1

func main() {
	fmt.Println("miriconf-backend ready for requests...")

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("mongo URI is not specified, set with MONGO_URI environment variable")
	}

	fmt.Println("docs accesible at localhost:8081/docs/")

	r := mux.NewRouter()

	// API Endpoints

	// Teams
	r.HandleFunc("/api/v1/teams/list", teams.ListTeams).Methods("GET")
	r.HandleFunc("/api/v1/teams/get/{id}", teams.GetTeams).Methods("GET")
	r.HandleFunc("/api/v1/teams", teams.CreateTeams).Methods("POST")
	r.HandleFunc("/api/v1/teams/{id}", teams.EditTeams).Methods("PUT")
	r.HandleFunc("/api/v1/teams/{id}", teams.DeleteTeams).Methods("DELETE")

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Start server
	log.Fatal(http.ListenAndServe(":8081", r))
}
