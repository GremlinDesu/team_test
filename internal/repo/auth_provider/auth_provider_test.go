package authprovider_test

import (
	"errors"
	authprovider "team-test/internal/repo/auth_provider"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testingSecret = []byte("testing")
var testingTime = 5 * time.Minute
var testingRef = time.Hour

func TestJwt(t *testing.T) {
	var idu int64 = 124434234242
	Auth := authprovider.AuthProvider{Key: testingSecret, JwtExp: testingTime, RefreshExp: testingRef}
	token, _, err := Auth.CreateJwt(idu)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	id, err := Auth.ValidateJwt(token)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	if id != idu {
		t.Errorf("bad id")
	}
}

func TestExpiration(t *testing.T) {
	var idu int64 = 124434234242
	Auth := authprovider.AuthProvider{Key: testingSecret, JwtExp: time.Second, RefreshExp: time.Hour}
	token, ref, err := Auth.CreateJwt(idu)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	time.Sleep(5 * time.Second)
	_, err = Auth.ValidateJwt(token)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("bad expiration %v", err)
	}
	token, err = Auth.Refresh(ref)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
	_, err = Auth.ValidateJwt(token)
	if err != nil {
		t.Errorf("%v", err.Error())
	}
}

func TestHashing(t *testing.T) {
	Auth := authprovider.AuthProvider{}
	var password = "desuwa"
	hashed, err := Auth.HashPassword(password)
	if err != nil {
		t.Errorf("hashing error desuwa %v", err)
	}
	if !Auth.ComparePassword(password, hashed) {
		t.Errorf("password not equal")
	}
}
