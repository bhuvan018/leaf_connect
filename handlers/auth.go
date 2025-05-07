package handlers

import (
        "bytes"
        "encoding/json"
        "io"
        "log"
        "net/http"
        "time"

        "github.com/plantexchange/app/models"
        "github.com/plantexchange/app/utils"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
        // Read the full request body for debugging
        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Error reading request body: " + err.Error()})
                return
        }
        
        // Log the request body for debugging
        log.Printf("Register request body: %s", string(bodyBytes))
        
        // Create a new reader with the same body contents
        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
        
        // Parse request body
        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body: " + err.Error()})
                return
        }
        
        // Log the parsed user for debugging
        log.Printf("Parsed user: %+v", user)

        // Validate required fields
        if user.Email == "" || user.Username == "" || user.Password == "" {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Email, username, and password are required"})
                return
        }

        // Check if email already exists
        if _, exists := utils.GetUserByEmail(user.Email); exists {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusConflict)
                json.NewEncoder(w).Encode(map[string]string{"message": "Email already registered"})
                return
        }

        // Check if username already exists
        if _, exists := utils.GetUserByUsername(user.Username); exists {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusConflict)
                json.NewEncoder(w).Encode(map[string]string{"message": "Username already taken"})
                return
        }

        // Hash password
        user.Password = utils.HashPassword(user.Password)

        // Set creation time
        user.CreatedAt = time.Now()
        user.LastLoginAt = time.Now()

        // Save user
        userID := utils.SaveUser(user)
        user.ID = userID

        // Create session
        session, _ := utils.SessionStore.Get(r, "session")
        session.Values["userID"] = user.ID
        session.Save(r, w)

        // Return user info (without password)
        userResponse := user.ToUserResponse()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
        log.Printf("Login attempt, cookies: %v", r.Cookies())
        
        // Parse request
        var credentials struct {
                Email    string `json:"email"`
                Password string `json:"password"`
        }
        
        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
                log.Printf("Error reading login request body: %v", err)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body: " + err.Error()})
                return
        }
        
        // Log the request body for debugging (strip password)
        log.Printf("Login request body: %s", string(bodyBytes))
        
        // Create a new reader with the same body contents
        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
        
        if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
                log.Printf("Error decoding login credentials: %v", err)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body: " + err.Error()})
                return
        }

        log.Printf("Login attempt for email: %s", credentials.Email)

        // Find user by email
        user, exists := utils.GetUserByEmail(credentials.Email)
        if !exists {
                log.Printf("User with email %s not found", credentials.Email)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
                return
        }

        log.Printf("Found user: %s (ID: %s)", user.Username, user.ID)

        // Check password
        if !utils.CheckPassword(credentials.Password, user.Password) {
                log.Printf("Invalid password for user: %s", user.Username)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or password"})
                return
        }

        log.Printf("Password validated for user: %s", user.Username)

        // Update last login time
        user.LastLoginAt = time.Now()
        utils.SaveUser(user)

        // Create session
        session, err := utils.SessionStore.Get(r, "session")
        if err != nil {
                log.Printf("Error getting session: %v", err)
        }
        
        log.Printf("Session before setting userID: %v", session.Values)
        session.Values["userID"] = user.ID
        
        err = session.Save(r, w)
        if err != nil {
                log.Printf("Error saving session: %v", err)
        }
        
        log.Printf("Session after setting userID: %v", session.Values)
        log.Printf("Set cookie header: %v", w.Header().Get("Set-Cookie"))

        // Return user info
        userResponse := user.ToUserResponse()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
        
        log.Printf("Login successful for user: %s", user.Username)
}

// Logout handles user logout
func Logout(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Clear session
        session.Values = make(map[interface{}]interface{})
        session.Save(r, w)

        // Return success
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// CheckAuth checks if a user is authenticated
func CheckAuth(w http.ResponseWriter, r *http.Request) {
        log.Printf("CheckAuth called, cookies: %v", r.Cookies())
        
        // Get current session
        session, err := utils.SessionStore.Get(r, "session")
        if err != nil {
                log.Printf("Error getting session: %v", err)
        }
        
        log.Printf("Session values: %v", session.Values)
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                log.Printf("No userID in session or not a string")
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
                return
        }

        log.Printf("Found userID in session: %s", userID)

        // Get user data
        user, exists := utils.GetUser(userID)
        if !exists {
                log.Printf("User with ID %s not found", userID)
                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
                return
        }

        log.Printf("User authenticated: %s (%s)", user.Username, user.ID)

        // User is authenticated
        response := map[string]interface{}{
                "authenticated": true,
                "user":          user.ToUserResponse(),
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(map[string]string{"message": "Not authenticated"})
                return
        }

        // Get user data
        user, exists := utils.GetUser(userID)
        if !exists {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"message": "User not found"})
                return
        }

        // Return user info
        userResponse := user.ToUserResponse()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
}
