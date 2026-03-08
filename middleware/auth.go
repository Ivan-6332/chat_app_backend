package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"chatapp-backend/config"
	"chatapp-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var jwksCache *JWKS
var jwksCacheTime time.Time

// Auth0Middleware validates JWT tokens from Auth0
func Auth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Authorization header is required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		token, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse(fmt.Sprintf("Invalid token: %v", err)))
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
			c.Set("user_id", claims["sub"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("Invalid token claims"))
			c.Abort()
			return
		}
	}
}

func validateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		// Get JWKS
		jwks, err := getJWKS()
		if err != nil {
			return nil, err
		}

		// Find matching key
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return convertKey(key)
			}
		}

		return nil, errors.New("unable to find appropriate key")
	})

	if err != nil {
		return nil, err
	}

	// Verify claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Verify audience
	aud, ok := claims["aud"].(string)
	if !ok {
		// Try as array
		if audArray, ok := claims["aud"].([]interface{}); ok {
			found := false
			for _, a := range audArray {
				if audStr, ok := a.(string); ok && audStr == config.AppConfig.Auth0Audience {
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("invalid audience")
			}
		} else {
			return nil, errors.New("invalid audience")
		}
	} else if aud != config.AppConfig.Auth0Audience {
		return nil, errors.New("invalid audience")
	}

	// Verify issuer
	iss := fmt.Sprintf("https://%s/", config.AppConfig.Auth0Domain)
	if issuer, ok := claims["iss"].(string); !ok || issuer != iss {
		return nil, errors.New("invalid issuer")
	}

	return token, nil
}

func getJWKS() (*JWKS, error) {
	// Check cache (cache for 1 hour)
	if jwksCache != nil && time.Since(jwksCacheTime) < time.Hour {
		return jwksCache, nil
	}

	// Fetch JWKS from Auth0
	url := fmt.Sprintf("https://%s/.well-known/jwks.json", config.AppConfig.Auth0Domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// Update cache
	jwksCache = &jwks
	jwksCacheTime = time.Now()

	return &jwks, nil
}

func convertKey(key JSONWebKey) (*rsa.PublicKey, error) {
	// Decode N
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}

	// Decode E
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	// Convert to big.Int
	n := new(big.Int).SetBytes(nBytes)

	// Convert E to int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}
