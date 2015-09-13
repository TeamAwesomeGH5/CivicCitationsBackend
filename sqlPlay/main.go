package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:my-secret-pw@/?charset=utf8")
	defer db.Close()
	db.Exec("CREATE DATABASE civicCitations")
	db.Exec("USE civicCitations")

	createForFile("citations.csv")
	fmt.Println("Finished citations, starting violations")
	createForFile("violations.csv")
	fmt.Println("Finished violations")
}

func createForFile(filename string) {
	table := "citations"
	if filename == "violations.csv" {
		table = "violations"
	}
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvFile := csv.NewReader(file)
	headers, err := csvFile.Read()
	if err != nil {
		log.Fatalf("Couldn't read header row: %s", err)
	}

	db, err := sql.Open("mysql", "root:my-secret-pw@/civicCitations?charset=utf8")
	if err != nil {
		log.Fatalf("There was an issue opening the database: %s", err)
	}
	defer db.Close()

	createTableSql := createTableStatement(filename, headers)
	fmt.Println(createTableSql)
	stmt, err := db.Prepare(createTableSql)
	if err != nil {
		log.Fatalf("There was an error running the statment: %s", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("error creating table %s: %s", createTableSql, err)
	}
	fmt.Println("Adding data to rows.")

	data, err := csvFile.ReadAll()
	if err != nil {
		log.Fatalf("There was an error reading all the records: %s", err)
	}
	for _, dataRow := range data {
		sql := fmt.Sprintf("INSERT into %s (%s) VALUES (\"%s\")", table, strings.Join(headers, ","), strings.Join(dataRow, "\",\""))
		statement, err := db.Prepare(sql)
		if err != nil {
			log.Fatalf("There was an error statement: %s\ninserting data %v\nError: %s", sql, dataRow, err)
		}
		statement.Exec()
	}

	//To create the table
	// db.Exec("CREATE DATABASE civicCitations")
	// db.Exec("USE civicCitations")
	// stmt, err := db.Prepare(createTableSql)
	// if err != nil {
	// 	log.Fatalf("There was an error running the statment: %s", err)
	// }
	// stmt.Exec()
}

func createTableStatement(filename string, headers []string) string {
	table := "citations"
	index := "cid"
	if filename == "violations.csv" {
		table = "violations"
		index = "vid"
	}
	statement := fmt.Sprintf("CREATE TABLE `%s` (\n", table)
	statement += fmt.Sprintf("`%s` INT(10) NOT NULL AUTO_INCREMENT,\n", index)
	for _, header := range headers {
		statement += fmt.Sprintf("`%s` VARCHAR(255) NULL DEFAULT NULL,\n", header)
	}
	statement += fmt.Sprintf("PRIMARY KEY(`%s`)\n", index)
	statement += ");"
	return statement
}
