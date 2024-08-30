package auth

// UserCredentials represents the login credentials provided by the user.
type UserCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
