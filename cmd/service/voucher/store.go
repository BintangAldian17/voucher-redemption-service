package voucher

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

func (s *Store) CreateVoucher(voucher types.Voucher) error {
	_, err := s.db.Exec("INSERT INTO vouchers (name, brand_id, cost_in_points, description) VALUES(?, ?, ?, ?)", voucher.Name, voucher.BrandID, voucher.CostInPoints, voucher.Description)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetVoucherByID(id int) (*types.Voucher, error) {
	rows, err := s.db.Query("SELECT * FROM vouchers WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	fmt.Println("rows", rows)
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("voucher with id %d not found", id)
	}

	voucher, err := scanRowIntoVoucher(rows)
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *Store) GetVoucherByName(name string) (*types.Voucher, error) {
	rows, err := s.db.Query("SELECT * FROM vouchers WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("voucher with name %s not found", name)
	}

	voucher, err := scanRowIntoVoucher(rows)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *Store) GetVouchersByBrandID(brandID int) ([]*types.Voucher, error) {
	rows, err := s.db.Query("SELECT * FROM vouchers WHERE brand_id = ?", brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vouchers := []*types.Voucher{}
	for rows.Next() {
		voucher, err := scanRowIntoVoucher(rows)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, voucher)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if len(vouchers) == 0 {
	// 	return nil, fmt.Errorf("no vouchers found for brand ID: %d", brandID)
	// }

	return vouchers, nil
}

func scanRowIntoVoucher(rows *sql.Rows) (*types.Voucher, error) {
	voucher := new(types.Voucher)

	err := rows.Scan(
		&voucher.ID,
		&voucher.BrandID,
		&voucher.Name,
		&voucher.Description,
		&voucher.CostInPoints,
		&voucher.CreatedAt,
		&voucher.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return voucher, nil
}
