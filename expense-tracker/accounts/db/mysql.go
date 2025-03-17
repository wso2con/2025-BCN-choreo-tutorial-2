package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jo/choreo-tutorial/accounts/config"
	"github.com/jo/choreo-tutorial/accounts/models"
)

// MySQLDB implements the Database interface for MySQL
type MySQLDB struct {
	db *sql.DB
}

// NewMySQLDB creates a new MySQL database connection
func NewMySQLDB(cfg *config.Config) (*MySQLDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQLDB{db: db}, nil
}

// CreateTables creates the necessary tables if they don't exist
func (m *MySQLDB) CreateTables() error {
	// Create bills table
	_, err := m.db.Exec(`
	CREATE TABLE IF NOT EXISTS bills (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		total DECIMAL(10, 2) NOT NULL DEFAULT 0,
		due_date DATE,
		paid BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		return err
	}

	// Create bill_items table
	_, err = m.db.Exec(`
	CREATE TABLE IF NOT EXISTS bill_items (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		bill_id BIGINT NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		amount DECIMAL(10, 2) NOT NULL,
		quantity INT NOT NULL DEFAULT 1,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (bill_id) REFERENCES bills(id) ON DELETE CASCADE
	)
	`)
	return err
}

// Close closes the database connection
func (m *MySQLDB) Close() error {
	return m.db.Close()
}

// GetBills returns all bills with summary information
func (m *MySQLDB) GetBills() ([]models.BillSummary, error) {
	rows, err := m.db.Query(`
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
		var dueDate sql.NullTime

		err := rows.Scan(
			&bill.ID,
			&bill.Title,
			&bill.Description,
			&bill.Total,
			&dueDate,
			&bill.Paid,
			&bill.CreatedAt,
			&bill.UpdatedAt,
			&bill.ItemCount,
		)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			bill.DueDate = dueDate.Time
		}

		bills = append(bills, bill)
	}

	return bills, nil
}

// GetBill returns a single bill with all its items
func (m *MySQLDB) GetBill(id int64) (*models.Bill, error) {
	// Get the bill
	var bill models.Bill
	var dueDate sql.NullTime

	err := m.db.QueryRow(`
	SELECT id, title, description, total, due_date, paid, created_at, updated_at
	FROM bills
	WHERE id = ?
	`, id).Scan(
		&bill.ID,
		&bill.Title,
		&bill.Description,
		&bill.Total,
		&dueDate,
		&bill.Paid,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if dueDate.Valid {
		bill.DueDate = dueDate.Time
	}

	// Get the bill items
	items, err := m.GetBillItems(id)
	if err != nil {
		return nil, err
	}
	bill.Items = items

	return &bill, nil
}

// CreateBill creates a new bill and its items
func (m *MySQLDB) CreateBill(billInput *models.BillInput) (int64, error) {
	// Start a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Parse due date
	var dueDate *time.Time
	if billInput.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", billInput.DueDate)
		if err != nil {
			return 0, fmt.Errorf("invalid due date format: %v", err)
		}
		dueDate = &parsedDate
	}

	// Calculate total from items
	var total float64
	for _, item := range billInput.Items {
		total += item.Amount * float64(item.Quantity)
	}

	// Insert bill
	result, err := tx.Exec(`
	INSERT INTO bills (title, description, total, due_date, paid)
	VALUES (?, ?, ?, ?, ?)
	`, billInput.Title, billInput.Description, total, dueDate, billInput.Paid)
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
func (m *MySQLDB) UpdateBill(id int64, billInput *models.BillInput) error {
	// Start a transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Parse due date
	var dueDate *time.Time
	if billInput.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", billInput.DueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format: %v", err)
		}
		dueDate = &parsedDate
	}

	// Calculate total from items
	var total float64
	for _, item := range billInput.Items {
		total += item.Amount * float64(item.Quantity)
	}

	// Update bill
	_, err = tx.Exec(`
	UPDATE bills
	SET title = ?, description = ?, total = ?, due_date = ?, paid = ?
	WHERE id = ?
	`, billInput.Title, billInput.Description, total, dueDate, billInput.Paid, id)
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
func (m *MySQLDB) DeleteBill(id int64) error {
	_, err := m.db.Exec("DELETE FROM bills WHERE id = ?", id)
	return err
}

// GetBillItems returns all items for a bill
func (m *MySQLDB) GetBillItems(billID int64) ([]models.BillItem, error) {
	rows, err := m.db.Query(`
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
func (m *MySQLDB) GetBillItem(id int64) (*models.BillItem, error) {
	var item models.BillItem
	err := m.db.QueryRow(`
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
func (m *MySQLDB) CreateBillItem(billID int64, itemInput *models.BillItemInput) (int64, error) {
	// Start a transaction
	tx, err := m.db.Begin()
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
func (m *MySQLDB) UpdateBillItem(id int64, itemInput *models.BillItemInput) error {
	// Start a transaction
	tx, err := m.db.Begin()
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
func (m *MySQLDB) DeleteBillItem(id int64) error {
	// Start a transaction
	tx, err := m.db.Begin()
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
