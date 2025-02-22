package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/BintangAldian17/voucher-redemption-service/cmd/service/brand"
	"github.com/BintangAldian17/voucher-redemption-service/cmd/service/customer"
	"github.com/BintangAldian17/voucher-redemption-service/cmd/service/transaction"
	"github.com/BintangAldian17/voucher-redemption-service/cmd/service/voucher"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	brandStore := brand.NewStore(s.db)
	voucherStore := voucher.NewStore(s.db)
	transactionStore := transaction.NewStore(s.db)
	customerStore := customer.NewStore(s.db)

	brandHandler := brand.NewHandler(brandStore)
	voucherHandler := voucher.NewHandler(voucherStore)
	transactionHandler := transaction.NewHandler(transactionStore, customerStore, voucherStore)

	voucherHandler.RegisterRoutes(subrouter)
	brandHandler.RegisterRoutes(subrouter)
	transactionHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
