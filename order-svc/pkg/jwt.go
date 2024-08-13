package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	User UserInfo
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Token struct {
	Token     string `json:"access_token"`
	ExpiresAt string `json:"expires_at"`
}

func GenerateToken(userInfo UserInfo) (Token, error) {
	expirationTime := time.Now().Add(200 * time.Minute)
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "golang-book.beneboba.me",
			Subject:   userInfo.ID,
		},
		User: UserInfo{
			ID:       userInfo.ID,
			Username: userInfo.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	accessToken, err := token.SignedString(signKey)

	if err != nil {
		return Token{}, err
	}

	return Token{
		Token:     accessToken,
		ExpiresAt: expirationTime.Format(time.DateTime),
	}, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
