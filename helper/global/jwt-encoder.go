package global

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)


func GenerateJwt() (string, error) {
	jwtSecret := viper.GetString("media-service.jwt-key")
	var sampleSecretKey = []byte(jwtSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	//claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["username"] = "rnr-svc"
	claims["fullname"] = "rnr svc"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}