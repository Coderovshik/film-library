package util

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnknownClaimsType = errors.New("unknown claims type, cannot proceed")
)

type JWTContextKey string

type UserClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func NewUserClaims(id int, username string, isAdmin bool) *UserClaims {
	return &UserClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "chat_app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
}

func NewJWTSignedString(claims jwt.Claims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseUserClaims(ss string, key string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(ss, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, ErrUnknownClaimsType
	}
}

func SetJWTCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		MaxAge:   int((24 * time.Hour).Seconds()),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	})
}
