package auth

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserCredentials represents the login credentials provided by the user.
type UserCredentials struct {
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

// Claims represents the structure of the JWT claims.
type Claims struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string
	jwt.StandardClaims
}
