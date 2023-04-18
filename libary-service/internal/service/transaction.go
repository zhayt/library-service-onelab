package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
)

type Transaction struct {
	url string
	log *zap.Logger
}

func NewTransaction(log *zap.Logger) *Transaction {
	return &Transaction{
		url: "localhost:8081/api/v1",
		log: log}
}

func (s *Transaction) CreateTransaction(transaction model.Transaction) (int, error) {
	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		var transactionID int

		if err := json.NewDecoder(resp.Body).Decode(&transactionID); err != nil {
			return 0, err
		}

		return transactionID, nil
	} else {
		return 0, errors.New("transaction service error")
	}
}

func (s *Transaction) CreateTransactionItem(item model.TransactionItem) error {
	return nil
}

func (s *Transaction) DeleteTransaction(transactionID int) error {
	return nil
}
