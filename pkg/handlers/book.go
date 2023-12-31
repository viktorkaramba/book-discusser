package handlers

import (
	"book-discusser/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createBook(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.UserCreateBook
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Book.Create(userId, input)
	if err != nil {
		//newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllBooks(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	books, err := h.services.Book.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if role == "ADMIN" {
		c.HTML(http.StatusOK, "admin_page.html", gin.H{
			"title":   "Admin Page",
			"payload": books,
		})
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "Home Page",
			"payload": books,
		})
	}
}

func (h *Handler) getBookByUserId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	books, err := h.services.Book.GetByUserId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) updateBook(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateBookInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(*input.ImageBook)
	if err := h.services.Book.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteBook(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Book.Delete(id)
	if err != nil {
		//newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
