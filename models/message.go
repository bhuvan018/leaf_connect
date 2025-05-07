package models

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	FromID    string    `json:"fromId"`
	ToID      string    `json:"toId"`
	ListingID string    `json:"listingId"`
	Content   string    `json:"content"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt"`
}

// MessageWithUser includes user information with the message
type MessageWithUser struct {
	Message
	FromUser UserResponse `json:"fromUser"`
	ToUser   UserResponse `json:"toUser"`
	Listing  Listing      `json:"listing"`
}

// Conversation represents a conversation between two users
type Conversation struct {
	UserID       string           `json:"userId"`
	Username     string           `json:"username"`
	ProfilePic   string           `json:"profilePic"`
	LastMessage  string           `json:"lastMessage"`
	LastActivity time.Time        `json:"lastActivity"`
	Messages     []MessageWithUser `json:"messages"`
	Unread       int              `json:"unread"`
}
