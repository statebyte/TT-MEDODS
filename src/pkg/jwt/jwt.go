package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTGenerator struct {
	AccessSecret   string
	RefreshSecret  string
	AccessExpires  time.Duration
	RefreshExpires time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Payload struct {
	UserID    uuid.UUID `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	jwt.RegisteredClaims
}

func (j *JWTGenerator) Init(secret string) {
	j.AccessSecret = secret
	j.RefreshSecret = secret
	j.AccessExpires = 15 * time.Minute
	j.RefreshExpires = 7 * 24 * time.Hour
}

func (j *JWTGenerator) generateToken(userID uuid.UUID, ipAddress string, expires time.Time, isAccessToken bool) (string, error) {
	var secret string
	if isAccessToken {
		secret = j.AccessSecret
	} else {
		secret = j.RefreshSecret
	}

	claims := Payload{
		UserID:    userID,
		IPAddress: ipAddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTGenerator) CreateTokenPair(userID uuid.UUID, ipAddress string) (*TokenPair, error) {
	tokenDetails := &TokenPair{}

	accessToken, err := j.generateToken(userID, ipAddress, time.Now().Add(j.AccessExpires), true)
	if err != nil {
		return nil, err
	}
	tokenDetails.AccessToken = accessToken

	refreshToken, err := j.generateToken(userID, ipAddress, time.Now().Add(j.RefreshExpires), false)
	if err != nil {
		return nil, err
	}
	tokenDetails.RefreshToken = refreshToken

	return tokenDetails, nil
}

func (j *JWTGenerator) VerifyToken(tokenString string, isAccessToken bool) (*Payload, error) {
	var secret string
	if isAccessToken {
		secret = j.AccessSecret
	} else {
		secret = j.RefreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Payload); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
