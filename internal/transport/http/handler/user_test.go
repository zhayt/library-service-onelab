package handler

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/zhayt/user-storage-service/internal/model"
	"github.com/zhayt/user-storage-service/internal/transport/http/handler/mocks"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ShowUser(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		getUserByIDErr error
		expectedStatus int
	}{
		{
			name:           "Success",
			id:             "1",
			getUserByIDErr: nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User Not Found",
			id:             "2",
			getUserByIDErr: sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/users/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.id)

			userService := mocks.NewIUserService(t)

			if tc.id == "1" {
				userService.On("GetUserByID", mock.Anything, mock.AnythingOfType("int")).
					Return(func(ctx context.Context, userID int) (model.User, error) {
						return model.User{}, nil
					})
			} else {
				userService.On("GetUserByID", mock.Anything, mock.AnythingOfType("int")).
					Return(func(ctx context.Context, userID int) (model.User, error) {
						return model.User{}, sql.ErrNoRows
					})
			}

			h := &Handler{
				log:  zap.NewExample(),
				user: userService,
			}

			err := h.ShowUser(c)
			if err != nil {
				t.Fatalf("unexpected error: want %v, got %v", tc.getUserByIDErr, err)
			}

			if rec.Code != tc.expectedStatus {
				t.Errorf("unexpected status code: want %d, got %d", tc.expectedStatus, rec.Code)
			}
		})
	}
}
