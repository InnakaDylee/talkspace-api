package middlewares

import (
	"errors"
	"time"

	"talkspace-api/app/configs"
	"talkspace-api/utils/constant"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var config *configs.Configuration

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}
	config, err = configs.LoadConfig()
	if err != nil {
		logrus.Fatalf("failed to load configuration: %v", err)
	}
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(config.JWT.JWT_SECRET),
		SigningMethod: "HS256",
	})
}

func GenerateToken(id string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWT.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(c echo.Context) (string, string, error) {
	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		role := claims["role"].(string)

		return id, role, nil
	}
	return "", "", errors.New(constant.ERROR_TOKEN_INVALID)
}

func GenerateVerifyToken(email string) (string, error) {
	godotenv.Load()
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT.JWT_SECRET))
}

func ExtractVerifyToken(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		return email, nil
	}
	return "", errors.New(constant.ERROR_TOKEN_INVALID)
}
