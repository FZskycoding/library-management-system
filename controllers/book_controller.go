package controllers

import (
	"library-sys/database"
	"library-sys/models"
	"library-sys/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LibraryController struct {
	bookService *services.BookService
}

// 改用函數來創建控制器
func DefaultController() LibraryController {
	return LibraryController{
		bookService: services.CreateBookService(database.GetDB()),
	}
}

// 查詢所有 book
func (lc LibraryController) GetAll(c *gin.Context) {
	books, err := lc.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookFetch.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// 建立book
func (lc LibraryController) Create(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrRequiredFields.Error()})
		return
	}

	if err := lc.bookService.CreateBook(&book); err != nil {
		if err == models.ErrDuplicateISBN {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookCreate.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, book)
}

// 查詢特定的書
func (lc LibraryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID.Error()})
		return
	}

	book, err := lc.bookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// 更新書籍訊息
func (lc LibraryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID.Error()})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := lc.bookService.UpdateBook(id, &updatedBook)
	if err != nil {
		if err == models.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookUpdate.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// 刪除書籍訊息
func (lc LibraryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID.Error()})
		return
	}

	err = lc.bookService.DeleteBook(id)
	if err != nil {
		if err == models.ErrBookNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookDelete.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// 借書
func (lc LibraryController) Borrow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID.Error()})
		return
	}

	var borrowRequest models.BorrowRequest
	if err := c.ShouldBindJSON(&borrowRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := lc.bookService.BorrowBook(id, &borrowRequest)
	if err != nil {
		switch err {
		case models.ErrBookNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case models.ErrBookBorrowed:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookUpdate.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// 還書
func (lc LibraryController) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID.Error()})
		return
	}

	var returnRequest models.ReturnRequest
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := lc.bookService.ReturnBook(id, &returnRequest)
	if err != nil {
		switch err {
		case models.ErrBookNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case models.ErrBookNotBorrowed, models.ErrWrongBorrower:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrBookUpdate.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}
