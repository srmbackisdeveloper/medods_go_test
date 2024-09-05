package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", s.healthHandler)
		v1.GET("/ping", s.pingHandler)
		v1.POST("/token", s.tokenHandler)
		v1.POST("/refresh-token", s.refreshTokenHandler)
	}

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (s *Server) tokenHandler(c *gin.Context) {
	tokens := Tokens{
		AccessToken:  "access t",
		RefreshToken: "refresh t",
	}

	c.JSON(http.StatusOK, tokens)
}

func (s *Server) refreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
