package database

import (
	"errors"
	"gorm.io/gorm"
	"server/internal/models"
)

func (s *service) SaveRefreshToken(userID, hashedRefreshToken string) error {
	var userToken models.UserToken

	err := s.db.Where("user_id = ?", userID).First(&userToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUserToken := models.UserToken{
				UserID:             userID,
				HashedRefreshToken: hashedRefreshToken,
			}
			return s.db.Create(&newUserToken).Error
		}
		return err
	}

	userToken.HashedRefreshToken = hashedRefreshToken
	return s.db.Save(&userToken).Error
}

func (s *service) GetRefreshToken(userID string) (string, error) {
	var userToken models.UserToken
	err := s.db.Where("user_id = ?", userID).First(&userToken).Error
	if err != nil {
		return "", err
	}

	return userToken.HashedRefreshToken, nil
}

func (s *service) UpdateRefreshToken(userID, newHashedToken string) error {
	var userToken models.UserToken
	err := s.db.Where("user_id = ?", userID).First(&userToken).Error
	if err != nil {
		return err
	}

	userToken.HashedRefreshToken = newHashedToken
	err = s.db.Save(&userToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateAccount(account *models.Account) error {
	err := s.db.Create(account).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAccountByEmail(email string) (*models.Account, error) {
	var account models.Account
	err := s.db.Where("email = ?", email).First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (s *service) GetUserByRefreshToken(refreshToken string) (*models.UserToken, error) {
	var userToken models.UserToken

	err := s.db.Where("hashed_refresh_token = ?", refreshToken).First(&userToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}

	return &userToken, nil
}

func (s *service) GetAccountByID(userId string) (*models.Account, error) {
	var account models.Account
	err := s.db.Where("id = ?", userId).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
