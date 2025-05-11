package Utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const signedKey = "supersecret"

func GenerateToken(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"expiry": time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(signedKey))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("not returning same secret key")
		}
		return []byte(signedKey), nil
	})
	if err != nil {
		return 0, errors.New("could not verify token")
	}
	tokenIsvalid := parsedToken.Valid

	if !tokenIsvalid {
		return 0, errors.New("token is invalid")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	userId := int64(claims["userID"].(float64))
	return userId, nil
}
