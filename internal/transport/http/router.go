package http

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"net/http"
)

func (s *Server) SetUpRoute() {
	v1 := s.App.Group("/api/v1")

	s.App.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	s.App.GET("/live", func(e echo.Context) error {
		return e.NoContent(http.StatusOK)
	})

	user := v1.Group("/users")
	user.POST("/sign-up", s.handler.SignUp)
	user.POST("/sign-in", s.handler.SignIn)
	user.GET("/:id", s.handler.ShowUser)

	setting := user.Group("/settings", s.mid.ValidateAuth)
	setting.PATCH("/profile", s.handler.UpdateUserFIO)
	setting.PATCH("/password", s.handler.UpdateUserPassword)
	setting.DELETE("/profile", s.handler.DeleteUser)

	book := v1.Group("/books")
	book.POST("", s.handler.CreateBook)
	book.GET("", s.handler.ShowAllBooks)
	book.GET("/:id", s.handler.ShowBook)
	book.PATCH("/:id", s.handler.UpdateBook)
	book.DELETE("/:id", s.handler.DeleteBook)

	history := v1.Group("/rents")
	history.POST("", s.handler.CreateBIHistory, s.mid.ValidateAuth)
	history.GET("", s.handler.ShowCurrentBorrowedBooks)
	history.GET("/months", s.handler.ShowBIHistoryLastMonth)
	history.PATCH("/:id", s.handler.UpdateBIHistory)
	history.DELETE("/:id", s.handler.DeleteBIHistory)
}
