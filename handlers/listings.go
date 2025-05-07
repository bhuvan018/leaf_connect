package handlers

import (
        "encoding/json"
        "net/http"
        "strings"
        "time"

        "github.com/gorilla/mux"

        "github.com/plantexchange/app/models"
        "github.com/plantexchange/app/utils"
)

// GetListings returns all listings, with optional filtering
func GetListings(w http.ResponseWriter, r *http.Request) {
        // Get query parameters for filtering
        queryParams := r.URL.Query()
        userID := queryParams.Get("userId")
        listingType := queryParams.Get("type")
        plantType := queryParams.Get("plantType")
        location := queryParams.Get("location")

        // Get all listings
        allListings := utils.GetListings()
        
        // Filter listings based on query parameters
        filteredListings := []models.ListingWithUser{}
        
        for _, listing := range allListings {
                // Apply filters
                if userID != "" && listing.UserID != userID {
                        continue
                }
                if listingType != "" && listing.Type != listingType {
                        continue
                }
                if plantType != "" && listing.PlantType != plantType {
                        continue
                }
                if location != "" && !strings.Contains(strings.ToLower(listing.Location), strings.ToLower(location)) {
                        continue
                }
                
                // Get user info
                user, exists := utils.GetUser(listing.UserID)
                if !exists {
                        continue
                }
                
                // Create listing with user info
                listingWithUser := models.ListingWithUser{
                        Listing: listing,
                        User:    user.ToUserResponse(),
                }
                
                filteredListings = append(filteredListings, listingWithUser)
        }
        
        // Return filtered listings
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(filteredListings)
}

// GetListing returns a specific listing by ID
func GetListing(w http.ResponseWriter, r *http.Request) {
        // Get listing ID from URL path
        vars := mux.Vars(r)
        listingID := vars["id"]

        // Find listing
        listing, exists := utils.GetListing(listingID)
        if !exists {
                http.Error(w, "Listing not found", http.StatusNotFound)
                return
        }

        // Get user info
        user, exists := utils.GetUser(listing.UserID)
        if !exists {
                http.Error(w, "Listing owner not found", http.StatusInternalServerError)
                return
        }

        // Create listing with user info
        listingWithUser := models.ListingWithUser{
                Listing: listing,
                User:    user.ToUserResponse(),
        }

        // Return listing with user info
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(listingWithUser)
}

// CreateListing creates a new listing
func CreateListing(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Parse request
        var listing models.Listing
        if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Set user ID and timestamps
        listing.UserID = userID
        listing.CreatedAt = time.Now()
        listing.UpdatedAt = time.Now()
        
        // Set default status
        if listing.Status == "" {
                listing.Status = "available"
        }

        // Validate required fields
        if listing.Title == "" || listing.Description == "" || listing.Type == "" {
                http.Error(w, "Title, description, and type are required", http.StatusBadRequest)
                return
        }

        // Save listing
        listingID := utils.SaveListing(listing)
        listing.ID = listingID

        // Return created listing
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(listing)
}

// UpdateListing updates an existing listing
func UpdateListing(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get listing ID from URL path
        vars := mux.Vars(r)
        listingID := vars["id"]

        // Find listing
        listing, exists := utils.GetListing(listingID)
        if !exists {
                http.Error(w, "Listing not found", http.StatusNotFound)
                return
        }

        // Check if user owns the listing
        if listing.UserID != userID {
                http.Error(w, "Unauthorized", http.StatusForbidden)
                return
        }

        // Parse request
        var updates struct {
                Title       *string   `json:"title"`
                Description *string   `json:"description"`
                Type        *string   `json:"type"`
                PlantType   *string   `json:"plantType"`
                Price       *float64  `json:"price"`
                TradeFor    *string   `json:"tradeFor"`
                Location    *string   `json:"location"`
                Images      *[]string `json:"images"`
                Status      *string   `json:"status"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Update fields if provided
        if updates.Title != nil {
                listing.Title = *updates.Title
        }
        if updates.Description != nil {
                listing.Description = *updates.Description
        }
        if updates.Type != nil {
                listing.Type = *updates.Type
        }
        if updates.PlantType != nil {
                listing.PlantType = *updates.PlantType
        }
        if updates.Price != nil {
                listing.Price = *updates.Price
        }
        if updates.TradeFor != nil {
                listing.TradeFor = *updates.TradeFor
        }
        if updates.Location != nil {
                listing.Location = *updates.Location
        }
        if updates.Images != nil {
                listing.Images = *updates.Images
        }
        if updates.Status != nil {
                listing.Status = *updates.Status
        }

        // Update timestamp
        listing.UpdatedAt = time.Now()

        // Save updated listing
        utils.SaveListing(listing)

        // Return updated listing
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(listing)
}

// DeleteListing deletes a listing
func DeleteListing(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get listing ID from URL path
        vars := mux.Vars(r)
        listingID := vars["id"]

        // Find listing
        listing, exists := utils.GetListing(listingID)
        if !exists {
                http.Error(w, "Listing not found", http.StatusNotFound)
                return
        }

        // Check if user owns the listing
        if listing.UserID != userID {
                http.Error(w, "Unauthorized", http.StatusForbidden)
                return
        }

        // Delete listing
        if success := utils.DeleteListing(listingID); !success {
                http.Error(w, "Failed to delete listing", http.StatusInternalServerError)
                return
        }

        // Return success
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// SearchListings searches for listings based on query
func SearchListings(w http.ResponseWriter, r *http.Request) {
        // Get search query
        query := r.URL.Query().Get("q")
        if query == "" {
                http.Error(w, "Search query is required", http.StatusBadRequest)
                return
        }

        // Convert query to lowercase for case-insensitive search
        queryLower := strings.ToLower(query)

        // Get all listings
        allListings := utils.GetListings()
        
        // Filter listings based on search query
        searchResults := []models.ListingWithUser{}
        
        for _, listing := range allListings {
                // Check if query matches title, description, or plant type
                titleMatch := strings.Contains(strings.ToLower(listing.Title), queryLower)
                descMatch := strings.Contains(strings.ToLower(listing.Description), queryLower)
                typeMatch := strings.Contains(strings.ToLower(listing.PlantType), queryLower)
                
                if titleMatch || descMatch || typeMatch {
                        // Get user info
                        user, exists := utils.GetUser(listing.UserID)
                        if !exists {
                                continue
                        }
                        
                        // Create listing with user info
                        listingWithUser := models.ListingWithUser{
                                Listing: listing,
                                User:    user.ToUserResponse(),
                        }
                        
                        searchResults = append(searchResults, listingWithUser)
                }
        }
        
        // Return search results
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(searchResults)
}

// ToggleFavorite adds or removes a listing from a user's favorites
func ToggleFavorite(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Parse request
        var request struct {
                ListingID string `json:"listingId"`
                Action    string `json:"action"` // "add" or "remove"
        }
        
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Validate listing exists
        if _, exists := utils.GetListing(request.ListingID); !exists {
                http.Error(w, "Listing not found", http.StatusNotFound)
                return
        }

        var success bool
        if request.Action == "add" {
                success = utils.AddFavorite(userID, request.ListingID)
        } else if request.Action == "remove" {
                success = utils.RemoveFavorite(userID, request.ListingID)
        } else {
                http.Error(w, "Invalid action, must be 'add' or 'remove'", http.StatusBadRequest)
                return
        }

        // Return result
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

// GetFavorites gets a user's favorite listings
func GetFavorites(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get favorite listing IDs
        favoriteIDs := utils.GetFavorites(userID)
        
        // Get favorite listings
        favoriteListings := []models.ListingWithUser{}
        
        for _, id := range favoriteIDs {
                listing, exists := utils.GetListing(id)
                if !exists {
                        continue
                }
                
                // Get user info
                user, exists := utils.GetUser(listing.UserID)
                if !exists {
                        continue
                }
                
                // Create listing with user info
                listingWithUser := models.ListingWithUser{
                        Listing: listing,
                        User:    user.ToUserResponse(),
                }
                
                favoriteListings = append(favoriteListings, listingWithUser)
        }
        
        // Return favorites
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(favoriteListings)
}
