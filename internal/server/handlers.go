package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"server/internal/helpers"
	"server/internal/models"
	"strings"
	"time"
)

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ApiResponse{Message: s.db.Health()})
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func (s *Server) createAccountHandler(c *gin.Context) {
	var req struct {
		Email  string `json:"email" binding:"required,email"`
		Number string `json:"number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Error: helpers.WrongRequest})
		return
	}

	existingAccount, err := s.db.GetAccountByEmail(req.Email)
	if err == nil && existingAccount != nil {
		c.JSON(http.StatusConflict, ApiResponse{Error: "account already exists"})
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.WrongDB})
		return
	}

	account := models.Account{
		Email:  req.Email,
		Number: req.Number,
	}

	err = s.db.CreateAccount(&account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.DbError("create account failure")})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) getAccountHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: "missing authorization header"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: "invalid authorization header format"})
		return
	}
	accessToken := tokenParts[1]

	userID, err := helpers.VerifyAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: err.Error()})
		return
	}

	account, err := s.db.GetAccountByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: "failed to retrieve account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Tokens

func (s *Server) tokenHandler(c *gin.Context) {
	userID := c.Query("guid")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ApiResponse{Error: helpers.GuidRequired})
		return
	}

	accessToken, err := helpers.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.FailureAT})
		return
	}

	// Generate refresh token
	refreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.FailureRT})
		return
	}

	// Hash the refresh token
	hashedRefreshToken, err := helpers.HashRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: "token hashing failure"})
		return
	}

	// Save or update the refresh token in the database
	err = s.db.SaveRefreshToken(userID, hashedRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.DbError("failed to save refresh token")})
		return
	}

	c.JSON(http.StatusOK, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *Server) refreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Error: helpers.WrongRequest})
		return
	}

	// Fetch token info from the database using the user_id linked to the refresh token
	tokenInfo, err := s.db.GetUserByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: err.Error()})
		return
	}

	// ID
	uid := tokenInfo.UserID

	// Check if the refresh token is expired (older than 1 week)
	if tokenInfo.CreatedAt.Before(time.Now().AddDate(0, 0, -7)) {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: "refresh token expired"})
		return
	}

	// Compare the provided refresh token with the hashed refresh token from the database
	//err = bcrypt.CompareHashAndPassword([]byte(tokenInfo.HashedRefreshToken), []byte(req.RefreshToken))
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, ApiResponse{Error: helpers.WrongRT})
	//	return
	//}

	if tokenInfo.HashedRefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, ApiResponse{Error: helpers.WrongRT})
		return
	}

	// If the refresh token is valid, generate new tokens
	newAccessToken, err := helpers.GenerateAccessToken(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.FailureAT})
		return
	}

	// Generate new refresh token
	newRefreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.FailureRT})
		return
	}

	// Hash the new refresh token
	hashedNewRefreshToken, err := helpers.HashRefreshToken(newRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: "token hashing failure"})
		return
	}

	// Update the stored refresh token with the new hashed refresh token
	err = s.db.UpdateRefreshToken(uid, hashedNewRefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Error: helpers.DbError("update refresh token failure")})
		return
	}

	// Respond with new tokens
	c.JSON(http.StatusOK, Tokens{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
