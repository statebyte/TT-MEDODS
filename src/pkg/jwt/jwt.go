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

type JWTResult struct {
	TokenUUID uuid.UUID
	Pair      TokenPair
}

type Payload struct {
	UserID    uuid.UUID `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	TokenUUID uuid.UUID
	jwt.RegisteredClaims
}

func (j *JWTGenerator) Init(secret string) {
	j.AccessSecret = secret
	j.RefreshSecret = secret
	j.AccessExpires = 15 * time.Minute
	j.RefreshExpires = 7 * 24 * time.Hour
}

func (j *JWTGenerator) generateToken(userID uuid.UUID, ipAddress string, expires time.Time, tokenUUID uuid.UUID, isAccessToken bool) (string, error) {
	var secret string
	if isAccessToken {
		secret = j.AccessSecret
	} else {
		secret = j.RefreshSecret
	}

	claims := Payload{
		UserID:    userID,
		IPAddress: ipAddress,
		TokenUUID: tokenUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	// HMAC - SHA512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTGenerator) CreateTokenPair(userID uuid.UUID, ipAddress string) (*JWTResult, error) {
	tokenDetails := &JWTResult{}

	tokenDetails.TokenUUID = uuid.New()

	accessToken, err := j.generateToken(userID, ipAddress, time.Now().Add(j.AccessExpires), tokenDetails.TokenUUID, true)
	if err != nil {
		return nil, err
	}
	tokenDetails.Pair.AccessToken = accessToken

	refreshToken, err := j.generateToken(userID, ipAddress, time.Now().Add(j.RefreshExpires), tokenDetails.TokenUUID, false)
	if err != nil {
		return nil, err
	}
	tokenDetails.Pair.RefreshToken = refreshToken

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

	if payload, ok := token.Claims.(*Payload); ok && token.Valid {
		return payload, nil
	}

	return nil, fmt.Errorf("invalid token")
}
