package utility

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"test/internal/config"
	"test/internal/model"
	"time"
)

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type MyClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {
	claims := MyClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Cfg.JWT.Issuer,
			Subject:   "auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}
	return claims, nil
}

func IsTokenValid(claims *MyClaims, tokenString string) bool {
	if time.Now().Unix() > claims.ExpiresAt.Unix() {

		{
			fmt.Println("Token expired")
		}

		return false
	}
	return !model.IsTokenBlackListed(tokenString)
}

// DeleteToken Maintain a blacklist of tokens to revoke access
func DeleteToken(tokenString string) error {
	blackList := model.BlackList{
		Token: tokenString,
	}
	err := model.InsertBlackList(&blackList)
	if err != nil {
		return err
	}
	return nil
}
