package postgres

import (
	"context"
	"fmt"
	"github.com/zhayt/user-storage-service/internal/model"
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestBIHistoryStorage_CreateBIHistory(t *testing.T) {
	type args struct {
		ctx       context.Context
		bIHistory model.BIHistory
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{ctx: context.Background(),
			bIHistory: model.BIHistory{UserID: 1, Books: []*model.RentalBooks{{ID: 1, Quantity: 2}, {ID: 2, Quantity: 5}}},
		}, false},
		{"error bookID not exist", args{ctx: context.Background(),
			bIHistory: model.BIHistory{UserID: 1, Books: []*model.RentalBooks{{ID: 1, Quantity: 2}, {ID: 5, Quantity: 5}}},
		}, true},
		{"error userID not exist", args{ctx: context.Background(),
			bIHistory: model.BIHistory{UserID: 5, Books: []*model.RentalBooks{{ID: 1, Quantity: 2}, {ID: 5, Quantity: 5}}},
		}, true},
	}

	// db test container
	dbContainer, db, err := SetupTestDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbContainer.Terminate(context.Background())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BIHistoryStorage{
				db:  db,
				log: zap.NewExample(),
			}

			if err := r.CreateBIHistory(tt.args.ctx, tt.args.bIHistory); (err != nil) != tt.wantErr {
				t.Errorf("CreateBIHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBIHistoryStorage_GetCurrentBorrowedBooks(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		want    []model.BorrowedBooks
		wantErr bool
	}{
		{"success", args{context.Background()}, []model.BorrowedBooks{
			{ID: 1, UserName: "Test Fio Test", BookName: "Test book", BookAuthor: "Test Author", Quantity: 2, CreatedAt: time.Time{}},
			{ID: 1, UserName: "Test Fio Test", BookName: "Test book2", BookAuthor: "Test Author", Quantity: 5, CreatedAt: time.Time{}}}, false},
	}

	// db test container
	dbContainer, db, err := SetupTestDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer dbContainer.Terminate(context.Background())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &BIHistoryStorage{
				db:  db,
				log: zap.NewExample(),
			}

			got, err := r.GetCurrentBorrowedBooks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentBorrowedBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("got date: ", got)
		})
	}
}
