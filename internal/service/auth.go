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

	"github.com/golang-jwt/jwt/v5"
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
		IsAdmin:  false,
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		l.Error("Error occurred when creating a user", zap.Error(err))

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

	// Get the user from the database
	user, err := s.repo.GetUser(ctx, email, password)
	if err != nil {
		l.Warn("Error when getting a user", zap.Error(err))
		return &TokenDetails{}, err
	}

	// Create the access token with the user ID as the subject
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&accessTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authConfig.JWT.AccessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.ID,
			user.Email,
		})

	// Create the refresh token with the user ID as the subject
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&refreshTokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authConfig.JWT.RefreshTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			user.ID,
		})

	// Sign the tokens using the secret key
	accessTokenJWT, err := accessToken.SignedString([]byte(s.authConfig.JWT.AccessTokenSecret))
	if err != nil {
		return nil, err
	}

	refreshTokenJWT, err := refreshToken.SignedString([]byte(s.authConfig.JWT.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	// Store the tokens in the token details struct and return it
	tokenDetails.UserID = user.ID
	tokenDetails.AccessToken = accessTokenJWT
	tokenDetails.RefreshToken = refreshTokenJWT
	tokenDetails.AccessTokenExpiresIn = s.authConfig.JWT.AccessTokenTTL
	tokenDetails.RefreshTokenExpiresIn = s.authConfig.JWT.RefreshTokenTTL

	return &tokenDetails, nil
}

// ParseToken parses the access token and returns the user ID
func (s *AuthService) ParseToken(accessToken string) (string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(accessToken, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used to sign the token (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.ErrBadSigningMethod
		}

		// Return the secret key used to sign the token
		return []byte(s.authConfig.JWT.AccessTokenSecret), nil
	})

	// Check if the token is valid
	if err != nil {
		return "", err
	}

	// Check if the token is expired or not active yet
	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return "", apperror.ErrBadClaimsType
	}

	return claims.UserId, nil
}

// SaveRefreshToken parses the refresh token and returns the user ID
func (s *AuthService) SaveRefreshToken(ctx context.Context, userID string, details *TokenDetails) error {
	// Save the refresh token in the database with the user ID as the key and the refresh token as the value
	// The refresh token is stored with an expiration time equal to the refresh token TTL
	// This is so that we can delete the refresh token from the database when it expires
	// The refresh token is also stored in the token details struct so that it can be returned to the client
	err := s.redisClient.Set(ctx, userID, details.RefreshToken, details.RefreshTokenExpiresIn).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteRefreshToken deletes the refresh token from the database
func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID string) error {
	// Delete the refresh token from the database
	// The refresh token is stored in the database with the user ID as the key and the refresh token as the value
	// When the user logs out, we delete the refresh token from the database
	// This way, the user will not be able to use the refresh token to get a new access token
	// The user will have to log in again to get a new refresh token
	err := s.redisClient.Del(ctx, userID).Err()
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken generates a new access token and refresh token pair
func generatePasswordHash(password string) (string, error) {
	// Generate a hash from the password using the bcrypt algorithm with the default cost (10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Return the hash as a string
	return string(bytes), nil
}
