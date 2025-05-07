package utils

import (
        "crypto/rand"
        "encoding/base64"
        "log"

        "github.com/gorilla/sessions"
        "golang.org/x/crypto/bcrypt"
)

var (
        // Store for sessions
        SessionStore = sessions.NewCookieStore([]byte("plant-exchange-secret-key"))
)

func init() {
        // Configure the session store
        SessionStore.Options = &sessions.Options{
                Path:     "/",
                MaxAge:   3600 * 24 * 30, // 30 days
                HttpOnly: true,
        }
        log.Println("Session store initialized with cookie options")
}

// GenerateToken generates a random token
func GenerateToken() (string, error) {
        b := make([]byte, 32)
        _, err := rand.Read(b)
        if err != nil {
                return "", err
        }
        return base64.URLEncoding.EncodeToString(b), nil
}

// HashPassword creates a hash of a password using bcrypt
func HashPassword(password string) string {
        hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                log.Printf("Error hashing password: %v", err)
                return ""
        }
        return string(hash)
}

// CheckPassword verifies a password against a hash
func CheckPassword(password string, hash string) bool {
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        return err == nil
}

// CreateSession creates a new session for a user
func CreateSession(userID string) string {
        token, _ := GenerateToken()
        // In a real implementation, you'd store this token in a database
        // with the user ID and expiration time
        return token
}

// ValidateSession validates a session token
func ValidateSession(token string) bool {
        // In a real implementation, you'd check if the token exists
        // and if it's not expired
        return token != ""
}

// GetUserIDFromToken gets the user ID associated with a token
func GetUserIDFromToken(token string) (string, bool) {
        // In a real implementation, you'd look up the user ID
        // in your database by the token
        // For now, this is a dummy implementation
        if token != "" {
                return "1", true // Return a dummy user ID
        }
        return "", false
}
