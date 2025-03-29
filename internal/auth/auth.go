package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakubwoy/go-file-share/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWTToken(userID int, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(cfg.JWTExpiration)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}
func GenerateRandomString(length int) string {
	b := make([]byte, length/2+1)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)[:length]
}
func ValidateJWTToken(tokenString string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
