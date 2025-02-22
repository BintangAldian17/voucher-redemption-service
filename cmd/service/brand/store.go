package brand

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

func (s *Store) CreateBrand(brand types.Brand) error {

	_, err := s.db.Exec("INSERT INTO brands (name, description) VALUES(?, ?)", brand.Name, brand.Description)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBrandByName(name string) (*types.Brand, error) {
	rows, err := s.db.Query("SELECT * FROM brands WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("brand dengan nama %s tidak ditemukan", name)
	}

	brand, err := scanRowIntoBrand(rows)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brand, nil
}

func scanRowIntoBrand(rows *sql.Rows) (*types.Brand, error) {
	brand := new(types.Brand)

	err := rows.Scan(
		&brand.ID,
		&brand.Name,
		&brand.Description,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return brand, nil
}
