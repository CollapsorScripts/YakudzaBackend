package token

import (
	"Yakudza/pkg/config"
	"Yakudza/pkg/database/models"
	"Yakudza/pkg/logger"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// tokenClaims - структура токена JWT
type tokenClaims struct {
	jwt.StandardClaims
	UserId uint64 `json:"userId"`
}

// IUser - интерфейс для доступа к полям субъекта токена JWT
type IUser interface {
	GetUserId() uint64
}

// GetUserId - возвращает ID пользователя из токена JWT
func (t *tokenClaims) GetUserId() uint64 {
	return t.UserId
}

// CreateToken - создание токена JWT
func CreateToken(user *models.User, cfg *config.Config) (string, error) {
	standartClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(cfg.JWT.Timeout).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	claims := &tokenClaims{
		StandardClaims: standartClaims,
		UserId:         user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		logger.Error("Ошибка при подписи токена: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ParseToken - парсит токен из строки
func ParseToken(accessToken string, cfg *config.Config) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprint("Неверная подпись"))
		}
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims, nil
}
