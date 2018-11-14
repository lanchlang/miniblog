package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)
const JwtKey="that is a secret"
const JwtIssuer="company"
// 产生json web token
func GenToken(expiresDuration int64,subject string) (string,error) {
	claims := &jwt.StandardClaims{
		Subject:subject,
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + expiresDuration),
		Issuer:    JwtIssuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(JwtKey)
	if err != nil {
		return "",err
	}
	return ss,nil
}


// 获取token
func GetToken(tokenStr string) (*jwt.Token,error) {
	token, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	return token,err
}