package handlers

import (
	"book-discusser/configs"
	"book-discusser/pkg/models"
	"book-discusser/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

const (
	salt = "3476tg3gpjhb39hgm3pphnb3g"
)

func (h *Handler) googleLogin(c *gin.Context) {

	oauthState := utils.GenerateStateOauthCookie(c.Writer)
	u := configs.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusSeeOther)
}

func (h *Handler) googleCallback(c *gin.Context) {

	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Writer.Header().Add("content-type", "application/json")

	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		fmt.Fprintf(c.Writer, "invalid oauth google state")
		return
	}
	// Exchange Auth Code for Tokens
	token, err := configs.AppConfig.GoogleLoginConfig.Exchange(
		context.Background(), code)

	// ERROR : Auth Code Exchange Failed
	if err != nil {
		fmt.Fprintf(c.Writer, "falied code exchange: %s", err.Error())
		return
	}

	// Fetch User Data from google server
	response, err := http.Get(configs.OauthGoogleUrlAPI + token.AccessToken)

	// ERROR : Unable to get user data from google
	if err != nil {
		fmt.Fprintf(c.Writer, "failed getting user info: %s", err.Error())
		return
	}

	// Parse user data JSON Object
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	var oauthResponse map[string]interface{}
	err = json.Unmarshal(contents, &oauthResponse)
	if err != nil {
		fmt.Fprintf(c.Writer, "failed read response: %s", err.Error())
		return
	}
	user := models.User{
		Name:  oauthResponse["name"].(string),
		Email: oauthResponse["email"].(string),
	}

	id, err := h.services.Authorization.CreateUser(user)
	cookie := http.Cookie{
		Name:  "userId",
		Value: string(id),
	}
	http.SetCookie(c.Writer, &cookie)
	if err != nil {
		return
	}
	// send back response to browser
	fmt.Fprintln(c.Writer, string(id))
}

type loginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) register(c *gin.Context) {
	var input models.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.services.GenerateSessionToken(input.ID, input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   session.ID,
		Expires: session.Expiry,
	})
	c.JSON(http.StatusOK, map[string]interface{}{
		"sessionToken": session.ID,
	})
}

func (h *Handler) login(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}
	var input loginInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.services.GenerateSessionToken(id, input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   session.ID,
		Expires: session.Expiry,
	})
	c.JSON(http.StatusOK, map[string]interface{}{
		"sessionToken": session.ID,
	})
}

func (h *Handler) logout(c *gin.Context) {
	sessionToken, err := h.GetCookie(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// remove the users session from the session map
	err = h.services.Authorization.DeleteSession(sessionToken.Value)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
}

func (h *Handler) registerPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register Page",
	})
}

func (h *Handler) loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login Page",
	})
}
