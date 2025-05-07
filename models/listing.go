package models

import (
	"time"
)

type Listing struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // plant, seed, cutting
	PlantType   string    `json:"plantType"` // indoor, outdoor, vegetable, herb, etc.
	Price       float64   `json:"price"`
	TradeFor    string    `json:"tradeFor"` // What the user is willing to trade for
	Location    string    `json:"location"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Status      string    `json:"status"` // available, pending, sold, traded
}

// ListingWithUser combines listing data with basic user information
type ListingWithUser struct {
	Listing
	User UserResponse `json:"user"`
}
