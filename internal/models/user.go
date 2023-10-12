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

type UserResponse struct {
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `gorm:"size:255;" json:"role"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}
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
type MerchantInfoResponse struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	MerchantName string    `gorm:"size:255;not null" json:"merchant_name"`
	Location     string    `gorm:"size:255;not null" json:"location"`
	UserId       uint      `gorm:"not null" json:"user_id"`
	UserName     string    `gorm:"size:255;not null" json:"user_name"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (mi *MerchantInfo) ToMerchantInfoResponse() MerchantInfoResponse {
	return MerchantInfoResponse{
		ID:           mi.ID,
		MerchantName: mi.MerchantName,
		Location:     mi.Location,
		UserId:       mi.UserId,
		UserName:     mi.User.FirstName + " " + mi.User.LastName,
		CreatedAt:    mi.CreatedAt,
		UpdatedAt:    mi.UpdatedAt,
	}
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

type MerchantOfferResponse struct {
	ID             uint `gorm:"primaryKey"`
	CardName       Card `gorm:"not null" json:"card_name"`
	DiscountRate   uint `gorm:"not null" json:"discount_rate"`
	MerchantInfoID uint `gorm:"not null" json:"merchant_info_id"`
}

func (mo *MerchantOffer) ToMerchantOfferResponse() MerchantOfferResponse {
	return MerchantOfferResponse{
		ID:             mo.ID,
		CardName:       mo.CardName,
		DiscountRate:   mo.DiscountRate,
		MerchantInfoID: mo.MerchantInfoID,
	}
}

// user will by subscription
type CardSubscription struct {
	gorm.Model
	ID       uint   `gorm:"primary_key" json:"id"`
	CardName string `gorm:"size:255;not null" json:"card_name"`
	Number   int64  `gorm:"size:255;not null" json:"number"`
	UserId   uint   `gorm:"not null" json:"user_id"`
	User     User   `gorm:"foreignKey:UserId"`
}

type CardSubscriptionResponse struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	CardName string `gorm:"size:255;not null" json:"card_name"`
	Number   int64  `gorm:"size:255;not null" json:"number"`
	UserId   uint   `gorm:"not null" json:"user_id"`
	UserName string `gorm:"size:255;not null" json:"user_name"`
}

func (cs *CardSubscription) ToCardSubscriptionResponse() CardSubscriptionResponse {
	return CardSubscriptionResponse{
		ID:       cs.ID,
		CardName: cs.CardName,
		Number:   cs.Number,
		UserId:   cs.UserId,
		UserName: cs.User.FirstName + " " + cs.User.LastName,
	}
}

// merchant and user will generate transactions
type Transaction struct {
	gorm.Model
	ID uint `gorm:"primary_key" json:"id"`
	// UserId             uint             `gorm:"not null" json:"user_id"`
	// User               User             `gorm:"foreignKey:UserId"`
	CardSubscription   CardSubscription `gorm:"foreignKey:CardSubscriptionID"`
	CardSubscriptionID uint             `gorm:"not null" json:"card_subscription_id"`
	MerchantOffer      MerchantOffer    `gorm:"foreignKey:MerchantOfferID"`
	MerchantOfferID    uint             `gorm:"not null" json:"merchant_offer_id"`
	Amount             float64          `gorm:"not null" json:"amount"`
	BillNumber         string           `gorm:"size:255;not null" json:"bill_number"`
	CreatedAt          time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type TransactionResponse struct {
	ID                 int       `json:"id"`
	CardName           string    `json:"card_name"`
	MerchantOfferID    uint      `gorm:"not null" json:"merchant_offer_id"`
	Amount             float64   `gorm:"not null" json:"amount"`
	BillNumber         string    `gorm:"size:255;not null" json:"bill_number"`
	CardSubscriptionID uint      `gorm:"not null" json:"card_subscription_id"`
	CreatedAt          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CardNumber         int64     `gorm:"size:255;not null" json:"card_number"`
	UserID             uint      `gorm:"not null" json:"user_id"`
	UserName           string    `gorm:"size:255;not null" json:"user_name"`
	MerchantName       string    `gorm:"size:255;not null" json:"merchant_name"`
	// Add more fields as needed
}

// ToTransactionResponse converts the Transaction struct to a TransactionResponse struct
func (t Transaction) ToTransactionResponse() TransactionResponse {
	return TransactionResponse{
		ID:                 int(t.ID),
		CardName:           t.CardSubscription.CardName,
		MerchantOfferID:    t.MerchantOfferID,
		Amount:             t.Amount,
		BillNumber:         t.BillNumber,
		CardSubscriptionID: t.CardSubscriptionID,
		CreatedAt:          t.CreatedAt,
		CardNumber:         t.CardSubscription.Number,
		UserID:             t.CardSubscription.UserId,
		UserName:           t.CardSubscription.User.FirstName + " " + t.CardSubscription.User.LastName,
		MerchantName:       t.MerchantOffer.MerchantInfo.MerchantName,
	}
}

type ValidateRequest struct {
	MerchantId uint `json:"merchant_id"`
	CardId     uint `json:"card_id"`
}
