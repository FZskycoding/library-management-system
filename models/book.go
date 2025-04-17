package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model // 這會自動添加 ID, CreatedAt, UpdatedAt, DeletedAt 欄位

	Title      string     `json:"title" binding:"required"`
	Author     string     `json:"author" binding:"required"`
	ISBN       string     `json:"isbn" binding:"required" gorm:"uniqueIndex"`
	Status     string     `json:"status"`
	BorrowedAt *time.Time `json:"borrowed_at,omitempty"`
	Borrower   string     `json:"borrower,omitempty"`
	Note       string     `json:"note,omitempty"`
}

type BorrowRequest struct {
	Borrower string `json:"borrower" binding:"required"`
	Note     string `json:"note" binding:"required"`
}

type ReturnRequest struct {
	Borrower string `json:"borrower" binding:"required"`
}

// 書籍狀態常量
const (
	StatusAvailable = "available" // 可借狀態
	StatusBorrowed  = "borrowed"  // 已借出狀態
)

// 錯誤定義
var (
	ErrBookNotFound    = errors.New("Book not found")
	ErrBookBorrowed    = errors.New("Book is already borrowed")
	ErrBookNotBorrowed = errors.New("Book is not borrowed")
	ErrInvalidID       = errors.New("invalid ID format")
	ErrWrongBorrower   = errors.New("only the borrower can return this book")
	ErrRequiredFields  = errors.New("add new book should include title,isbn,author")
	ErrDuplicateISBN   = errors.New("book with this ISBN already exists")
	ErrBookUpdate      = errors.New("error updating book")
	ErrBookFetch       = errors.New("error fetching books")
	ErrBookDelete      = errors.New("error deleting book")
	ErrBookCreate      = errors.New("error creating book")
)

