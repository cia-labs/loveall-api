package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser     Role = "user"
	RoleMerchant Role = "merchant"
	RoleAdmin    Role = "admin"
)

type User struct {
	gorm.Model
	ID        uint      `gorm:"primary_key" json:"id"`
	FirstName string    `gorm:"size:255;not null" json:"first_name"`
	LastName  string    `gorm:"size:255;not null" json:"last_name"`
	Email     string    `gorm:"size:255;unique;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	Role      string    `gorm:"size:255;" json:"role"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type MerchantInfo struct {
	gorm.Model
	ID           uint      `gorm:"primary_key" json:"id"`
	MerchantName string    `gorm:"size:255;not null" json:"merchant_name"`
	Location     string    `gorm:"size:255;not null" json:"location"`
	UserId       uint      `gorm:"not null" json:"user_id"`
	User         User      `gorm:"foreignKey:UserId"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// merchant will create offers
type MerchantOffer struct {
	gorm.Model
	ID             uint         `gorm:"primaryKey"`
	CardName       Card         `gorm:"not null" json:"card_name"`
	DiscountRate   uint         `gorm:"not null" json:"discount_rate"`
	MerchantInfoID uint         `gorm:"not null" json:"merchant_info_id"`
	MerchantInfo   MerchantInfo `gorm:"foreignKey:MerchantInfoID"`
}

// user will by subscription
type CardSubscription struct {
	gorm.Model
	ID       uint   `gorm:"primary_key" json:"id"`
	CardName string `gorm:"size:255;not null" json:"card_name"`
	Number   string `gorm:"size:255;not null" json:"number"`
	UserId   uint   `gorm:"not null" json:"user_id"`
	User     User   `gorm:"foreignKey:UserId"`
}

// merchant and user will generate transactions
type Transaction struct {
	gorm.Model
	ID                 uint             `gorm:"primary_key" json:"id"`
	UserId             uint             `gorm:"not null" json:"user_id"`
	User               User             `gorm:"foreignKey:UserId"`
	CardSubscription   CardSubscription `gorm:"foreignKey:CardSubscriptionID"`
	CardSubscriptionID uint             `gorm:"not null" json:"card_subscription_id"`
	MerchantOffer      MerchantOffer    `gorm:"foreignKey:MerchantOfferID"`
	MerchantOfferID    uint             `gorm:"not null" json:"merchant_offer_id"`
	Amount             float64          `gorm:"not null" json:"amount"`
	CreatedAt          time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
