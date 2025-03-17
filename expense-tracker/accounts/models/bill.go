package models

import (
	"encoding/json"
	"time"
)

// Bill represents a financial bill
type Bill struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Total       float64     `json:"total"`
	DueDate     time.Time   `json:"due_date"`
	Paid        bool        `json:"paid"`
	Items       []BillItem  `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// BillItem represents an item within a bill
type BillItem struct {
	ID          int64     `json:"id"`
	BillID      int64     `json:"bill_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BillInput represents the JSON input for creating/updating a bill
type BillInput struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	DueDate     string          `json:"due_date"` // ISO format (YYYY-MM-DD)
	Paid        bool            `json:"paid"`
	Items       []BillItemInput `json:"items"`
}

// BillItemInput represents the JSON input for creating/updating a bill item
type BillItemInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Quantity    int     `json:"quantity"`
}

// FromJSON converts a JSON string to a BillInput
func (b *BillInput) FromJSON(data []byte) error {
	return json.Unmarshal(data, b)
}

// ToJSON converts a Bill to a JSON string
func (b *Bill) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// BillSummary represents a summary of a bill with total amount
type BillSummary struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Total       float64   `json:"total"`
	DueDate     time.Time `json:"due_date"`
	Paid        bool      `json:"paid"`
	ItemCount   int       `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
