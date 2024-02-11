package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func generateCustomerTokenPair(customerID uuid.UUID) (TokenPair, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return TokenPair{}, errors.New("jwt secret not found, cannot generate token pair")
	}

	accesToken := jwt.New(jwt.SigningMethodHS256)

	atClaims := accesToken.Claims.(jwt.MapClaims)
	atClaims["sub"] = customerID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	atSigned, err := accesToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = customerID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rtSigned, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  atSigned,
		RefreshToken: rtSigned,
	}, nil
}

func getCustomerIDFromToken(token string) (uuid.UUID, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return uuid.Nil, errors.New("jwt secret not found, cannot validate token")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid token")
	}

	customerID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, errors.New("invalid token")
	}

	return customerID, nil
}
