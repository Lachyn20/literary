package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hemra-siirow/literary/internal/infrastructure/config"
	"github.com/hemra-siirow/literary/internal/domain/repository"
)

type JWTAdapter struct {
	secret              string
	accessTokenExpiry   time.Duration
	refreshTokenExpiry  time.Duration
}

func NewJWTAdapter(cfg *config.Config) *JWTAdapter {
	return &JWTAdapter{
		secret:             cfg.JWTSecret,
		accessTokenExpiry:  cfg.JWTAccessTokenExpiry,
		refreshTokenExpiry: cfg.JWTRefreshTokenExpiry,
	}
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTAdapter) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := jwtClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTAdapter) GenerateRefreshToken() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (j *JWTAdapter) ValidateToken(token string) (*repository.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*jwtClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	return repository.NewTokenClaims(userID, claims.Role), nil
}
