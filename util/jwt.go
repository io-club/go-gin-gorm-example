package util

import "github.com/dgrijalva/jwt-go"

func GenerateToken(userID uint) (string, error) {
	// 生成 jwt
	jwtSecret := []byte("secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CheckToken(tokenString string) (uint, error) {
	jwtSecret := []byte("secret")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}
	userID := uint(claims["userID"].(float64))
	return userID, nil
}
