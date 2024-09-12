package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile("etc/secrets/private.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadPublicKey() (*rsa.PublicKey, error) {
	publicKeyBytes, err := os.ReadFile("etc/secrets/public.key")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func GenerateToken(userID string) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func TokenValid(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)

	publicKey, err := loadPublicKey()
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return "", err
	}

	userID, err := extractTokenID(token)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func ExtractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func extractTokenID(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		sub, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("Invalid token")
		}

		return sub, nil
	}

	return "", nil
}
