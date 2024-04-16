package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"pustaka-api/book"
)

type bookHandler struct {
	bookService book.Service
}

func NewBookHandler(bookService book.Service) *bookHandler {
	return &bookHandler{bookService}
}

func (h *bookHandler) GetBooks(c *gin.Context) {
	rawBooks, err := h.bookService.FindAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})

		return
	}

	var books []book.BookResponse
	for _, b := range rawBooks {
		bookResponse := nconvertToBookResponse(b)
		books = append(books, bookResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})
}

func (h *bookHandler) GetBook(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	b, err := h.bookService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertToBookResponse(b),
	})
}

func (h *bookHandler) CreateBook(c *gin.Context) {
	var bookRequest book.BookRequest

	err := c.ShouldBindJSON(&bookRequest)
	if err != nil {
		errorMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg := fmt.Sprintf("Error on field %s, condition %s", e.Field(), e.ActualTag())
			errorMsgs = append(errorMsgs, errorMsg)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errorMsgs,
		})

		return
	}

	book, err := h.bookService.Create(bookRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertToBookResponse(book),
	})
}

func (h *bookHandler) UpdateBook(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	var bookRequest book.BookRequest
	err := c.ShouldBindJSON(&bookRequest)
	if err != nil {
		errorMsgs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg := fmt.Sprintf("Error on field %s, condition %s", e.Field(), e.ActualTag())
			errorMsgs = append(errorMsgs, errorMsg)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errorMsgs,
		})

		return
	}

	book, updateError := h.bookService.Update(id, bookRequest)
	if updateError != nil {
		if errors.Is(updateError, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"errors": "Record not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": updateError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertToBookResponse(book),
	})
}

func (h *bookHandler) DeleteBook(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	book, updateError := h.bookService.Delete(id)
	if updateError != nil {
		if errors.Is(updateError, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"errors": "Record not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": updateError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertToBookResponse(book),
	})
}

func convertToBookResponse(b book.Book) book.BookResponse {
	return book.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Price:       b.Price,
		Description: b.Description,
		Rating:      b.Rating,
	}
}

func nconvertToBookResponse(b *book.Book) book.BookResponse {
	return book.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Price:       b.Price,
		Description: b.Description,
		Rating:      b.Rating,
	}
}
