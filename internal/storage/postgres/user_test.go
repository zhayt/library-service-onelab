package postgres

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

func TestUserStorage_GetUserByEmail(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{context.Background(), "existsemail@mail.ru"}, false},
		{"email not exists", args{context.Background(), "notexistsemail@mail.ru"}, true},
	}

	// up test db container and get pool connection
	dbContainer, db, err := SetupTestDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer dbContainer.Terminate(context.Background())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &UserStorage{
				db:  db,
				log: zap.NewExample(),
			}

			_, err = r.GetUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserStorage_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user model.User
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"success", args{context.Background(),
			model.User{
				FIO:      "Test User",
				Email:    "newemail@mail.ru",
				Password: "sad"}}, 2, false,
		},
		{
			"email exists", args{context.Background(),
				model.User{
					FIO:      "User Exists",
					Email:    "existsemail@mail.ru",
					Password: "asd",
				}}, 0, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// init logger and up test db container
			l, _ := logger.Init(&config.Config{Level: "dev"})

			dbContainer, db, err := SetupTestDatabase()
			if err != nil {
				log.Fatal(err)
			}
			defer dbContainer.Terminate(context.Background())

			r := &UserStorage{
				db:  db,
				log: l,
			}
			got, err := r.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func SetupTestDatabase() (testcontainers.Container, *sqlx.DB, error) {
	// Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "test_db",
			"POSTGRES_PASSWORD": "qwerty",
			"POSTGRES_USER":     "onelab",
		},
	}

	// Start PostgreSQL container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}

	// Get host and port of PostgreSQL container
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	// Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://onelab:qwerty@%v:%v/test_db", host, port.Port())

	db, err := sqlx.Connect("pgx", dbURI)
	if err != nil {
		return nil, nil, err
	}

	qr, err := os.ReadFile("./migrations_test/000001_init.up.sql")
	if err != nil {
		return dbContainer, db, err
	}

	if _, err := db.Exec(string(qr)); err != nil {
		return dbContainer, db, err
	}

	return dbContainer, db, err
}
