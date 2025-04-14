package models

import (
	"time"
)

type Library struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	ISBN       string     `json:"isbn"`
	Status     string     `json:"status"`
	BorrowedAt *time.Time `json:"borrowed_at,omitempty"`
	Borrower   string     `json:"borrower,omitempty"`
	Note       string     `json:"note,omitempty"`
}

var Libraries = make([]Library, 0)
