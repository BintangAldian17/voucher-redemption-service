package types

import (
	"database/sql"
	"time"
)

type Brand struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Voucher struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	BrandID      int       `json:"brand_id"`
	CostInPoints int       `json:"cost_in_points"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Customer struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TotalPoints int       `json:"total_points"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Transaction struct {
	ID              int       `json:"id"`
	CustomerID      int       `json:"customer_id"`
	TransactionDate time.Time `json:"transaction_date"`
	TotalPoints     int       `json:"total_points"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type TransactionItem struct {
	ID            int       `json:"id"`
	TransactionID int       `json:"transaction_id"`
	VoucherID     int       `json:"voucher_id"`
	Quantity      int       `json:"quantity"`
	CostInPoints  int       `json:"cost_in_points"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateVoucherPayload struct {
	Name         string `json:"name" validate:"required"`
	BrandID      int    `json:"brand_id" validate:"required"`
	CostInPoints int    `json:"cost_in_points" validate:"required,min=1"`
	Description  string `json:"description,omitempty"`
}

type CreateBrandPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
}

type RedemptionPayload struct {
	CustomerID int              `json:"customer_id" validate:"required"`
	Vouchers   []RedemptionItem `json:"vouchers" validate:"required,min=1"`
}

type RedemptionItem struct {
	VoucherID int `json:"voucher_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type TransactionDetails struct {
	ID               int               `json:"id"`
	CustomerID       int               `json:"customer_id"`
	TransactionDate  time.Time         `json:"transaction_date"`
	TotalPoints      int               `json:"total_points"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	TransactionItems []TransactionItem `json:"transaction_items"`
}

type TransactionStore interface {
	RedeemVouchers(tx *sql.Tx, customer *Customer, transactionItems []TransactionItem, totalPointsNeeded int) (int, error)
	Begin() (*sql.Tx, error)
	GetTransactionWithItemsByID(id int) (*TransactionDetails, error)
}
type TransactionItemStore interface{}

type CustomerStore interface {
	GetCustomerByID(id int) (*Customer, error)
}

type VoucherStore interface {
	GetVoucherByName(name string) (*Voucher, error)
	GetVoucherByID(id int) (*Voucher, error)
	GetVouchersByBrandID(brandID int) ([]*Voucher, error)
	CreateVoucher(brand Voucher) error
}

type BrandStore interface {
	GetBrandByName(name string) (*Brand, error)
	CreateBrand(brand Brand) error
}
