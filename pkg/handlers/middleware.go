package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	userCtx    = "userId"
	roleCtx    = "role"
	expiryTime = 3600
)

func (h *Handler) userIdentity(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	role := session.Get("role")
	fmt.Println(c.Request.URL.Path)
	if sessionID != nil && c.Request.URL.Path != "/auth/logout" {
		c.Set(userCtx, sessionID)
		c.Set(roleCtx, role)
		c.Redirect(http.StatusTemporaryRedirect, "/api/books")
	}
}

func (h *Handler) isLogin(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	role := session.Get("role")
	if sessionID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/view/login")
		return
	}
	c.Set(userCtx, sessionID)
	c.Set(roleCtx, role)
}

func (h *Handler) home(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	if sessionID == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/view/login")
		return
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/api/books")
	}
}

func (h *Handler) isExpiry(c *gin.Context) bool {
	session := sessions.Default(c)
	sessionID := session.Get("id")
	if sessionID != nil {
		sessionTimeExpires := session.Get("expires")
		fmt.Println(sessionTimeExpires)
		timeExpires, err := time.Parse(time.DateTime, sessionTimeExpires.(string))
		fmt.Println(timeExpires)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return false
		}
		if time.Now().After(timeExpires) {
			return true
		}
	}
	return false
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not of valid type")
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}

func getRole(c *gin.Context) (string, error) {
	role, ok := c.Get(roleCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return "", errors.New("user id not found")
	}
	userRole, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not of valid type")
		return "", errors.New("user id not found")
	}
	return userRole, nil
}
