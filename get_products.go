package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/howeyc/gopass"
	_ "github.com/lib/pq"
)

func getPassword(prompt string) string {
	fmt.Print(prompt)
	passwordBytes, err := gopass.GetPasswd()
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

func main() {
	corePassword := getPassword("Core db password: ")
	rcPassword := getPassword("RC db password: ")
	fmt.Println()

	username := "postgres"
	host := "dev-core.cuthntaqfsrx.us-east-1.rds.amazonaws.com"
	database := "dev_avant"
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", username, corePassword, host, database)
	coreDb := getDbConnection("postgres", connectionString)
	defer coreDb.Close()

	fmt.Println("Getting core products...")
	coreRows, err := coreDb.Query("select productid, coalesce(name, '') as name from product order by productid")
	if err != nil {
		log.Fatal(err)
	}
	defer coreRows.Close()

	var productID int
	var name string
	for coreRows.Next() {
		err := coreRows.Scan(&productID, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(productID, name)
	}
	err = coreRows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	username = "avant"
	host = "mssql-development.cuthntaqfsrx.us-east-1.rds.amazonaws.com"
	database = "IIRWIN-RC"
	connectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", username, rcPassword, host, database)
	rcDb := getDbConnection("sqlserver", connectionString)
	defer rcDb.Close()

	fmt.Println("Getting RC products...")
	rcRows, err := rcDb.Query("select ProductID, coalesce(ProductCode, '') as ProductCode from Product order by ProductID")
	if err != nil {
		log.Fatal(err)
	}
	defer rcRows.Close()

	for rcRows.Next() {
		err := rcRows.Scan(&productID, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(productID, name)
	}
	err = rcRows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone!")
}
