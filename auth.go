package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWT Claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt (server-side hashing)
// Note: Password comes already hashed from client (PBKDF2)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash verifies a password against its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "techstore",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	// Optionally store session in database for tracking/blacklisting
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}
	DB.Create(session)

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// getCurrentUser extracts user from JWT token in request
func getCurrentUser(r *http.Request) (*User, error) {
	// Try to get token from Authorization header first
	authHeader := r.Header.Get("Authorization")
	var tokenString string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Fallback to cookie for browser requests
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			return nil, err
		}
		tokenString = cookie.Value
	}

	// Validate JWT
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// Get user from database
	var user User
	result := DB.Where("id = ?", claims.UserID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// authMiddleware protects routes requiring authentication
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := getCurrentUser(r)
		if err != nil {
			// Check if it's an API request or browser request
			if r.Header.Get("Content-Type") == "application/json" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Store user in request context for handlers to use
		// For now, we'll just call the next handler
		// In production, you'd use context.WithValue
		next(w, r)
	}
}

// registerHandler handles user registration
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"` // Already hashed with PBKDF2 from client
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Check if user exists
	var existingUser User
	result := DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email already registered"})
		return
	}

	// Hash the already-hashed password with bcrypt (double hashing for extra security)
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		return
	}

	// Create user
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		CreatedAt:    time.Now(),
	}

	if err := DB.Create(user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Set cookie for browser
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   token,
		"user": map[string]string{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// loginHandler handles user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"` // Already hashed with PBKDF2 from client
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Find user
	var user User
	result := DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Verify password (client sent PBKDF2 hash, we check against bcrypt hash)
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Set cookie for browser
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   token,
		"user": map[string]string{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// logoutHandler handles user logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get token to blacklist it
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		// Delete session from database
		DB.Where("token = ?", cookie.Value).Delete(&Session{})
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "auth_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// profileHandler shows user profile
func profileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	tmpl.Execute(w, user)
}

// updatePasswordHandler allows users to change their password
func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := getCurrentUser(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"` // PBKDF2 hashed
		NewPassword     string `json:"new_password"`     // PBKDF2 hashed
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// Verify current password
	if !CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	newHashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Server error"})
		return
	}

	// Update password
	if err := DB.Model(user).Update("password_hash", newHashedPassword).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update password"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password updated successfully",
	})
}
