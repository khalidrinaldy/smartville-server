package helper

import "github.com/dgrijalva/jwt-go"

func JwtGenerator(username, key string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return err.Error()
	}
	return tokenString
}