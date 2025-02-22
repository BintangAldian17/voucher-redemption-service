package transaction

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/BintangAldian17/voucher-redemption-service/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *Store) RedeemVouchers(tx *sql.Tx, customer *types.Customer, transactionItems []types.TransactionItem, totalPointsNeeded int) (int, error) {

	transaction := types.Transaction{
		CustomerID:      customer.ID,
		TotalPoints:     totalPointsNeeded,
		TransactionDate: time.Now(),
	}

	transactionID, err := s.CreateTransaction(tx, transaction)
	if err != nil {
		return 0, err
	}

	for _, item := range transactionItems {
		item.TransactionID = transactionID
		if err := s.CreateTransactionItem(tx, item); err != nil {
			return 0, err
		}
	}

	// Update total points customer
	newTotalPoint := customer.TotalPoints - totalPointsNeeded
	err = s.UpdateCustomerTotalPoints(tx, customer.ID, newTotalPoint)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

func (s *Store) CreateTransaction(tx *sql.Tx, transaction types.Transaction) (int, error) {
	query := `
        INSERT INTO transactions (customer_id, transaction_date, total_points)
        VALUES (?, ?, ?)
    `

	var stmt *sql.Stmt
	var err error

	if tx != nil {
		stmt, err = tx.Prepare(query)
	} else {
		stmt, err = s.db.Prepare(query)
	}

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(transaction.CustomerID, transaction.TransactionDate, transaction.TotalPoints)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}

func (s *Store) CreateTransactionItem(tx *sql.Tx, item types.TransactionItem) error {
	query := `
        INSERT INTO transaction_items (transaction_id, voucher_id, quantity, cost_in_points)
        VALUES (?, ?, ?, ?)
    `

	var stmt *sql.Stmt
	var err error

	if tx != nil {
		stmt, err = tx.Prepare(query)
	} else {
		stmt, err = s.db.Prepare(query)
	}

	if err != nil {
		return err // Kembalikan error jika prepare statement gagal
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.TransactionID, item.VoucherID, item.Quantity, item.CostInPoints)
	return err // Kembalikan error jika eksekusi query gagal
}

func (s *Store) UpdateCustomerTotalPoints(tx *sql.Tx, customerID int, newTotalPoints int) error {
	query := `
        UPDATE customers
        SET total_points = ?
        WHERE id = ?
    `

	var stmt *sql.Stmt
	var err error

	if tx != nil {
		stmt, err = tx.Prepare(query)
	} else {
		stmt, err = s.db.Prepare(query)
	}

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newTotalPoints, customerID)
	return err
}

// Signature function mengembalikan *types.TransactionDetails
func (s *Store) GetTransactionWithItemsByID(id int) (*types.TransactionDetails, error) {
	query := `
    SELECT 
        t.id,
        t.customer_id,
        t.transaction_date,
        t.total_points,
        t.created_at,
        t.updated_at,
        (
            SELECT JSON_ARRAYAGG(
                JSON_OBJECT(
                    'id', ti.id,
                    'transaction_id', ti.transaction_id,
                    'voucher_id', ti.voucher_id,
                    'quantity', ti.quantity,
                    'cost_in_points', ti.cost_in_points
                )
            )
            FROM transaction_items ti
            WHERE ti.transaction_id = t.id
        ) AS transaction_items
    FROM transactions t
    WHERE t.id = ?
    `

	row := s.db.QueryRow(query, id)
	var transaction types.TransactionDetails
	var itemsJSON sql.NullString

	err := row.Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.TransactionDate,
		&transaction.TotalPoints,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&itemsJSON,
	)
	if err != nil {
		return nil, err
	}

	if itemsJSON.Valid {
		err = json.Unmarshal([]byte(itemsJSON.String), &transaction.TransactionItems)
		if err != nil {
			return nil, err
		}
	}

	return &transaction, nil
}
