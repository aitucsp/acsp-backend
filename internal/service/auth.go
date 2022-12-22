package service

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorestapi/internal/apperror"
	"gorestapi/internal/dto"
	"gorestapi/internal/logs"
	"gorestapi/internal/model"
	"gorestapi/internal/repository"

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

func (s *AuthService) CreateUser(ctx context.Context, userDto dto.CreateUser) (string, error) {
	logs.Log().Info("Creating a user...")

	user := model.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: generatePasswordHash(userDto.Password),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	logs.Log().Info("Generating a token...")

	user, err := s.repo.GetUser(ctx, email, password)
	if err != nil {
		logs.Log().Warn("Error happened: %s", err.Error())
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

func generatePasswordHash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logs.Log().Fatalf("Error occurred when generating a hash of a password")
	}

	return string(bytes)
}
