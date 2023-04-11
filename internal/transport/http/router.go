package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) SetUpRoute() {
	v1 := s.App.Group("/api/v1")

	v1.GET("/live", func(e echo.Context) error {
		return e.NoContent(http.StatusOK)
	})

	user := v1.Group("/users")
	user.POST("/sign-up", s.handler.CreateUser)
	user.POST("/sign-in", s.handler.GetUser)

	setting := user.Group("/settings", s.mid.ValidateAuth)
	setting.POST("/profile", s.handler.UpdateUser)
	setting.DELETE("/profile", s.handler.DeleteUser)
	setting.POST("/password", s.handler.UpdateUserPassword)

	book := v1.Group("/books")
	book.POST("/", s.handler.CreateBook)
}
