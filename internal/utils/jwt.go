package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	error2 "github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/spf13/viper"
	"time"
)

type JwtCustomClaims struct {
	ID   uuid.UUID
	Role string
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID          uuid.UUID
	Role        string
	AccessToken string
	jwt.RegisteredClaims
}

func CreateAccessToken(user_id uuid.UUID, role string) (accessToken string, custErr error2.CustomError) {
	claims := &JwtCustomClaims{
		ID:   user_id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("EXPIRY")))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("SECRET_JWT")))
	if err != nil {
		custErr = error2.CustomError{
			Cause:   err,
			Service: "Utils",
			Code:    INTERNAL_SERVER_ERROR,
		}
		return "", custErr
	}
	return t, custErr
}

func CreateRefreshToken(user_id uuid.UUID, accessToken string, role string) (refreshToken string, custErr error2.CustomError) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID:          user_id,
		Role:        role,
		AccessToken: accessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(viper.GetInt("EXPIRY")))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(viper.GetString("SECRET_JWT")))
	if err != nil {
		custErr = error2.CustomError{
			Cause:   err,
			Service: "Utils",
			Code:    INTERNAL_SERVER_ERROR,
		}
		return "", custErr
	}
	return rt, custErr
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["id"].(string), nil
}
