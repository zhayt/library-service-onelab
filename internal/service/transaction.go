package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Transaction struct {
	transactionURL string
	itemURL        string
	deleteURL      string
	log            *zap.Logger
}

func NewTransaction(log *zap.Logger) *Transaction {
	return &Transaction{
		transactionURL: "http://localhost:8081/api/v1/transactions",
		itemURL:        "http://localhost:8081/api/v1/transactions/items",
		deleteURL:      "http://localhost:8081/api/v1/transactions/",
		log:            log}
}

func (s *Transaction) CreateTransaction(transaction model.Transaction) (int, error) {
	jsonData, _ := json.Marshal(transaction)

	req, err := prepareRequest(http.MethodPost, s.transactionURL, jsonData)
	if err != nil {
		return 0, fmt.Errorf("pepare request error: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request error: %w", err)
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
	jsonData, _ := json.Marshal(item)

	req, err := prepareRequest(http.MethodPost, s.itemURL, jsonData)
	if err != nil {
		return fmt.Errorf("pepare request error: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("transaction service error")
	}
}

func (s *Transaction) DeleteTransaction(transactionID int) error {
	req, err := prepareRequest(http.MethodDelete, s.deleteURL+strconv.Itoa(transactionID), []byte{})
	if err != nil {
		return fmt.Errorf("pepare request error: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request error: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("transaction service error")
	}
}

func prepareRequest(method string, url string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
