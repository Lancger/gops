package util

import (
	"gops/backend/pkg/setting"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	NickName string `json:"nickname"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(username, nickname string) (string, error) {
	nowTime := time.Now()
	expriedTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		username,
		nickname,
		jwt.StandardClaims{
			ExpiresAt: expriedTime.Unix(),
			Issuer:    "jwt-go",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析Token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// // GetUsernameByToken 解析Token
func GetUsernameByToken(token string) (username, nickname string, err error) {
	username = ""
	nickname = ""
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			username = claims.Username
			nickname = claims.NickName
			err = nil
			return
		}
	}
	return
}
