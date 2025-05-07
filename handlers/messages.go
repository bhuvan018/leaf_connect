package handlers

import (
        "encoding/json"
        "net/http"
        "time"

        "github.com/gorilla/mux"

        "github.com/plantexchange/app/models"
        "github.com/plantexchange/app/utils"
)

// GetMessages gets all messages for the current user
func GetMessages(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get all messages for this user
        allMessages := utils.GetMessagesByUser(userID)
        
        // Enhance messages with user and listing information
        messagesWithInfo := []models.MessageWithUser{}
        
        for _, msg := range allMessages {
                // Get from user
                fromUser, fromExists := utils.GetUser(msg.FromID)
                
                // Get to user
                toUser, toExists := utils.GetUser(msg.ToID)
                
                // Get listing
                listing, listingExists := utils.GetListing(msg.ListingID)
                
                if fromExists && toExists && listingExists {
                        msgWithInfo := models.MessageWithUser{
                                Message:  msg,
                                FromUser: fromUser.ToUserResponse(),
                                ToUser:   toUser.ToUserResponse(),
                                Listing:  listing,
                        }
                        
                        messagesWithInfo = append(messagesWithInfo, msgWithInfo)
                }
        }
        
        // Return messages
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(messagesWithInfo)
}

// GetMessage gets a specific message by ID
func GetMessage(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get message ID from URL path
        vars := mux.Vars(r)
        messageID := vars["id"]

        // Find message
        msg, exists := utils.GetMessage(messageID)
        if !exists {
                http.Error(w, "Message not found", http.StatusNotFound)
                return
        }

        // Check if user is part of the conversation
        if msg.FromID != userID && msg.ToID != userID {
                http.Error(w, "Unauthorized", http.StatusForbidden)
                return
        }

        // Mark message as read if recipient is viewing it
        if msg.ToID == userID && !msg.Read {
                utils.MarkMessageAsRead(messageID)
                msg.Read = true
        }

        // Get from user
        fromUser, fromExists := utils.GetUser(msg.FromID)
        
        // Get to user
        toUser, toExists := utils.GetUser(msg.ToID)
        
        // Get listing
        listing, listingExists := utils.GetListing(msg.ListingID)
        
        if !fromExists || !toExists || !listingExists {
                http.Error(w, "Message references missing data", http.StatusInternalServerError)
                return
        }

        // Create message with additional info
        msgWithInfo := models.MessageWithUser{
                Message:  msg,
                FromUser: fromUser.ToUserResponse(),
                ToUser:   toUser.ToUserResponse(),
                Listing:  listing,
        }

        // Return message
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(msgWithInfo)
}

// SendMessage sends a new message
func SendMessage(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        fromID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Parse request
        var msg models.Message
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
        }

        // Validate required fields
        if msg.ToID == "" || msg.ListingID == "" || msg.Content == "" {
                http.Error(w, "Recipient, listing, and message content are required", http.StatusBadRequest)
                return
        }

        // Check if recipient exists
        if _, exists := utils.GetUser(msg.ToID); !exists {
                http.Error(w, "Recipient not found", http.StatusNotFound)
                return
        }

        // Check if listing exists
        if _, exists := utils.GetListing(msg.ListingID); !exists {
                http.Error(w, "Listing not found", http.StatusNotFound)
                return
        }

        // Set sender and timestamp
        msg.FromID = fromID
        msg.CreatedAt = time.Now()
        msg.Read = false

        // Save message
        messageID := utils.SaveMessage(msg)
        msg.ID = messageID

        // Return created message
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(msg)
}

// GetConversations gets all conversations for the current user
func GetConversations(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get all messages for this user
        allMessages := utils.GetMessagesByUser(userID)
        
        // Group messages by conversation partner
        conversationPartners := make(map[string][]models.Message)
        
        for _, msg := range allMessages {
                var partnerID string
                if msg.FromID == userID {
                        partnerID = msg.ToID
                } else {
                        partnerID = msg.FromID
                }
                
                conversationPartners[partnerID] = append(conversationPartners[partnerID], msg)
        }
        
        // Create conversation summaries
        conversations := []models.Conversation{}
        
        for partnerID, messages := range conversationPartners {
                // Get partner user info
                partner, exists := utils.GetUser(partnerID)
                if !exists {
                        continue
                }
                
                // Find the most recent message
                var lastMessage models.Message
                var lastActivity time.Time
                unreadCount := 0
                
                for _, msg := range messages {
                        if msg.CreatedAt.After(lastActivity) {
                                lastMessage = msg
                                lastActivity = msg.CreatedAt
                        }
                        
                        // Count unread messages sent to the current user
                        if msg.ToID == userID && !msg.Read {
                                unreadCount++
                        }
                }
                
                // Create conversation
                conversation := models.Conversation{
                        UserID:       partnerID,
                        Username:     partner.Username,
                        ProfilePic:   partner.ProfilePic,
                        LastMessage:  lastMessage.Content,
                        LastActivity: lastActivity,
                        Unread:       unreadCount,
                }
                
                conversations = append(conversations, conversation)
        }
        
        // Return conversations
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(conversations)
}

// GetConversation gets all messages between current user and another user
func GetConversation(w http.ResponseWriter, r *http.Request) {
        // Get current session
        session, _ := utils.SessionStore.Get(r, "session")
        
        // Check if userID exists in session
        userID, ok := session.Values["userID"].(string)
        if !ok {
                http.Error(w, "Not authenticated", http.StatusUnauthorized)
                return
        }

        // Get partner ID from URL path
        vars := mux.Vars(r)
        partnerID := vars["userId"]

        // Check if partner exists
        partner, exists := utils.GetUser(partnerID)
        if !exists {
                http.Error(w, "User not found", http.StatusNotFound)
                return
        }

        // Get messages between these users
        messages := utils.GetMessagesBetweenUsers(userID, partnerID)
        
        // Enhance messages with user and listing information
        messagesWithInfo := []models.MessageWithUser{}
        
        for _, msg := range messages {
                // Get listing
                listing, listingExists := utils.GetListing(msg.ListingID)
                if !listingExists {
                        continue
                }
                
                // Mark as read if this user is the recipient
                if msg.ToID == userID && !msg.Read {
                        utils.MarkMessageAsRead(msg.ID)
                        msg.Read = true
                }
                
                // Create message with additional info
                fromUser, _ := utils.GetUser(msg.FromID)
                toUser, _ := utils.GetUser(msg.ToID)
                msgWithInfo := models.MessageWithUser{
                        Message:  msg,
                        FromUser: fromUser.ToUserResponse(),
                        ToUser:   toUser.ToUserResponse(),
                        Listing:  listing,
                }
                
                messagesWithInfo = append(messagesWithInfo, msgWithInfo)
        }
        
        // Create conversation
        conversation := models.Conversation{
                UserID:       partnerID,
                Username:     partner.Username,
                ProfilePic:   partner.ProfilePic,
                Messages:     messagesWithInfo,
                Unread:       0, // All messages marked as read
        }
        
        // Return conversation
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(conversation)
}
