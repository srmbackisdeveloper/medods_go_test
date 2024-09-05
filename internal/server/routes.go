package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/health", s.healthHandler)
	r.GET("/ping", s.pingHandler)
	r.POST("/token", s.tokenHandler)
	r.POST("/refresh-token", s.refreshTokenHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (s *Server) tokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) refreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
