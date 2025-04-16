package controllers

import (
	"library-sys/database"
	"library-sys/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LibraryController struct{}

// 獲取所有 book
func (lt LibraryController) GetAll(c *gin.Context) {
	db := database.GetDB()
	var books []models.Book

	result := db.Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// 建立book
func (lt LibraryController) Create(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrRequiredFields})
		return
	}

	// 先檢查 ISBN 是否存在
	var existingBook models.Book
	db := database.GetDB()
	result := db.Where("isbn = ?", book.ISBN).First(&existingBook)
	if result.Error == nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":models.ErrDuplicateISBN})
		return
	}


	book.Status = models.StatusAvailable
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating book"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// 查詢特定的書
func (lt LibraryController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	var book models.Book
	db := database.GetDB()
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
		return
	}

	c.JSON(http.StatusOK, book)
}

// 更新書籍訊息
func (lt LibraryController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	var book models.Book
	db := database.GetDB()

	// 檢查書籍是否存在
	if err := db.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
		return
	}

	// 讀取更新的資料
	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 只更新允許的欄位
	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.ISBN = updatedBook.ISBN

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// 刪除書籍訊息
func (lt LibraryController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	db := database.GetDB()
	result := db.Delete(&models.Book{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}
	// 沒有任何記錄被刪除，表示要刪除的 ID 不存在
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// 借書
func (lt LibraryController) Borrow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidID})
		return
	}

	var borrowRequest models.BorrowRequest
	if err := c.ShouldBindJSON(&borrowRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var book models.Book

	// 使用交易確保操作的一致性 
	tx := db.Begin()

	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
		return
	}

	if book.Status == models.StatusBorrowed {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrBookBorrowed})
		return
	}

	now := time.Now()
	book.Status = models.StatusBorrowed
	book.Borrower = borrowRequest.Borrower
	book.Note = borrowRequest.Note
	book.BorrowedAt = &now

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, book)
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

	db := database.GetDB()
	var book models.Book

	// 使用事務確保操作的一致性
	tx := db.Begin()
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrBookNotFound})
		return
	}

	if book.Status == models.StatusAvailable {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrBookNotBorrowed})
		return
	}

	if book.Borrower != returnRequest.Borrower {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrWrongBorrower})
		return
	}

	book.Status = models.StatusAvailable
	book.Borrower = ""
	book.Note = ""
	book.BorrowedAt = nil

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, book)
}
