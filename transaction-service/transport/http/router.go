package http

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func (s *Server) SetUpRoute() {
	v1 := s.App.Group("/api/v1")

	s.App.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	s.App.GET("/live", func(e echo.Context) error {
		return e.NoContent(http.StatusOK)
	})

	// transaction endpoints
	transaction := v1.Group("/transactions")

	transaction.POST("", s.handler.CreateTransaction)
	transaction.POST("/items", s.handler.CreateTransactionItem)
	transaction.DELETE("/:id", s.handler.DeleteTransaction)
}
