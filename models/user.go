package models

import (
        "time"
)

type User struct {
        ID          string    `json:"id"`
        Email       string    `json:"email"`
        Username    string    `json:"username"`
        Password    string    `json:"password"` // Used for registration/login, not exposed in responses
        Name        string    `json:"name"`
        Location    string    `json:"location"`
        Bio         string    `json:"bio"`
        ProfilePic  string    `json:"profilePic"`
        CreatedAt   time.Time `json:"createdAt"`
        LastLoginAt time.Time `json:"lastLoginAt"`
        Favorites   []string  `json:"favorites"` // Array of listing IDs
}

// UserResponse is a struct to return user data without sensitive information
type UserResponse struct {
        ID         string    `json:"id"`
        Username   string    `json:"username"`
        Name       string    `json:"name"`
        Location   string    `json:"location"`
        Bio        string    `json:"bio"`
        ProfilePic string    `json:"profilePic"`
        CreatedAt  time.Time `json:"createdAt"`
}

// ToUserResponse converts a User to a UserResponse
func (u *User) ToUserResponse() UserResponse {
        return UserResponse{
                ID:         u.ID,
                Username:   u.Username,
                Name:       u.Name,
                Location:   u.Location,
                Bio:        u.Bio,
                ProfilePic: u.ProfilePic,
                CreatedAt:  u.CreatedAt,
        }
}
