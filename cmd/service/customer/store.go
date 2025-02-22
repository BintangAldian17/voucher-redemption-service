package customer

import (
	"database/sql"
	"fmt"

	"github.com/BintangAldian17/voucher-redemption-service/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetCustomerByID(id int) (*types.Customer, error) {
	rows, err := s.db.Query("SELECT * FROM customers WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	fmt.Println("rows", id)
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("customer with id %d not found", id)
	}

	customer, err := scanRowIntoCustomer(rows)
	fmt.Println("ashiap", customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func scanRowIntoCustomer(rows *sql.Rows) (*types.Customer, error) {
	customer := new(types.Customer)

	err := rows.Scan(
		&customer.ID,
		&customer.Name,
		&customer.TotalPoints,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return customer, nil
}
