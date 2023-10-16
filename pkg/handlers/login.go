package handlers

import (
	"book-discusser/configs"
	"book-discusser/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) googleLogin(c *gin.Context) {

	oauthState := utils.GenerateStateOauthCookie(c.Writer)
	u := configs.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusSeeOther)
}
