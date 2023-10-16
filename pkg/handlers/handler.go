package handlers

import (
	"book-discusser/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.GET("/google_login", h.googleLogin)
		auth.GET("/google_callback", h.googleCallback)
	}

	api := router.Group("/api")
	{
		books := api.Group("/books")
		{
			books.POST("/", h.createBook)
			books.GET("/", h.getAllBooks)
			books.GET("/:id", h.getBookByUserId)
			books.PUT("/:id", h.updateBook)
			books.DELETE("/:id", h.deleteComment)
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
