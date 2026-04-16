package models

import (
	"time"
	"github.com/google/uuid"
)

type Item struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type Invoice struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	InvoiceNumber   string          `gorm:"uniqueIndex;not null" json:"invoice_number"`
	SenderName      string          `gorm:"not null" json:"sender_name"`
	SenderAddress   string          `gorm:"type:text;not null" json:"sender_address"`
	ReceiverName    string          `gorm:"not null" json:"receiver_name"`
	ReceiverAddress string          `gorm:"type:text;not null" json:"receiver_address"`
	TotalAmount     float64         `gorm:"not null" json:"total_amount"`
	CreatedBy       uint            `gorm:"not null" json:"created_by"` // ID dari User (Admin/Kerani)
	CreatedAt       time.Time       `json:"created_at"`
	Details         []InvoiceDetail `gorm:"foreignKey:InvoiceID" json:"details"`
}

type InvoiceDetail struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:uuid;not null" json:"invoice_id"`
	ItemID    uint      `gorm:"not null" json:"item_id"`
	ItemName  string    `gorm:"not null" json:"item_name"` // Nama saat transaksi
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"not null" json:"price"`    // Harga snapshot
	Subtotal  float64   `gorm:"not null" json:"subtotal"`
}