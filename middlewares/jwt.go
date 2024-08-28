package middlewares

import (
    "errors"
    "net/http"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/sirupsen/logrus"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        logrus.Fatalf("failed to load configuration: %v", err)
    }
}

func JWTMiddleware(verifyToken bool) echo.MiddlewareFunc {
    return middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte(os.Getenv("JWT_SECRET")),
        Skipper: func(c echo.Context) bool {
            if verifyToken {
                email, err := ExtractVerifyToken(c)
                if err != nil {
                    c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid verification token"})
                    return true
                }
                c.Set("email", email)
            } else {
                id, role, err := ExtractToken(c)
                if err != nil {
                    c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
                    return true
                }
                c.Set("id", id)
                c.Set("role", role)
            }
            return false
        },
    })
}

func GenerateToken(id string, role string) (string, error) {
    logrus.Infof("generating token for user with ID: %s, Role: %s", id, role)

    tokenClaims := jwt.MapClaims{}
    tokenClaims["authorized"] = true
    tokenClaims["id"] = id
    tokenClaims["role"] = role
    tokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        logrus.Errorf("error generating token: %v", err)
        return "", err
    }

    logrus.Infof("token generated successfully: %s", tokenString)
    return tokenString, nil
}

func ExtractToken(c echo.Context) (string, string, error) {
    tokenString := c.Request().Header.Get("Authorization")
    if tokenString == "" {
        return "", "", errors.New("missing authorization token")
    }

    tokenString = tokenString[len("Bearer "):]
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        jwtSecret := []byte(os.Getenv("JWT_SECRET"))
        return jwtSecret, nil
    })

    if err != nil || !token.Valid {
        return "", "", errors.New("invalid authorization token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", "", errors.New("invalid authorization token")
    }

    id, ok := claims["id"].(string)
    if !ok {
        return "", "", errors.New("id claim not found in token")
    }

    role, ok := claims["role"].(string)
    if !ok {
        return "", "", errors.New("role claim not found in token")
    }

    logrus.Infof("Extracted ID: %s, Role: %s from token", id, role)
    return id, role, nil
}

func GenerateVerifyToken(email string) (string, error) {
    godotenv.Load()
    claims := jwt.MapClaims{}
    claims["email"] = email
    claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ExtractVerifyToken(c echo.Context) (string, error) {
    tokenString := c.Request().Header.Get("Authorization")
    if tokenString == "" {
        return "", errors.New("missing authorization token")
    }

    tokenString = tokenString[len("Bearer "):]
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        jwtSecret := []byte(os.Getenv("JWT_SECRET"))
        return jwtSecret, nil
    })

    if err != nil || !token.Valid {
        return "", errors.New("invalid authorization token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("invalid token claims")
    }

    email, ok := claims["email"].(string)
    if !ok {
        return "", errors.New("email claim not found in token")
    }

    return email, nil
}
