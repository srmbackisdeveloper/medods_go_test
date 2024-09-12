package database

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"server/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() string
	SaveRefreshToken(userID, hashedRefreshToken string) error
	GetRefreshToken(userID string) (string, error)
	UpdateRefreshToken(userID, newHashedToken string) error
	CreateAccount(account *models.Account) error
	GetAccountByEmail(email string) (*models.Account, error)
	GetAccountByID(userId string) (*models.Account, error)
	GetUserByRefreshToken(refreshToken string) (*models.UserToken, error)
}

type service struct {
	db *gorm.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate(&models.Account{}, &models.UserToken{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return &service{db: db}
}

func (s *service) Health() string {
	sqlDB, err := s.db.DB()
	if err != nil {
		log.Fatalf("FATAL!\ndb connection failure: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return "Database is healthy!"
}
