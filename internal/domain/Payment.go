package domain

import "time"

type Payment struct {
	ID            uint      `gorm:"PrimaryKey" json:"id"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	TransactionId uint      `json:"transaction_id"`
	CustomerId    string    `json:"customer_id"`
	PaymentId     string    `json:"payment_id"`
	Status        string    `json:"status"`
	Response      string    `json:"response"`
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `gorm:"default:current_timestamp"`
}
