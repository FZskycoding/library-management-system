package services

import (
	"library-sys/models"
	"time"

	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func CreateBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

// GetAllBooks 獲取所有書籍
func (s *BookService) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	result := s.db.Find(&books)
	return books, result.Error
}

// CreateBook 創建新書籍
func (s *BookService) CreateBook(book *models.Book) error {
	// 檢查 ISBN 是否存在
	var existingBook models.Book
	result := s.db.Where("isbn = ?", book.ISBN).First(&existingBook)
	if result.Error == nil {
		return models.ErrDuplicateISBN
	}

	book.Status = models.StatusAvailable
	return s.db.Create(book).Error
}

// GetBookByID 通過 ID 獲取書籍
func (s *BookService) GetBookByID(id int) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// UpdateBook 更新書籍信息
func (s *BookService) UpdateBook(id int, updatedBook *models.Book) (*models.Book, error) {
	book, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.ISBN = updatedBook.ISBN

	if err := s.db.Save(book).Error; err != nil {
		return nil, err
	}
	return book, nil
}

// DeleteBook 刪除書籍
func (s *BookService) DeleteBook(id int) error {
	result := s.db.Delete(&models.Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrBookNotFound
	}
	return nil
}

// BorrowBook 借書
func (s *BookService) BorrowBook(id int, borrowRequest *models.BorrowRequest) (*models.Book, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var book models.Book
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if book.Status == models.StatusBorrowed {
		tx.Rollback()
		return nil, models.ErrBookBorrowed
	}

	now := time.Now()
	book.Status = models.StatusBorrowed
	book.Borrower = borrowRequest.Borrower
	book.Note = borrowRequest.Note
	book.BorrowedAt = &now

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &book, nil
}

// ReturnBook 還書
func (s *BookService) ReturnBook(id int, returnRequest *models.ReturnRequest) (*models.Book, error) {
	tx := s.db.Begin() //建立一個transaction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var book models.Book
	if err := tx.First(&book, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//檢查還書人是否等於借書人
	if book.Borrower != returnRequest.Borrower {
		tx.Rollback()
		return nil, models.ErrWrongBorrower
	}

	//更新書本狀態，並存回資料庫
	book.Status = models.StatusAvailable
	book.Borrower = ""
	book.Note = ""
	book.BorrowedAt = nil
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &book, nil
}
