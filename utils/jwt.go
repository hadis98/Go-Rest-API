package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecretkey"

func GenerateToken(email string, userId int64) (string, error) {
	fmt.Println("************[GenerateToken]************")
	fmt.Println("user id is : ", userId)
	fmt.Println("email is : ", email)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), //the lifetime of the token is 2hours
	})

	return token.SignedString([]byte(secretKey)) // we want to get a single string to attach it to the future requests
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// * we check weather the recieved token was signed by the signingMethod that we have used
		_, ok := t.Method.(*jwt.SigningMethodHMAC) //*this is a type checking syntax where we can access a type like a property on a field
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// * we can extract the data that we used in the process of token creation
	// email ,ok:= claims["email"].(string)
	// userId, ok := claims["userId"].(float64)
	// if !ok {
	// 	return 0, errors.New("invalid token claims")
	// }
	userId := int64(claims["userId"].(float64))
	return userId, nil
}
