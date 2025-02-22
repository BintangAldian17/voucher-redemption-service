package main

import (
	"fmt"
	"log"

	"github.com/BintangAldian17/voucher-redemption-service/configs"
	"github.com/BintangAldian17/voucher-redemption-service/db"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysqlDriver.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dbConn, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	seedQueries := []string{
		`INSERT INTO brands (name, description) VALUES 
			('Brand A', 'Description for Brand A'),
			('Brand B', 'Description for Brand B');`,
		`INSERT INTO vouchers (brand_id, name, description, cost_in_points) VALUES 
			(1, 'Voucher A1', '10% off discount', 100),
			(1, 'Voucher A2', '20% off discount', 200),
			(2, 'Voucher B1', 'Free shipping voucher', 50),
			(2, 'Voucher B2', '10% off voucher', 150);`,
		`INSERT INTO customers (name, total_points) VALUES 
			('Customer One', 1000),
			('Customer Two', 500),
			('Customer Three', 750);`,
	}

	for _, query := range seedQueries {
		if _, err := dbConn.Exec(query); err != nil {
			log.Fatalf("Error executing seed query: %v\nQuery: %s", err, query)
		}
	}

	fmt.Println("Seeder executed successfully!")
}
