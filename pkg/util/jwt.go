package util

import (
	"MIS/pkg/settings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(settings.AppSettings.JwtSecret)

type Claims struct {
	UserUuid string `json:"userUuid"`
	jwt.StandardClaims
}

// GenerateToken 生成jwt token
func GenerateToken(uuid string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour) // token过期时间3h

	claims := Claims{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			Issuer:    "Mis",             // 签发人
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret) // 用给定的加密方法和SecretKey对前面两部分加密，添在token的最后一段，至此token生成完毕

	return token, err
}

// ParseToken 解析jwt token
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
