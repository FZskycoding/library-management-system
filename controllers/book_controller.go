package controllers

import (
	"library-sys/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LibraryController struct{}

// GetAll 獲取所有 book
func (t LibraryController) GetAll(c *gin.Context) {

	c.JSON(http.StatusOK, models.Libraries)
}

// 建立book
func (lt LibraryController) Create(c *gin.Context) {
	var book models.Library
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = len(models.Libraries) + 1
	// book.BorrowedAt = time.Now()
	book.Status = "available"
	models.Libraries = append(models.Libraries, book)
	c.JSON(http.StatusCreated, book)
}

// 查詢特定的書
func (lt LibraryController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, book := range models.Libraries {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// 更新書籍訊息
func (lt LibraryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var updatedBook models.Library

	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			updatedBook.ID = id

			if updatedBook.Title != "" {
				models.Libraries[i].Title = updatedBook.Title
			}
			if updatedBook.Author != "" {
				models.Libraries[i].Author = updatedBook.Author
			}
			if updatedBook.ISBN != "" {
				models.Libraries[i].ISBN = updatedBook.ISBN
			}
			if updatedBook.Status != "" {
				models.Libraries[i].Status = updatedBook.Status
			}
			// models.Libraries[i] = updatedBook
			c.JSON(http.StatusOK, updatedBook)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// 刪除書籍訊息
func (lt LibraryController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, book := range models.Libraries {
		if id == book.ID {
			models.Libraries = append(models.Libraries[:i], models.Libraries[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// 借書
func (lt LibraryController) Borrow(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var borrowRequest models.BorrowRequest
	now := time.Now()
	if err := c.ShouldBindJSON(&borrowRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			//檢查是否已被借出
			if models.Libraries[i].Status == "borrowed" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Book is already borrowed"})
				return
			}
			models.Libraries[i].Status = "borrowed"
			models.Libraries[i].Borrower = borrowRequest.Borrower
			models.Libraries[i].Note = borrowRequest.Note
			models.Libraries[i].BorrowedAt = &now
			c.JSON(http.StatusOK, models.Libraries[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// 還書
func (lt LibraryController) Return(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var returnRequest models.ReturnRequest
	now := time.Now()
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {

			models.Libraries[i].Status = "available"
			models.Libraries[i].Borrower = ""
			models.Libraries[i].Note = ""
			models.Libraries[i].BorrowedAt = &now
			c.JSON(http.StatusOK, models.Libraries[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}
