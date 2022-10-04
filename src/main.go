package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MiriConf/miriconf-backend/controller"

	_ "github.com/MiriConf/miriconf-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	count := 5
	for count > 0 {
		fmt.Printf("miriconf-backend is starting in %v...\n", count)
		time.Sleep(time.Second)
		count--
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("mongo URI is not specified, set with MONGO_URI environment variable")
	}

	fmt.Println("docs accesible at /swagger/index.html")

	mainRouter := gin.Default()
	mainController := controller.NewController()

	v1 := mainRouter.Group("/api/v1")
	{
		teams := v1.Group("/teams")
		{
			// root teams
			teams.GET("/all", mainController.ListTeams)
			//teams.POST("", mainController.AddTeam)
			//teams.GET(":id", mainController.ShowTeam)
			//teams.PATCH(":id", mainController.UpdateTeam)
			//teams.DELETE(":id", mainController.DeleteTeam)
			// teams meta
			//teams.GET(":id/members", mainController.ListTeamMembers)
		}
	}
	mainRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mainRouter.Run(":8081")
}
