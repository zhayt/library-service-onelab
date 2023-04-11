package model

import "time"

type BIHistory struct {
	ID         int       `json:"id"`
	BookID     int       `json:"bookID"`
	UserID     int       `json:"userID"`
	IssueDate  time.Time `json:"issueDate"`
	ReturnDate time.Time `json:"returnDate"`
}
