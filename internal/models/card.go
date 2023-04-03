package models

type Card string

const (
	Basic    Card = "Basic"
	Premium  Card = "Premium"
	Platinum Card = "Platinum"
)

// type Card struct {
// 	ID          uint      `gorm:"primary_key" json:"id"`
// 	Number      string    `gorm:"size:255;not null" json:"number"`
// 	CVV         string    `gorm:"size:255;not null" json:"cvv"`
// 	ExpiryMonth string    `gorm:"size:2;not null" json:"expiry_month"`
// 	ExpiryYear  string    `gorm:"size:4;not null" json:"expiry_year"`
// 	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
// 	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
// 	UserID      uint      `json:"-"`
// 	// ... other fields as necessary
// }
