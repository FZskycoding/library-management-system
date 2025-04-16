package models

import (
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

// 其他錯誤訊息常量
const (
	ErrBookNotFound    = "Book not found"
	ErrInvalidID       = "Invalid ID format"
	ErrBookBorrowed    = "Book is already borrowed"
	ErrBookNotBorrowed = "Book is not borrowed"
	ErrWrongBorrower   = "Only the borrower can return this book"
	ErrRequiredFields  = "新增書籍必須包含title、isbn、author!"
	ErrDuplicateISBN   = "book with this ISBN already exists"
)
