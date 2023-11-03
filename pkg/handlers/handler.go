package handlers

import (
	"book-discusser/pkg/models"
	"book-discusser/pkg/service"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) goForm(c *gin.Context) {
	c.HTML(http.StatusOK, "book_form.html", gin.H{
		"title": "Add Book Page",
	})
}

func (h *Handler) createSession(c *gin.Context, id int, user models.User) (string, error) {
	session := sessions.Default(c)
	session.Set("id", id)
	session.Set("email", user.Email)
	session.Set("role", user.Role)
	expiryTimeSession := time.Now().Add(time.Second * expiryTime).Format(time.DateTime)
	session.Set("expires", expiryTimeSession)
	session.Options(sessions.Options{
		Path:     "/",
		Secure:   false,
		MaxAge:   expiryTime,
		HttpOnly: true,
	})
	err := session.Save()
	if err != nil {
		fmt.Println(err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return "", err
	}
	return session.ID(), nil
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.LoadHTMLGlob("./resources/templates/*")
	router.Static("/css", "./resources/css")
	router.Static("/img", "./resources/img")
	router.GET("/", h.home)
	auth := router.Group("/auth")
	{
		authView := auth.Group("/view")
		authView.Use(h.userIdentity)
		{
			authView.Static("/css", "./resources/css")
			authView.Static("/img", "./resources/img")
			authView.GET("/register", h.registerPage)
			authView.GET("/login", h.loginPage)

		}
		auth.Static("/css", "./resources/css")
		auth.Static("/img", "./resources/img")
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.GET("/logout", h.logout)
		auth.GET("/google_login", h.googleLogin)
		auth.GET("/google_callback", h.googleCallback)
	}
	api := router.Group("/api")
	api.Use(h.isLogin)
	{
		api.Static("/css", "./resources/css")
		api.Static("/img", "./resources/img")
		books := api.Group("/books")
		books.Use(h.isLogin)
		{
			books.Static("/css", "./resources/css")
			books.Static("/img", "./resources/img")
			books.POST("/", h.createBook)
			books.GET("/", h.getAllBooks)
			books.GET("/:id", h.getBookByUserId)
			books.PUT("/:id", h.updateBook)
			books.DELETE("/:id", h.deleteBook)
		}
		comments := api.Group("/comments")
		comments.Use(h.isLogin)
		{
			comments.Static("/css", "./resources/css")
			comments.Static("/img", "./resources/img")
			comments.POST("/", h.createComment)
			comments.GET("/", h.getAllComments)
			comments.PUT("/:id", h.updateComment)
			comments.DELETE("/:id", h.deleteComment)
		}
		bookComments := books.Group("/comments")
		bookComments.Use(h.isLogin)
		{
			bookComments.Static("/css", "./resources/css")
			bookComments.Static("/img", "./resources/img")
			bookComments.GET("/:id", h.getCommentByBookId)
		}
		api.GET("/add-book", h.goForm)
	}
	return router
}
