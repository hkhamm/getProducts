package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/howeyc/gopass"
	_ "github.com/lib/pq"
)

func getPassword(prompt string) string {
	passwordBytes, err := gopass.GetPasswdPrompt(prompt, true, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	return string(passwordBytes)
}

func getDbConnection(dbType string, connectionString string) *sql.DB {
	db, err := sql.Open(dbType, connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func printProducts(products *sql.Rows) {
	var productID int
	var name string
	for products.Next() {
		err := products.Scan(&productID, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(productID, name)
	}
	err := products.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func getCoreProducts(password string) {
	username := "postgres"
	host := "dev-core.cuthntaqfsrx.us-east-1.rds.amazonaws.com"
	database := "dev_avant"
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, host, database)
	db := getDbConnection("postgres", connectionString)
	defer db.Close()

	fmt.Println("Getting core products...")
	rows, err := db.Query("select productid, coalesce(name, '') as name from product order by productid")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	printProducts(rows)
}

func getRcProducts(password string) {
	username := "avant"
	host := "mssql-development.cuthntaqfsrx.us-east-1.rds.amazonaws.com"
	database := "IIRWIN-RC"
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, password, host, database)
	db := getDbConnection("sqlserver", connectionString)
	defer db.Close()

	fmt.Println("Getting RC products...")
	rows, err := db.Query("select ProductID, coalesce(ProductCode, '') as ProductCode from Product order by ProductID")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	printProducts(rows)
}

func main() {
	corePassword := getPassword("Core db password: ")
	rcPassword := getPassword("RC db password: ")
	fmt.Println()

	getCoreProducts(corePassword)
	getRcProducts(rcPassword)

	fmt.Println("Done!")
}
