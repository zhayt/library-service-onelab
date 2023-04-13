package model

import "time"

type BIHistory struct {
	ID         int       `json:"id"`
	BookID     int       `json:"bookID"`
	UserID     int       `json:"userID"`
	CreatedAt  time.Time `json:"issueDate"`
	ReturnDate time.Time `json:"returnDate"`
}

type BorrowedBooks struct {
	ID         int       `json:"id"`
	UserName   string    `json:"userName" db:"fio"`
	BookName   string    `json:"bookName" db:"name"`
	BookAuthor string    `json:"bookAuthor" db:"author"`
	CreatedAt  time.Time `json:"issueDate" db:"created_at"`
}
