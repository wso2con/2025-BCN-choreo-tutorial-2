package db

import (
	"errors"

	"github.com/jo/choreo-tutorial/accounts/config"
	"github.com/jo/choreo-tutorial/accounts/models"
)

// Common errors
var (
	ErrNotFound = errors.New("record not found")
)

// Database is the interface for database operations
type Database interface {
	// Bills
	GetBills() ([]models.BillSummary, error)
	GetBill(id int64) (*models.Bill, error)
	CreateBill(bill *models.BillInput) (int64, error)
	UpdateBill(id int64, bill *models.BillInput) error
	DeleteBill(id int64) error

	// BillItems
	GetBillItems(billID int64) ([]models.BillItem, error)
	GetBillItem(id int64) (*models.BillItem, error)
	CreateBillItem(billID int64, item *models.BillItemInput) (int64, error)
	UpdateBillItem(id int64, item *models.BillItemInput) error
	DeleteBillItem(id int64) error

	// Database management
	CreateTables() error
	Close() error
}

// InitDB initializes the database based on the configuration
func InitDB(cfg *config.Config) (Database, error) {
	switch cfg.DBType {
	case "mysql":
		return NewMySQLDB(cfg)
	case "sqlite":
		return NewSQLiteDB(cfg)
	default:
		return NewSQLiteDB(cfg) // Default to SQLite
	}
}
