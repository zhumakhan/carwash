package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
	"errors"
)

const WrongUUIDErrorText = "Wrong UUID structure"

type TokenClaim struct {
	UUID      string
	IssuedAt  int64
	ExpiresAt int64
}

type WrongUUID struct{}

func (wu WrongUUID) Error() string {
	return WrongUUIDErrorText
}


func (tc *TokenClaim) Valid() error {
	if !IsUUID(tc.UUID) {
		return WrongUUID{}
	}
	return nil
}


//Generates both tokens
func GenerateTokens(uuid string) (string, string) {
	now := time.Now()
	issuedAt := now.Unix()
	refreshExpiresAt := now.AddDate(0, 6, 0).Unix()
	accessExpiresAt := now.AddDate(0, 6, 0).Unix()
	accessTC := &TokenClaim{UUID: uuid, IssuedAt: issuedAt, ExpiresAt: accessExpiresAt}
	refreshTC := &TokenClaim{UUID: uuid, IssuedAt: issuedAt, ExpiresAt: refreshExpiresAt}
	accessToken, _ := GenerateToken(accessTC, os.Getenv("access_token_password"))
	refreshToken, _ := GenerateToken(refreshTC, os.Getenv("refresh_token_password"))
	return accessToken, refreshToken
}

//Generates one token according to token password and claims
func GenerateToken(tc *TokenClaim, tokenPass string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tc)
	tokenString, err := token.SignedString([]byte(tokenPass))
	return tokenString, err
}

//Checks correctness of token structure
func ValidateToken(tokenStr string, tokenPass string) (error, *TokenClaim) {
	tc := &TokenClaim{}
	_, err := jwt.ParseWithClaims(tokenStr, tc, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenPass), nil
	})
	if err != nil {
		return err, nil
	}
	tokenStr2, _ := GenerateToken(tc, tokenPass)
	if tokenStr != tokenStr2 {
		return errors.New("Fuck you"), nil
	}
	return nil, tc
}

func GetTokenFromString(tokenStr string, tokenPass string) *TokenClaim {
	tc := &TokenClaim{}
	jwt.ParseWithClaims(tokenStr, tc, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenPass), nil
	})
	return tc
}
