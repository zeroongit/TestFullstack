package models

import (
	"time"
	"github.com/google/uuid"
)

type Item struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}

type Invoice struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	InvoiceNumber string `gorm:"unique;not null" json:"invoice_number"`
	SenderName string `json:"sender_name"`
	SenderAddress string `json:"sender_address"`
	ReceiverName string `json:"receiver_name"`
	ReceiverAddress string `json:"receiver_address"`
	TotalAmount float64 `json:"total_amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Details []InvoiceDetail `gorm:"foreignKey:InvoiceID" json:"details"`
}

type InvoiceDetail struct {
	ID uint `gorm:"primaryKey" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:uuid" json:"invoice_id"`
	Name string `json:"name"`
	ItemID uint `json:"item_id"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}