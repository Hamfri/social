package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func MockNewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, aud, iss}
}

// TRICK the compiler with a fake
// impl of JWTAuthenticator
type MockAuthenticator struct {
	secret string
	// aud    string
	// iss    string
}

func (m *MockAuthenticator) GenerateToken(claims jwt.Claims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(m.secret))
	return &tokenStr, nil
}

func (m *MockAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})

}
