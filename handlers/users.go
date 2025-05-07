package handlers

import (
        "encoding/json"
        "net/http"

        "github.com/gorilla/mux"
        "github.com/plantexchange/app/utils"
)

// GetUser gets a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
        // Get user ID from URL path
        vars := mux.Vars(r)
        userID := vars["id"]

        // Find user
        user, exists := utils.GetUser(userID)
        if !exists {
                http.Error(w, "User not found", http.StatusNotFound)
                return
        }

        // Return user info (without sensitive data)
        userResponse := user.ToUserResponse()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
}

// UpdateUser updates a user's profile
func UpdateUser(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        currentUserID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get user ID from URL path
        vars := mux.Vars(r)
        userID := vars["id"]

        // Check if user is updating their own profile
        if currentUserID != userID {
                http.Error(w, "Unauthorized", http.StatusForbidden)
                return
        }

        // Find user
        user, exists := utils.GetUser(userID)
        if !exists {
                http.Error(w, "User not found", http.StatusNotFound)
                return
        }

        // Parse request
        var updates struct {
                Name       *string `json:"name"`
                Location   *string `json:"location"`
                Bio        *string `json:"bio"`
                ProfilePic *string `json:"profilePic"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Update fields if provided
        if updates.Name != nil {
                user.Name = *updates.Name
        }
        if updates.Location != nil {
                user.Location = *updates.Location
        }
        if updates.Bio != nil {
                user.Bio = *updates.Bio
        }
        if updates.ProfilePic != nil {
                user.ProfilePic = *updates.ProfilePic
        }

        // Save updated user
        utils.SaveUser(user)

        // Return updated user info
        userResponse := user.ToUserResponse()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
}

// GetUserListings gets listings by a user
func GetUserListings(w http.ResponseWriter, r *http.Request) {
        // Get user ID from URL path
        vars := mux.Vars(r)
        userID := vars["id"]

        // Find user
        _, exists := utils.GetUser(userID)
        if !exists {
                http.Error(w, "User not found", http.StatusNotFound)
                return
        }

        // Get listings by this user
        listings := utils.GetListingsByUser(userID)

        // Return listings
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(listings)
}
