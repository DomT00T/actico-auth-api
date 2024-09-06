package models

import "time"

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordHash      string    `json:"-"`
	UserType          string    `json:"user_type" binding:"required,oneof=musician venue_owner admin"`
	FirstName         string    `json:"first_name" binding:"required"`
	LastName          string    `json:"last_name" binding:"required"`
	PhoneNumber       string    `json:"phone_number" binding:"required"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	Gender            string    `json:"gender" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastLogin         time.Time `json:"last_login"`
	IsVerified        string    `json:"is_verified"`
	IsPremium         string    `json:"is_premium"`
	PremiumExpiryDate time.Time `json:"premium_expiry_date"`
}