package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jo/choreo-tutorial/accounts/config"
	"github.com/jo/choreo-tutorial/accounts/models"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB implements the Database interface for SQLite
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(cfg *config.Config) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

// CreateTables creates the necessary tables if they don't exist
func (s *SQLiteDB) CreateTables() error {
	// Create bills table
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS bills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		total REAL NOT NULL DEFAULT 0,
		due_date DATE,
		paid INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	// Create bill_items table
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS bill_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bill_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		amount REAL NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 1,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE
	)
	`)
	if err != nil {
		return err
	}

	// Create update trigger for bills
	_, err = s.db.Exec(`
	CREATE TRIGGER IF NOT EXISTS bills_update_trigger
	AFTER UPDATE ON bills
	FOR EACH ROW
	BEGIN
		UPDATE bills SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;
	`)
	if err != nil {
		return err
	}

	// Create update trigger for bill_items
	_, err = s.db.Exec(`
	CREATE TRIGGER IF NOT EXISTS bill_items_update_trigger
	AFTER UPDATE ON bill_items
	FOR EACH ROW
	BEGIN
		UPDATE bill_items SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;
	`)
	return err
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// GetBills returns all bills with summary information
func (s *SQLiteDB) GetBills() ([]models.BillSummary, error) {
	rows, err := s.db.Query(`
	SELECT b.id, b.title, b.description, b.total, b.due_date, b.paid, b.created_at, b.updated_at, COUNT(i.id) as item_count
	FROM bills b
	LEFT JOIN bill_items i ON b.id = i.bill_id
	GROUP BY b.id
	ORDER BY b.due_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.BillSummary
	for rows.Next() {
		var bill models.BillSummary
		var dueDate sql.NullString
		var paid int

		err := rows.Scan(
			&bill.ID,
			&bill.Title,
			&bill.Description,
			&bill.Total,
			&dueDate,
			&paid,
			&bill.CreatedAt,
			&bill.UpdatedAt,
			&bill.ItemCount,
		)
		if err != nil {
			return nil, err
		}

		bill.Paid = paid == 1

		if dueDate.Valid && dueDate.String != "" {
			parsedDate, err := time.Parse("2006-01-02", dueDate.String)
			if err == nil {
				bill.DueDate = parsedDate
			}
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

// GetBill returns a single bill with all its items
func (s *SQLiteDB) GetBill(id int64) (*models.Bill, error) {
	// Get the bill
	var bill models.Bill
	var dueDate sql.NullString
	var paid int

	err := s.db.QueryRow(`
	SELECT id, title, description, total, due_date, paid, created_at, updated_at
	FROM bills
	WHERE id = ?
	`, id).Scan(
		&bill.ID,
		&bill.Title,
		&bill.Description,
		&bill.Total,
		&dueDate,
		&paid,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	bill.Paid = paid == 1

	if dueDate.Valid && dueDate.String != "" {
		parsedDate, err := time.Parse("2006-01-02", dueDate.String)
		if err == nil {
			bill.DueDate = parsedDate
		}
	}

	// Get the bill items
	items, err := s.GetBillItems(id)
	if err != nil {
		return nil, err
	}
	bill.Items = items

	return &bill, nil
}

// CreateBill creates a new bill and its items
func (s *SQLiteDB) CreateBill(billInput *models.BillInput) (int64, error) {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Parse due date
	var dueDate *string
	if billInput.DueDate != "" {
		// Validate date format
		_, err := time.Parse("2006-01-02", billInput.DueDate)
		if err != nil {
			return 0, fmt.Errorf("invalid due date format: %v", err)
		}
		dueDate = &billInput.DueDate
	}

	// Calculate total from items
	var total float64
	for _, item := range billInput.Items {
		total += item.Amount * float64(item.Quantity)
	}

	// SQLite uses integers for boolean (0=false, 1=true)
	paidInt := 0
	if billInput.Paid {
		paidInt = 1
	}

	// Insert bill
	result, err := tx.Exec(`
	INSERT INTO bills (title, description, total, due_date, paid)
	VALUES (?, ?, ?, ?, ?)
	`, billInput.Title, billInput.Description, total, dueDate, paidInt)
	if err != nil {
		return 0, err
	}

	// Get the bill ID
	billID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert bill items
	for _, item := range billInput.Items {
		_, err = tx.Exec(`
		INSERT INTO bill_items (bill_id, name, description, amount, quantity)
		VALUES (?, ?, ?, ?, ?)
		`, billID, item.Name, item.Description, item.Amount, item.Quantity)
		if err != nil {
			return 0, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	return billID, err
}

// UpdateBill updates an existing bill and its items
func (s *SQLiteDB) UpdateBill(id int64, billInput *models.BillInput) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Parse due date
	var dueDate *string
	if billInput.DueDate != "" {
		// Validate date format
		_, err := time.Parse("2006-01-02", billInput.DueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format: %v", err)
		}
		dueDate = &billInput.DueDate
	}

	// Calculate total from items
	var total float64
	for _, item := range billInput.Items {
		total += item.Amount * float64(item.Quantity)
	}

	// SQLite uses integers for boolean (0=false, 1=true)
	paidInt := 0
	if billInput.Paid {
		paidInt = 1
	}

	// Update bill
	_, err = tx.Exec(`
	UPDATE bills
	SET title = ?, description = ?, total = ?, due_date = ?, paid = ?
	WHERE id = ?
	`, billInput.Title, billInput.Description, total, dueDate, paidInt, id)
	if err != nil {
		return err
	}

	// Delete existing items
	_, err = tx.Exec("DELETE FROM bill_items WHERE bill_id = ?", id)
	if err != nil {
		return err
	}

	// Insert new items
	for _, item := range billInput.Items {
		_, err = tx.Exec(`
		INSERT INTO bill_items (bill_id, name, description, amount, quantity)
		VALUES (?, ?, ?, ?, ?)
		`, id, item.Name, item.Description, item.Amount, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteBill deletes a bill and its items
func (s *SQLiteDB) DeleteBill(id int64) error {
	_, err := s.db.Exec("DELETE FROM bills WHERE id = ?", id)
	return err
}

// GetBillItems returns all items for a bill
func (s *SQLiteDB) GetBillItems(billID int64) ([]models.BillItem, error) {
	rows, err := s.db.Query(`
	SELECT id, bill_id, name, description, amount, quantity, created_at, updated_at
	FROM bill_items
	WHERE bill_id = ?
	`, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.BillItem
	for rows.Next() {
		var item models.BillItem
		err := rows.Scan(
			&item.ID,
			&item.BillID,
			&item.Name,
			&item.Description,
			&item.Amount,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetBillItem returns a single bill item
func (s *SQLiteDB) GetBillItem(id int64) (*models.BillItem, error) {
	var item models.BillItem
	err := s.db.QueryRow(`
	SELECT id, bill_id, name, description, amount, quantity, created_at, updated_at
	FROM bill_items
	WHERE id = ?
	`, id).Scan(
		&item.ID,
		&item.BillID,
		&item.Name,
		&item.Description,
		&item.Amount,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// CreateBillItem creates a new bill item
func (s *SQLiteDB) CreateBillItem(billID int64, itemInput *models.BillItemInput) (int64, error) {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert item
	result, err := tx.Exec(`
	INSERT INTO bill_items (bill_id, name, description, amount, quantity)
	VALUES (?, ?, ?, ?, ?)
	`, billID, itemInput.Name, itemInput.Description, itemInput.Amount, itemInput.Quantity)
	if err != nil {
		return 0, err
	}

	// Get the item ID
	itemID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Update bill total
	_, err = tx.Exec(`
	UPDATE bills
	SET total = (SELECT SUM(amount * quantity) FROM bill_items WHERE bill_id = ?)
	WHERE id = ?
	`, billID, billID)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	return itemID, err
}

// UpdateBillItem updates an existing bill item
func (s *SQLiteDB) UpdateBillItem(id int64, itemInput *models.BillItemInput) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the bill ID
	var billID int64
	err = tx.QueryRow("SELECT bill_id FROM bill_items WHERE id = ?", id).Scan(&billID)
	if err != nil {
		return err
	}

	// Update item
	_, err = tx.Exec(`
	UPDATE bill_items
	SET name = ?, description = ?, amount = ?, quantity = ?
	WHERE id = ?
	`, itemInput.Name, itemInput.Description, itemInput.Amount, itemInput.Quantity, id)
	if err != nil {
		return err
	}

	// Update bill total
	_, err = tx.Exec(`
	UPDATE bills
	SET total = (SELECT SUM(amount * quantity) FROM bill_items WHERE bill_id = ?)
	WHERE id = ?
	`, billID, billID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteBillItem deletes a bill item
func (s *SQLiteDB) DeleteBillItem(id int64) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the bill ID
	var billID int64
	err = tx.QueryRow("SELECT bill_id FROM bill_items WHERE id = ?", id).Scan(&billID)
	if err != nil {
		return err
	}

	// Delete item
	_, err = tx.Exec("DELETE FROM bill_items WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Update bill total
	_, err = tx.Exec(`
	UPDATE bills
	SET total = (SELECT SUM(amount * quantity) FROM bill_items WHERE bill_id = ?)
	WHERE id = ?
	`, billID, billID)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
