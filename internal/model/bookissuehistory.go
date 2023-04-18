package model

import "time"

type BIHistory struct {
	ID         int            `json:"id"`
	Books      []*RentalBooks `json:"bookID"`
	UserID     int            `json:"userID"`
	CreatedAt  time.Time      `json:"issueDate"`
	ReturnDate time.Time      `json:"returnDate"`
}

type RentalBooks struct {
	ID       int `json:"ID"`
	Quantity int `json:"quantity"`
}

type BorrowedBooks struct {
	ID         int       `json:"id"`
	UserName   string    `json:"userName" db:"fio"`
	BookName   string    `json:"bookName" db:"name"`
	BookAuthor string    `json:"bookAuthor" db:"author"`
	Quantity   int       `json:"quantity" db:"quantity"`
	CreatedAt  time.Time `json:"issueDate" db:"created_at"`
}
