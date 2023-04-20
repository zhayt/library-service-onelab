package model

import "time"

type Transaction struct {
	ID        uint      `json:"id"`
	UserName  string    `json:"name"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type TransactionItem struct {
	ID            uint `json:"ID"`
	TransactionID uint `json:"transactionID"`
	Book          *Book
}
