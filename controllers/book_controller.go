package controllers

import (
	"library-sys/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LibraryController struct{}

// 獲取所有 book
func (t LibraryController) GetAll(c *gin.Context) {

	c.JSON(http.StatusOK, models.Libraries)
}

// 建立book
func (lt LibraryController) Create(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrRequiredFields})
		return
	}

	book.ID = len(models.Libraries) + 1
	book.Status = models.StatusAvailable
	models.Libraries = append(models.Libraries, book)
	c.JSON(http.StatusCreated, book)
}

// 查詢特定的書
func (lt LibraryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	for _, book := range models.Libraries {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
}

// 更新書籍訊息
func (lt LibraryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	var updatedBook models.Book

	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			// 保留原有的 ID
			updatedBook.ID = book.ID
			// 保留借閱相關資訊
			updatedBook.Status = book.Status
			updatedBook.BorrowedAt = book.BorrowedAt
			updatedBook.Borrower = book.Borrower
			updatedBook.Note = book.Note

			// 直接更新整個結構體
			models.Libraries[i] = updatedBook
			c.JSON(http.StatusOK, updatedBook)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
}

// 刪除書籍訊息
func (lt LibraryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			models.Libraries = append(models.Libraries[:i], models.Libraries[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
}

// 借書
func (lt LibraryController) Borrow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	var borrowRequest models.BorrowRequest
	now := time.Now()
	if err := c.ShouldBindJSON(&borrowRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			//檢查是否已被借出
			if models.Libraries[i].Status == models.StatusBorrowed {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrBookBorrowed})
				return
			}
			models.Libraries[i].Status = models.StatusBorrowed
			models.Libraries[i].Borrower = borrowRequest.Borrower
			models.Libraries[i].Note = borrowRequest.Note
			models.Libraries[i].BorrowedAt = &now
			c.JSON(http.StatusOK, models.Libraries[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
}

// 還書
func (lt LibraryController) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}
	var returnRequest models.ReturnRequest
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, book := range models.Libraries {
		if id == book.ID {
			// 檢查書是否已被借出
			if book.Status == models.StatusAvailable {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrBookNotBorrowed})
				return
			}

			// 檢查還書的人是否為借書的人
			if book.Borrower != returnRequest.Borrower {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrWrongBorrower})
				return
			}
			models.Libraries[i].Status = models.StatusAvailable
			models.Libraries[i].Borrower = ""
			models.Libraries[i].Note = ""
			models.Libraries[i].BorrowedAt = nil
			c.JSON(http.StatusOK, models.Libraries[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
}
