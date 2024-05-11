package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

// Claims represents the structure of the JWT claims.
type Claims struct {
	UserID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token based on the provided username and user ID.
// It returns the signed token string or an error if the token generation fails.
func GenerateToken(username *string, userID *primitive.ObjectID) (string, error) {

	// Retrieve the JWT secret key from environment variables.
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	// Define the token claims.
	claims := Claims{
		UserID:   *userID,
		Username: *username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			Subject:   "authentication",
		},
	}

	// Create a new JWT token with the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken verifies the provided JWT token string.
// It returns the claims if the token is valid or an error otherwise.
func VerifyToken(tokenString *string) (*Claims, error) {

	// Retrieve the JWT secret key from environment variables.
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	// Parse and verify the token.
	token, err := jwt.ParseWithClaims(*tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	// Check for errors during token parsing or verification.
	if err != nil {
		fmt.Println("Error parsing token : ", token)
		return nil, err
	}

	// Extract the claims from the token.
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

// get userID and username (logged in) from token cookie/headers
func GetUserIDAndUsernameFromContext(c *gin.Context) (string, string, error) {
	// Get token from cookie/headers
	token, err := c.Cookie("token")
	if err != nil {
		token := c.GetHeader("Authorization")
		if token == "" {
			return "", "", err
		}

	}
	claims, err := VerifyToken(&token)
	if err != nil {
		return "", "", err
	}
	return claims.UserID.Hex(), claims.Username, nil
}
