package handlers

import (
	"book-discusser/configs"
	"book-discusser/pkg/models"
	"book-discusser/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (h *Handler) googleLogin(c *gin.Context) {
	oauthState := utils.GenerateStateOauthCookie(c.Writer)
	u := configs.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
}

func (h *Handler) googleCallback(c *gin.Context) {
	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Writer.Header().Add("content-type", "application/json")
	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
		newErrorResponse(c, http.StatusInternalServerError, "invalid oauth google state")
		return
	}
	// Exchange Auth Code for Tokens
	token, err := configs.AppConfig.GoogleLoginConfig.Exchange(
		context.Background(), code)
	// ERROR : Auth Code Exchange Failed
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("falied code exchange: %s", err.Error()))
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
		Role:  "USER",
	}
	isUserLogin, err := h.services.Authorization.GetUserByEmail(user.Email)
	if err != nil {
		if isUserLogin == nil {
			id, err := h.services.Authorization.CreateUser(user)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, "failed to login")
				return
			}
			_, err = h.createSession(c, id, user)
			if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("falied code exchange: %s", err.Error()))
			return
		}
	} else {
		_, err = h.createSession(c, isUserLogin.ID, *isUserLogin)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	// send back response to browser
	http.Redirect(c.Writer, c.Request, "/api/books", http.StatusPermanentRedirect)
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
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	sessionId, err := h.createSession(c, id, models.User{ID: id, Email: input.Email, Role: input.Role})
	if err != nil {
		fmt.Println(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"sessionToken": sessionId,
	})
}

func (h *Handler) login(c *gin.Context) {
	var input loginInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.services.Authorization.GetUser(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = h.createSession(c, user.ID, models.User{ID: user.ID, Email: input.Email, Role: user.Role})
	if err != nil {
		fmt.Println(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User Sign In successfully",
	})
}

func (h *Handler) logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cookie, err := c.Request.Cookie("mysession")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{
		"message": "User LogOut successfully",
	})
}

func (h *Handler) refreshSession(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := h.services.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cookie, err := c.Request.Cookie("mysession")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(c.Writer, cookie)
	sessionId, err := h.createSession(c, user.ID, *user)
	fmt.Println("Session refresh with id:", sessionId)
	c.Next()
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
