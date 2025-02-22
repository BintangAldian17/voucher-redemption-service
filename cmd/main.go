package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/BintangAldian17/voucher-redemption-service/cmd/api"
	"github.com/BintangAldian17/voucher-redemption-service/configs"
	"github.com/BintangAldian17/voucher-redemption-service/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	port := fmt.Sprintf(":%s", configs.Envs.Port)

	server := api.NewAPIServer(port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
