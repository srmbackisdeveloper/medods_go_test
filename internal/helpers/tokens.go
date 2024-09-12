package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtKey = os.Getenv("JWT_SECRET")
)

func GenerateAccessToken(userID string) (string, error) {
	// SHA512 Jwt
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(refreshToken), nil
}

func HashRefreshToken(refreshToken string) (string, error) {
	//hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	//if err != nil {
	//	return "", err
	//}

	return refreshToken, nil
}

func VerifyAccessToken(accessToken string) (string, error) {
	// Parse and validate the JWT token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	// Check if the token is valid and extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				return "", errors.New("token expired")
			}
		}

		// Extract user ID from the token claims (assumed to be in "sub")
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}

		return "", errors.New("user ID not found in token")
	}

	return "", errors.New("invalid token")
}
