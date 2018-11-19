package util

import (
	"errors"
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
	ss, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return "",err
	}
	return ss,nil
}



// 获取token
func GetToken(tokenStr string) (*jwt.Token,error) {
	token, err := jwt.ParseWithClaims(tokenStr,&jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	return token,err
}

// 获取token
func GetClaims(tokenStr string) (*jwt.StandardClaims,error) {
	token, err := jwt.ParseWithClaims(tokenStr,&jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err!=nil{
		return nil,err
	}
	if !token.Valid{
		return nil,errors.New("非法claims")
	}

	if claims,ok:= token.Claims.(*jwt.StandardClaims);ok{
		return claims,nil
	}
	return nil,errors.New("非法claims")
}