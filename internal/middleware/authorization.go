package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/LilLebowski/shortener/cmd/shortener/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type User struct {
	ID    string
	IsNew bool
}

func Authorization(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo, err := getUserIDFromCookie(ctx, config)
		if err != nil {
			code := http.StatusInternalServerError
			contentType := ctx.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				ctx.Header("Content-Type", "application/json")
				ctx.JSON(code, gin.H{
					"message": fmt.Sprintf("Unauthorized %s", err),
					"code":    code,
				})
			} else {
				ctx.String(code, fmt.Sprintf("Unauthorized %s", err))
			}
			ctx.Abort()
			return
		}
		ctx.Set("userID", userInfo.ID)
		ctx.Set("new", userInfo.IsNew)
	}
}

func getUserIDFromCookie(ctx *gin.Context, config *config.Config) (*User, error) {
	token, err := ctx.Cookie("userID")
	isNew := false
	if err != nil {
		token, err = buildJWTString(config)
		if err != nil {
			return nil, err
		}
		isNew = true
		ctx.SetCookie("userID", token, 3600, "/", "localhost", false, true)
	}
	userID, err := getUserID(token, config.SecretKey)
	if err != nil {
		return nil, err
	}
	userInfo := getNewUser(userID, isNew)

	return userInfo, nil
}

func buildJWTString(config *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpire)),
		},
		UserID: uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string, secretKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("token is not valid")
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

func getNewUser(id string, isNew bool) *User {
	return &User{
		ID:    id,
		IsNew: isNew,
	}
}
