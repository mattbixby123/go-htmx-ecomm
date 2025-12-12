package main

import (
	"time"

	"gorm.io/gorm"
)

// User represents a registered user
type User struct {
	ID           string `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Session represents a JWT session (optional - for token blacklisting)
type Session struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"not null;index"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

// Product represents a product in the store
type Product struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       int64 `gorm:"not null"` // in cents
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CartItem represents an item in a user's cart
type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    string  `gorm:"not null;index"`
	ProductID string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Order represents a completed order
type Order struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"not null;index"`
	Total     int64  `gorm:"not null"`
	Status    string `gorm:"not null"`
	PaymentID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem represents a product in an order
type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   string  `gorm:"not null;index"`
	ProductID string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Price     int64   `gorm:"not null"` // Price at time of purchase
	Product   Product `gorm:"foreignKey:ProductID"`
	CreatedAt time.Time
}

// TableName overrides for GORM
func (User) TableName() string      { return "users" }
func (Session) TableName() string   { return "sessions" }
func (Product) TableName() string   { return "products" }
func (CartItem) TableName() string  { return "cart_items" }
func (Order) TableName() string     { return "orders" }
func (OrderItem) TableName() string { return "order_items" }

// BeforeCreate hooks for UUID generation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return nil
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateUUID()
	}
	return nil
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = generateUUID()
	}
	return nil
}

func generateUUID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
