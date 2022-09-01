package main

import (
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

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	mainRouter := gin.Default()

	mainController := controller.NewController()

	v1 := mainRouter.Group("/api/v1")
	{
		teams := v1.Group("/teams")
		{
			// root teams
			teams.GET("", mainController.ListTeams)
			//teams.POST("", mainController.AddTeam)
			//teams.GET(":id", mainController.ShowTeam)
			//teams.PATCH(":id", mainController.UpdateTeam)
			//teams.DELETE(":id", mainController.DeleteTeam)
			// teams meta
			//teams.GET(":id/members", mainController.ListTeamMembers)
		}
	}
	mainRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mainRouter.Run(":8080")
}
