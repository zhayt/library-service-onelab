package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/service/mocks"
	"go.uber.org/zap"
	"regexp"
	"testing"
)

func TestMatchesPatternTableDriven(t *testing.T) {
	type args struct {
		value   string
		pattern *regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Valid email", args{"example@email.ru", EmailRX}, false},
		{"Invalid email missing .", args{"example@emailru", EmailRX}, true},
		{"Invalid email missing @", args{"exampleemailru", EmailRX}, true},
		{"Invalid email has only домен", args{"@mail.ru", EmailRX}, true},
		{"Valid email", args{"example@gemail.com", EmailRX}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := matchesPattern(tt.args.value, tt.args.pattern); (err != nil) != tt.wantErr {
				t.Errorf("matchesPattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
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
		{name: "success", args: args{ctx: nil, user: model.User{FIO: "Aybek", Email: "example@mail.ru", Password: "asd"}}, want: 2, wantErr: false},
		{name: "email already exists", args: args{ctx: nil, user: model.User{FIO: "Test User", Email: "existexample@mail.ru", Password: "asd"}}, want: 0, wantErr: true},
		{name: "short name", args: args{ctx: nil, user: model.User{FIO: "sh", Email: "example@email.ru", Password: "asd"}}, want: 0, wantErr: true},
		{name: "invalid email", args: args{ctx: nil, user: model.User{FIO: "I don't have email", Email: "invalid@mailru", Password: "sad"}}, want: 0, wantErr: true},
	}

	l, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userStorage := mocks.NewIUserStorage(t)

			if tt.args.user.Email == "existexample@mail.ru" {
				userStorage.
					On("CreateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (int, error) {
						return 0, errors.New("user with this email already exists")
					})
			} else {
				userStorage.
					On("CreateUser", mock.Anything, mock.AnythingOfType("model.User")).
					Return(func(ctx context.Context, user model.User) (int, error) {
						return 2, nil
					})
			}

			s := &UserService{
				user: userStorage,
				log:  l,
			}
			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
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
