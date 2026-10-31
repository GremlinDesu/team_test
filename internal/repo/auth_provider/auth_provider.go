package authprovider

import (
	"fmt"
	"os"
	"team-test/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	Key        []byte
	JwtExp     time.Duration
	RefreshExp time.Duration
}

func NewAuthProvider(config config.Config) *AuthProvider {
	key := os.Getenv("AUTH_SECRET")
	// it was loaded in main.go, but to be more secure we don,t want to store it in config, so we load it here again
	if key == "" {
		panic("AUTH_SECRET environment variable is not set")
	}
	return &AuthProvider{
		Key:        []byte(key),
		JwtExp:     time.Duration(config.AuthJwtExpire) * time.Minute,
		RefreshExp: time.Duration(config.AuthRefreshExp) * time.Hour,
	}
}

func (A *AuthProvider) CreateJwt(id int64) (string, string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     now.Add(A.JwtExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	basicToken, err := token.SignedString(A.Key)
	if err != nil {
		return "", "", err
	}
	claimsref := jwt.MapClaims{
		"user_id": id,
		"exp":     now.Add(A.RefreshExp).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsref)
	refreshToken, err := refresh.SignedString(A.Key)
	return basicToken, refreshToken, nil
}

func (A *AuthProvider) ValidateJwt(tokenstr string) (int64, error) {
	parsed, err := jwt.Parse(tokenstr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {

			return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Header["alg"])

		}
		return A.Key, nil
	})
	if err != nil {

		return 0, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		idf := claims["user_id"].(float64)
		id := int64(idf)
		return id, nil
	} else {
		return 0, ErrInvalidToken
	}
}

func (A *AuthProvider) Refresh(Refresh string) (string, error) {
	id, err := A.ValidateJwt(Refresh)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(A.JwtExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(A.Key)
}

func (A *AuthProvider) ComparePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (A *AuthProvider) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
