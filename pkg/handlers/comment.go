package handlers

import (
	"book-discusser/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createComment(c *gin.Context) {
	userId, err := getUserId(c)
	fmt.Println(userId)
	var input models.Comment
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Comment.Create(userId, input.ID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllCommentsResponse struct {
	Data []models.Comment `json:"data"`
}

func (h *Handler) getAllComments(c *gin.Context) {

	comments, err := h.services.Comment.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllCommentsResponse{
		Data: comments,
	})
}

type getAllUserCommentsResponse struct {
	Data []models.UsersComments `json:"data"`
}

func (h *Handler) getCommentByBookId(c *gin.Context) {
	userId, err := getUserId(c)
	user, err := h.services.Authorization.GetUserById(userId)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	comments, err := h.services.Comment.GetByBookId(id)
	if err != nil && comments != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user.Role == "ADMIN" {
		c.HTML(http.StatusOK, "admin_comments_page.html", gin.H{
			"title":   "Admin Comments Page",
			"payload": comments,
		})
	} else {
		c.HTML(http.StatusOK, "comment.html", gin.H{
			"title":     "Home Page",
			"userEmail": user.Email,
			"payload":   comments,
		})
	}
}

func (h *Handler) updateComment(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateCommentInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Comment.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteComment(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Comment.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
