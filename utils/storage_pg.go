package utils

import (
        "database/sql"
        "log"
        "strconv"
        "time"

        "github.com/plantexchange/app/models"
)

// PostgreSQL storage implementation

// GetUsers retrieves all users from the database
func GetUsers() []models.User {
        rows, err := GetDB().Query(`
                SELECT id, email, username, password, name, location, bio, profile_pic, created_at, last_login_at
                FROM users
        `)
        if err != nil {
                log.Printf("Error getting users: %v", err)
                return []models.User{}
        }
        defer rows.Close()

        users := []models.User{}
        for rows.Next() {
                var user models.User
                var id int
                err := rows.Scan(&id, &user.Email, &user.Username, &user.Password, &user.Name, &user.Location, &user.Bio, &user.ProfilePic, &user.CreatedAt, &user.LastLoginAt)
                if err != nil {
                        log.Printf("Error scanning user row: %v", err)
                        continue
                }
                user.ID = strconv.Itoa(id)

                // Get user favorites
                user.Favorites = GetFavorites(user.ID)

                users = append(users, user)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating user rows: %v", err)
        }

        return users
}

// GetUser retrieves a user by ID from the database
func GetUser(id string) (models.User, bool) {
        var user models.User
        var dbID int

        userID, err := strconv.Atoi(id)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return models.User{}, false
        }

        err = GetDB().QueryRow(`
                SELECT id, email, username, password, name, location, bio, profile_pic, created_at, last_login_at
                FROM users
                WHERE id = $1
        `, userID).Scan(&dbID, &user.Email, &user.Username, &user.Password, &user.Name, &user.Location, &user.Bio, &user.ProfilePic, &user.CreatedAt, &user.LastLoginAt)

        if err != nil {
                if err == sql.ErrNoRows {
                        return models.User{}, false
                }
                log.Printf("Error getting user by ID: %v", err)
                return models.User{}, false
        }

        user.ID = strconv.Itoa(dbID)
        
        // Get user favorites
        user.Favorites = GetFavorites(user.ID)

        return user, true
}

// GetUserByEmail retrieves a user by email from the database
func GetUserByEmail(email string) (models.User, bool) {
        var user models.User
        var id int

        err := GetDB().QueryRow(`
                SELECT id, email, username, password, name, location, bio, profile_pic, created_at, last_login_at
                FROM users
                WHERE email = $1
        `, email).Scan(&id, &user.Email, &user.Username, &user.Password, &user.Name, &user.Location, &user.Bio, &user.ProfilePic, &user.CreatedAt, &user.LastLoginAt)

        if err != nil {
                if err == sql.ErrNoRows {
                        return models.User{}, false
                }
                log.Printf("Error getting user by email: %v", err)
                return models.User{}, false
        }

        user.ID = strconv.Itoa(id)
        
        // Get user favorites
        user.Favorites = GetFavorites(user.ID)

        return user, true
}

// GetUserByUsername retrieves a user by username from the database
func GetUserByUsername(username string) (models.User, bool) {
        var user models.User
        var id int

        err := GetDB().QueryRow(`
                SELECT id, email, username, password, name, location, bio, profile_pic, created_at, last_login_at
                FROM users
                WHERE username = $1
        `, username).Scan(&id, &user.Email, &user.Username, &user.Password, &user.Name, &user.Location, &user.Bio, &user.ProfilePic, &user.CreatedAt, &user.LastLoginAt)

        if err != nil {
                if err == sql.ErrNoRows {
                        return models.User{}, false
                }
                log.Printf("Error getting user by username: %v", err)
                return models.User{}, false
        }

        user.ID = strconv.Itoa(id)
        
        // Get user favorites
        user.Favorites = GetFavorites(user.ID)

        return user, true
}

// SaveUser saves a user to the database
func SaveUser(user models.User) string {
        // If the user has no ID, insert a new user
        if user.ID == "" {
                var id int
                err := GetDB().QueryRow(`
                        INSERT INTO users (email, username, password, name, location, bio, profile_pic, created_at, last_login_at)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                        RETURNING id
                `, user.Email, user.Username, user.Password, user.Name, user.Location, user.Bio, user.ProfilePic, user.CreatedAt, user.LastLoginAt).Scan(&id)

                if err != nil {
                        log.Printf("Error creating user: %v", err)
                        return ""
                }

                return strconv.Itoa(id)
        }

        // User has an ID, update existing user
        userID, err := strconv.Atoi(user.ID)
        if err != nil {
                log.Printf("Invalid user ID for update: %v", err)
                return ""
        }

        _, err = GetDB().Exec(`
                UPDATE users
                SET email = $1, username = $2, password = $3, name = $4, location = $5, bio = $6, profile_pic = $7, last_login_at = $8
                WHERE id = $9
        `, user.Email, user.Username, user.Password, user.Name, user.Location, user.Bio, user.ProfilePic, user.LastLoginAt, userID)

        if err != nil {
                log.Printf("Error updating user: %v", err)
                return ""
        }

        return user.ID
}

// GetListings retrieves all listings from the database
func GetListings() []models.Listing {
        rows, err := GetDB().Query(`
                SELECT l.id, l.user_id, l.title, l.description, l.type, l.plant_type, l.price,
                           l.trade_for, l.location, l.created_at, l.updated_at, l.status
                FROM listings l
                ORDER BY l.created_at DESC
        `)
        if err != nil {
                log.Printf("Error getting listings: %v", err)
                return []models.Listing{}
        }
        defer rows.Close()

        listings := []models.Listing{}
        for rows.Next() {
                var listing models.Listing
                var id, userID int
                err := rows.Scan(&id, &userID, &listing.Title, &listing.Description, &listing.Type, &listing.PlantType, &listing.Price,
                        &listing.TradeFor, &listing.Location, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status)
                if err != nil {
                        log.Printf("Error scanning listing row: %v", err)
                        continue
                }
                listing.ID = strconv.Itoa(id)
                listing.UserID = strconv.Itoa(userID)

                // Get images for the listing
                images, err := getListingImages(id)
                if err != nil {
                        log.Printf("Error getting images for listing %d: %v", id, err)
                } else {
                        listing.Images = images
                }

                listings = append(listings, listing)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating listing rows: %v", err)
        }

        return listings
}

// getListingImages retrieves all images for a listing
func getListingImages(listingID int) ([]string, error) {
        rows, err := GetDB().Query(`
                SELECT image_url FROM listing_images
                WHERE listing_id = $1
                ORDER BY id
        `, listingID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var images []string
        for rows.Next() {
                var imageURL string
                if err := rows.Scan(&imageURL); err != nil {
                        return nil, err
                }
                images = append(images, imageURL)
        }

        if err = rows.Err(); err != nil {
                return nil, err
        }

        return images, nil
}

// GetListing retrieves a listing by ID from the database
func GetListing(id string) (models.Listing, bool) {
        var listing models.Listing
        var dbID, userID int

        listingID, err := strconv.Atoi(id)
        if err != nil {
                log.Printf("Invalid listing ID: %v", err)
                return models.Listing{}, false
        }

        err = GetDB().QueryRow(`
                SELECT id, user_id, title, description, type, plant_type, price,
                           trade_for, location, created_at, updated_at, status
                FROM listings
                WHERE id = $1
        `, listingID).Scan(&dbID, &userID, &listing.Title, &listing.Description, &listing.Type, &listing.PlantType, &listing.Price,
                &listing.TradeFor, &listing.Location, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status)

        if err != nil {
                if err == sql.ErrNoRows {
                        return models.Listing{}, false
                }
                log.Printf("Error getting listing by ID: %v", err)
                return models.Listing{}, false
        }

        listing.ID = strconv.Itoa(dbID)
        listing.UserID = strconv.Itoa(userID)

        // Get images for the listing
        images, err := getListingImages(dbID)
        if err != nil {
                log.Printf("Error getting images for listing %d: %v", dbID, err)
        } else {
                listing.Images = images
        }

        return listing, true
}

// GetListingsByUser retrieves all listings by a user from the database
func GetListingsByUser(userID string) []models.Listing {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return []models.Listing{}
        }

        rows, err := GetDB().Query(`
                SELECT l.id, l.user_id, l.title, l.description, l.type, l.plant_type, l.price,
                           l.trade_for, l.location, l.created_at, l.updated_at, l.status
                FROM listings l
                WHERE l.user_id = $1
                ORDER BY l.created_at DESC
        `, userIDInt)
        if err != nil {
                log.Printf("Error getting listings by user: %v", err)
                return []models.Listing{}
        }
        defer rows.Close()

        listings := []models.Listing{}
        for rows.Next() {
                var listing models.Listing
                var id, dbUserID int
                err := rows.Scan(&id, &dbUserID, &listing.Title, &listing.Description, &listing.Type, &listing.PlantType, &listing.Price,
                        &listing.TradeFor, &listing.Location, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status)
                if err != nil {
                        log.Printf("Error scanning listing row: %v", err)
                        continue
                }
                listing.ID = strconv.Itoa(id)
                listing.UserID = strconv.Itoa(dbUserID)

                // Get images for the listing
                images, err := getListingImages(id)
                if err != nil {
                        log.Printf("Error getting images for listing %d: %v", id, err)
                } else {
                        listing.Images = images
                }

                listings = append(listings, listing)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating listing rows: %v", err)
        }

        return listings
}

// SaveListing saves a listing to the database
func SaveListing(listing models.Listing) string {
        tx, err := GetDB().Begin()
        if err != nil {
                log.Printf("Error starting transaction: %v", err)
                return ""
        }
        defer func() {
                if err != nil {
                        tx.Rollback()
                }
        }()

        // If the listing has no ID, insert a new listing
        if listing.ID == "" {
                userID, err := strconv.Atoi(listing.UserID)
                if err != nil {
                        log.Printf("Invalid user ID: %v", err)
                        return ""
                }

                var id int
                err = tx.QueryRow(`
                        INSERT INTO listings (user_id, title, description, type, plant_type, price,
                                                                 trade_for, location, created_at, updated_at, status)
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                        RETURNING id
                `, userID, listing.Title, listing.Description, listing.Type, listing.PlantType, listing.Price,
                        listing.TradeFor, listing.Location, listing.CreatedAt, listing.UpdatedAt, listing.Status).Scan(&id)

                if err != nil {
                        log.Printf("Error creating listing: %v", err)
                        return ""
                }

                // Save images
                for _, imageURL := range listing.Images {
                        _, err = tx.Exec(`
                                INSERT INTO listing_images (listing_id, image_url)
                                VALUES ($1, $2)
                        `, id, imageURL)
                        if err != nil {
                                log.Printf("Error saving listing image: %v", err)
                                return ""
                        }
                }

                err = tx.Commit()
                if err != nil {
                        log.Printf("Error committing transaction: %v", err)
                        return ""
                }

                return strconv.Itoa(id)
        }

        // Listing has an ID, update existing listing
        listingID, err := strconv.Atoi(listing.ID)
        if err != nil {
                log.Printf("Invalid listing ID for update: %v", err)
                return ""
        }

        userID, err := strconv.Atoi(listing.UserID)
        if err != nil {
                log.Printf("Invalid user ID for update: %v", err)
                return ""
        }

        listing.UpdatedAt = time.Now()

        _, err = tx.Exec(`
                UPDATE listings
                SET user_id = $1, title = $2, description = $3, type = $4, plant_type = $5,
                        price = $6, trade_for = $7, location = $8, updated_at = $9, status = $10
                WHERE id = $11
        `, userID, listing.Title, listing.Description, listing.Type, listing.PlantType,
                listing.Price, listing.TradeFor, listing.Location, listing.UpdatedAt, listing.Status, listingID)

        if err != nil {
                log.Printf("Error updating listing: %v", err)
                return ""
        }

        // Delete existing images and add new ones
        _, err = tx.Exec(`DELETE FROM listing_images WHERE listing_id = $1`, listingID)
        if err != nil {
                log.Printf("Error deleting listing images: %v", err)
                return ""
        }

        for _, imageURL := range listing.Images {
                _, err = tx.Exec(`
                        INSERT INTO listing_images (listing_id, image_url)
                        VALUES ($1, $2)
                `, listingID, imageURL)
                if err != nil {
                        log.Printf("Error saving listing image: %v", err)
                        return ""
                }
        }

        err = tx.Commit()
        if err != nil {
                log.Printf("Error committing transaction: %v", err)
                return ""
        }

        return listing.ID
}

// DeleteListing deletes a listing from the database
func DeleteListing(id string) bool {
        listingID, err := strconv.Atoi(id)
        if err != nil {
                log.Printf("Invalid listing ID: %v", err)
                return false
        }

        // Start a transaction
        tx, err := GetDB().Begin()
        if err != nil {
                log.Printf("Error starting transaction: %v", err)
                return false
        }
        defer func() {
                if err != nil {
                        tx.Rollback()
                }
        }()

        // Delete images first (cascade will handle this, but being explicit)
        _, err = tx.Exec(`DELETE FROM listing_images WHERE listing_id = $1`, listingID)
        if err != nil {
                log.Printf("Error deleting listing images: %v", err)
                return false
        }

        // Delete the listing
        result, err := tx.Exec(`DELETE FROM listings WHERE id = $1`, listingID)
        if err != nil {
                log.Printf("Error deleting listing: %v", err)
                return false
        }

        // Check if the listing was actually deleted
        rowsAffected, err := result.RowsAffected()
        if err != nil {
                log.Printf("Error getting rows affected: %v", err)
                return false
        }

        // Commit the transaction
        err = tx.Commit()
        if err != nil {
                log.Printf("Error committing transaction: %v", err)
                return false
        }

        return rowsAffected > 0
}

// GetMessages retrieves all messages from the database
func GetMessages() []models.Message {
        rows, err := GetDB().Query(`
                SELECT id, from_id, to_id, listing_id, content, read, created_at
                FROM messages
                ORDER BY created_at
        `)
        if err != nil {
                log.Printf("Error getting messages: %v", err)
                return []models.Message{}
        }
        defer rows.Close()

        messages := []models.Message{}
        for rows.Next() {
                var message models.Message
                var id, fromID, toID int
                var listingID sql.NullInt64
                err := rows.Scan(&id, &fromID, &toID, &listingID, &message.Content, &message.Read, &message.CreatedAt)
                if err != nil {
                        log.Printf("Error scanning message row: %v", err)
                        continue
                }
                message.ID = strconv.Itoa(id)
                message.FromID = strconv.Itoa(fromID)
                message.ToID = strconv.Itoa(toID)
                if listingID.Valid {
                        message.ListingID = strconv.FormatInt(listingID.Int64, 10)
                } else {
                        message.ListingID = ""
                }

                messages = append(messages, message)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating message rows: %v", err)
        }

        return messages
}

// GetMessage retrieves a message by ID from the database
func GetMessage(id string) (models.Message, bool) {
        var message models.Message
        var dbID, fromID, toID int
        var listingID sql.NullInt64

        messageID, err := strconv.Atoi(id)
        if err != nil {
                log.Printf("Invalid message ID: %v", err)
                return models.Message{}, false
        }

        err = GetDB().QueryRow(`
                SELECT id, from_id, to_id, listing_id, content, read, created_at
                FROM messages
                WHERE id = $1
        `, messageID).Scan(&dbID, &fromID, &toID, &listingID, &message.Content, &message.Read, &message.CreatedAt)

        if err != nil {
                if err == sql.ErrNoRows {
                        return models.Message{}, false
                }
                log.Printf("Error getting message by ID: %v", err)
                return models.Message{}, false
        }

        message.ID = strconv.Itoa(dbID)
        message.FromID = strconv.Itoa(fromID)
        message.ToID = strconv.Itoa(toID)
        if listingID.Valid {
                message.ListingID = strconv.FormatInt(listingID.Int64, 10)
        } else {
                message.ListingID = ""
        }

        return message, true
}

// GetMessagesByUser retrieves all messages for a user from the database
func GetMessagesByUser(userID string) []models.Message {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return []models.Message{}
        }

        rows, err := GetDB().Query(`
                SELECT id, from_id, to_id, listing_id, content, read, created_at
                FROM messages
                WHERE from_id = $1 OR to_id = $1
                ORDER BY created_at
        `, userIDInt)
        if err != nil {
                log.Printf("Error getting messages by user: %v", err)
                return []models.Message{}
        }
        defer rows.Close()

        messages := []models.Message{}
        for rows.Next() {
                var message models.Message
                var id, fromID, toID int
                var listingID sql.NullInt64
                err := rows.Scan(&id, &fromID, &toID, &listingID, &message.Content, &message.Read, &message.CreatedAt)
                if err != nil {
                        log.Printf("Error scanning message row: %v", err)
                        continue
                }
                message.ID = strconv.Itoa(id)
                message.FromID = strconv.Itoa(fromID)
                message.ToID = strconv.Itoa(toID)
                if listingID.Valid {
                        message.ListingID = strconv.FormatInt(listingID.Int64, 10)
                } else {
                        message.ListingID = ""
                }

                messages = append(messages, message)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating message rows: %v", err)
        }

        return messages
}

// GetMessagesBetweenUsers retrieves all messages between two users from the database
func GetMessagesBetweenUsers(user1ID, user2ID string) []models.Message {
        user1IDInt, err := strconv.Atoi(user1ID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return []models.Message{}
        }

        user2IDInt, err := strconv.Atoi(user2ID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return []models.Message{}
        }

        rows, err := GetDB().Query(`
                SELECT id, from_id, to_id, listing_id, content, read, created_at
                FROM messages
                WHERE (from_id = $1 AND to_id = $2) OR (from_id = $2 AND to_id = $1)
                ORDER BY created_at
        `, user1IDInt, user2IDInt)
        if err != nil {
                log.Printf("Error getting messages between users: %v", err)
                return []models.Message{}
        }
        defer rows.Close()

        messages := []models.Message{}
        for rows.Next() {
                var message models.Message
                var id, fromID, toID int
                var listingID sql.NullInt64
                err := rows.Scan(&id, &fromID, &toID, &listingID, &message.Content, &message.Read, &message.CreatedAt)
                if err != nil {
                        log.Printf("Error scanning message row: %v", err)
                        continue
                }
                message.ID = strconv.Itoa(id)
                message.FromID = strconv.Itoa(fromID)
                message.ToID = strconv.Itoa(toID)
                if listingID.Valid {
                        message.ListingID = strconv.FormatInt(listingID.Int64, 10)
                } else {
                        message.ListingID = ""
                }

                messages = append(messages, message)
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating message rows: %v", err)
        }

        return messages
}

// SaveMessage saves a message to the database
func SaveMessage(msg models.Message) string {
        fromID, err := strconv.Atoi(msg.FromID)
        if err != nil {
                log.Printf("Invalid from user ID: %v", err)
                return ""
        }

        toID, err := strconv.Atoi(msg.ToID)
        if err != nil {
                log.Printf("Invalid to user ID: %v", err)
                return ""
        }

        var listingIDParam interface{} = nil
        if msg.ListingID != "" {
                listingIDInt, err := strconv.Atoi(msg.ListingID)
                if err != nil {
                        log.Printf("Invalid listing ID: %v", err)
                        return ""
                }
                listingIDParam = listingIDInt
        }

        // If the message has no ID, insert a new message
        if msg.ID == "" {
                var id int
                err := GetDB().QueryRow(`
                        INSERT INTO messages (from_id, to_id, listing_id, content, read, created_at)
                        VALUES ($1, $2, $3, $4, $5, $6)
                        RETURNING id
                `, fromID, toID, listingIDParam, msg.Content, msg.Read, msg.CreatedAt).Scan(&id)

                if err != nil {
                        log.Printf("Error creating message: %v", err)
                        return ""
                }

                return strconv.Itoa(id)
        }

        // Message has an ID, update existing message
        messageID, err := strconv.Atoi(msg.ID)
        if err != nil {
                log.Printf("Invalid message ID for update: %v", err)
                return ""
        }

        _, err = GetDB().Exec(`
                UPDATE messages
                SET from_id = $1, to_id = $2, listing_id = $3, content = $4, read = $5
                WHERE id = $6
        `, fromID, toID, listingIDParam, msg.Content, msg.Read, messageID)

        if err != nil {
                log.Printf("Error updating message: %v", err)
                return ""
        }

        return msg.ID
}

// MarkMessageAsRead marks a message as read in the database
func MarkMessageAsRead(id string) bool {
        messageID, err := strconv.Atoi(id)
        if err != nil {
                log.Printf("Invalid message ID: %v", err)
                return false
        }

        result, err := GetDB().Exec(`
                UPDATE messages
                SET read = true
                WHERE id = $1
        `, messageID)

        if err != nil {
                log.Printf("Error marking message as read: %v", err)
                return false
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                log.Printf("Error getting rows affected: %v", err)
                return false
        }

        return rowsAffected > 0
}

// GetFavorites retrieves all favorite listing IDs for a user from the database
func GetFavorites(userID string) []string {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return []string{}
        }

        rows, err := GetDB().Query(`
                SELECT listing_id
                FROM favorites
                WHERE user_id = $1
        `, userIDInt)
        if err != nil {
                log.Printf("Error getting favorites: %v", err)
                return []string{}
        }
        defer rows.Close()

        favoriteIDs := []string{}
        for rows.Next() {
                var listingID int
                err := rows.Scan(&listingID)
                if err != nil {
                        log.Printf("Error scanning favorite row: %v", err)
                        continue
                }
                favoriteIDs = append(favoriteIDs, strconv.Itoa(listingID))
        }

        if err = rows.Err(); err != nil {
                log.Printf("Error iterating favorite rows: %v", err)
        }

        return favoriteIDs
}

// AddFavorite adds a listing to a user's favorites in the database
func AddFavorite(userID, listingID string) bool {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return false
        }

        listingIDInt, err := strconv.Atoi(listingID)
        if err != nil {
                log.Printf("Invalid listing ID: %v", err)
                return false
        }

        // Check if already favorited
        var count int
        err = GetDB().QueryRow(`
                SELECT COUNT(*)
                FROM favorites
                WHERE user_id = $1 AND listing_id = $2
        `, userIDInt, listingIDInt).Scan(&count)

        if err != nil {
                log.Printf("Error checking for existing favorite: %v", err)
                return false
        }

        if count > 0 {
                // Already favorited
                return false
        }

        // Add to favorites
        _, err = GetDB().Exec(`
                INSERT INTO favorites (user_id, listing_id)
                VALUES ($1, $2)
        `, userIDInt, listingIDInt)

        if err != nil {
                log.Printf("Error adding favorite: %v", err)
                return false
        }

        return true
}

// RemoveFavorite removes a listing from a user's favorites in the database
func RemoveFavorite(userID, listingID string) bool {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return false
        }

        listingIDInt, err := strconv.Atoi(listingID)
        if err != nil {
                log.Printf("Invalid listing ID: %v", err)
                return false
        }

        result, err := GetDB().Exec(`
                DELETE FROM favorites
                WHERE user_id = $1 AND listing_id = $2
        `, userIDInt, listingIDInt)

        if err != nil {
                log.Printf("Error removing favorite: %v", err)
                return false
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                log.Printf("Error getting rows affected: %v", err)
                return false
        }

        return rowsAffected > 0
}

// IsFavorite checks if a listing is in a user's favorites in the database
func IsFavorite(userID, listingID string) bool {
        userIDInt, err := strconv.Atoi(userID)
        if err != nil {
                log.Printf("Invalid user ID: %v", err)
                return false
        }

        listingIDInt, err := strconv.Atoi(listingID)
        if err != nil {
                log.Printf("Invalid listing ID: %v", err)
                return false
        }

        var count int
        err = GetDB().QueryRow(`
                SELECT COUNT(*)
                FROM favorites
                WHERE user_id = $1 AND listing_id = $2
        `, userIDInt, listingIDInt).Scan(&count)

        if err != nil {
                log.Printf("Error checking if favorite: %v", err)
                return false
        }

        return count > 0
}