package controller

import (
	"net/http"

	"github.com/MiriConf/miriconf-backend/model"
	"github.com/gin-gonic/gin"
)

// ListTeams godoc
// @Summary      List teams
// @Description  get teams
// @Tags         teams
// @Produce      json
// @Success      200  {array}   model.Team
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /teams/all [get]
func (c *Controller) ListTeams(ctx *gin.Context) {
	teams := model.TeamsAll()
	ctx.JSON(http.StatusOK, teams)
}
