package handlers

import (
	"book-discusser/pkg/service"
	"book-discusser/pkg/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) GetCookie(c *gin.Context) (*http.Cookie, error) {
	s, err := c.Request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			newErrorResponse(c, http.StatusUnauthorized, "empty cookie")
			return nil, err
		}
		// For any other type of error, return a bad request status
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return nil, err
	}
	return s, nil
}

func (h *Handler) GetSession(c *gin.Context) (*sessions.Session, error) {
	// We then get the session from our session map
	sessionToken, err := h.GetCookie(c)
	userSession, err := h.services.Authorization.GetSession(sessionToken.Value)
	if err != nil {
		// If the session token is not present in database, return an unauthorized error
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return nil, err
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.IsExpired() {
		err := h.services.Authorization.DeleteSession(userSession.ID)
		if err != nil {
			return nil, err
		}
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return nil, err
	}
	return userSession, nil
}

func (h *Handler) goForm(c *gin.Context) {
	c.HTML(http.StatusOK, "book_form.html", gin.H{
		"title": "Add Book Page",
	})
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("./resources/templates/*")
	router.Static("/css", "./resources/css")
	router.Static("/img", "./resources/img")
	auth := router.Group("/auth")
	{
		router.GET("/register", h.registerPage)
		router.POST("/register", h.register)
		router.GET("/login", h.loginPage)
		router.POST("/login", h.login)
		router.GET("/logout", h.logout)
		auth.GET("/google_login", h.googleLogin)
		auth.GET("/google_callback", h.googleCallback)
	}

	api := router.Group("/api", h.userIdentity)
	{
		books := api.Group("/books")
		{
			books.POST("/", h.createBook)
			books.GET("/", h.getAllBooks)
			books.GET("/:id", h.getBookByUserId)
			books.PUT("/:id", h.updateBook)
			books.DELETE("/:id", h.deleteComment)
			router.GET("/add_book", h.goForm)
		}

		comments := api.Group("/comments")
		{
			comments.POST("/", h.createComment)
			comments.GET("/", h.getAllComments)
			comments.GET("/:id", h.getCommentByBookId)
			comments.PUT("/:id", h.updateComment)
			comments.DELETE("/:id", h.deleteComment)
		}
	}
	return router
}
