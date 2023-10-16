package handlers

import (
	"book-discusser/configs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

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

	// send back response to browser
	fmt.Fprintln(c.Writer, string(contents))
}
