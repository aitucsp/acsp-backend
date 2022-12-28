package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/model"
	"acsp/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 24 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, userDto dto.CreateUser) (int, error) {
	l := logging.LoggerFromContext(ctx)
	l.Info("Creating a user...")

	generatedHash, err := generatePasswordHash(userDto.Password)
	if err != nil {
		l.Error("Error occurred when generating hash password", zap.Error(err))
		return -1, err
	}

	user := model.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: generatedHash,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	l := logging.LoggerFromContext(ctx)
	l.Info("Generating a token...")

	user, err := s.repo.GetUser(ctx, email, password)
	if err != nil {
		l.Warn("Error when getting a user", zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.ID,
		user.Email,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrBadSigningMethod
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", apperror.ErrBadClaimsType
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
