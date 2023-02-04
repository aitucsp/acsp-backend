package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"acsp/internal/apperror"
	"acsp/internal/config"
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

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
}

type TokenDetails struct {
	UserID                string        `json:"-"`
	AccessToken           string        `json:"access_token"`
	RefreshToken          string        `json:"refresh_token"`
	AccessTokenExpiresIn  time.Duration `json:"-"`
	RefreshTokenExpiresIn time.Duration `json:"-"`
}

type AuthService struct {
	repo        repository.Authorization
	roles       repository.Roles
	redisClient *redis.Client
	authConfig  config.AuthConfig
}

func NewAuthService(repo repository.Authorization, rolesRepo repository.Roles, r *redis.Client, a config.AuthConfig) *AuthService {
	return &AuthService{
		repo:        repo,
		roles:       rolesRepo,
		redisClient: r,
		authConfig:  a,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, userDto dto.CreateUser) error {
	l := logging.LoggerFromContext(ctx)
	l.Info("Creating a user...")

	user, err := s.repo.GetByEmail(ctx, userDto.Email)
	if user != nil || err == nil {
		return apperror.ErrEmailAlreadyExists
	}

	generatedHash, err := generatePasswordHash(userDto.Password)
	if err != nil {
		l.Error("Error occurred when generating hash password", zap.Error(err))

		return err
	}

	newUser := model.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: generatedHash,
	}

	userID, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		l.Error("Error occurred when creating a user", zap.Error(err))

		return err
	}

	err = s.roles.SaveUserRole(ctx, userID, 1)
	if err != nil {
		l.Error("Error occurred when adding a role to a user", zap.Error(err))

		return err
	}

	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	userId, _ := strconv.Atoi(userID)

	user, err := s.repo.GetByID(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s *AuthService) GenerateTokenPair(ctx context.Context, email, password string) (*TokenDetails, error) {
	l := logging.LoggerFromContext(ctx)
	l.Info("Generating a token...")

	var tokenDetails TokenDetails

	user, err := s.repo.GetUser(ctx, email, password)
	if err != nil {
		l.Warn("Error when getting a user", zap.Error(err))
		return &TokenDetails{}, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&accessTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authConfig.JWT.AccessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.ID,
			user.Email,
		})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&refreshTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authConfig.JWT.RefreshTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.ID,
		})

	accessTokenJWT, err := accessToken.SignedString([]byte(s.authConfig.JWT.AccessTokenSecret))
	if err != nil {
		return nil, err
	}

	refreshTokenJWT, err := refreshToken.SignedString([]byte(s.authConfig.JWT.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	tokenDetails.UserID = user.ID
	tokenDetails.AccessToken = accessTokenJWT
	tokenDetails.RefreshToken = refreshTokenJWT
	tokenDetails.AccessTokenExpiresIn = time.Minute * s.authConfig.JWT.AccessTokenTTL
	tokenDetails.RefreshTokenExpiresIn = time.Minute * s.authConfig.JWT.RefreshTokenTTL

	return &tokenDetails, nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrBadSigningMethod
		}

		return []byte(s.authConfig.JWT.AccessTokenSecret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return "", apperror.ErrBadClaimsType
	}

	return claims.UserId, nil
}

func (s *AuthService) SaveRefreshToken(ctx context.Context, userID string, details *TokenDetails) error {
	err := s.redisClient.Set(ctx, userID, details.RefreshToken, details.RefreshTokenExpiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID string) error {
	err := s.redisClient.Del(ctx, userID).Err()
	if err != nil {
		return err
	}

	return nil
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
