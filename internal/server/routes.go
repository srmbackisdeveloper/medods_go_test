package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/account", s.createAccountHandler)
		v1.GET("/account", s.getAccountHandler)
		v1.GET("/health", s.healthHandler)
		v1.GET("/ping", s.pingHandler)
		v1.POST("/token", s.tokenHandler)
		v1.POST("/refresh-token", s.refreshTokenHandler)
	}

	return r
}
