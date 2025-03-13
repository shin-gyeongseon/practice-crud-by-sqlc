package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	minSecretKeySize = 32
)

type JWTMaker struct {
	scretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size : must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", Payload{}, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload) // Q. 이게 어떻게 토큰을 생성한다는걸 알아야할까 ? 
	token, err2 := jwtToken.SignedString([]byte(maker.scretKey))
	if err2 != nil {
		return "", Payload{}, err2
	}

	return token, *payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrExpiredToken
		}
		return []byte(maker.scretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
