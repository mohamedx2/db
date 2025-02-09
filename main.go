package main

import (
	"bufio"
	"db/api"
	"db/database"
	"fmt"
	"log"
	"os"
)

func main() {
	db := database.NewDatabase("MyDB")
	server := api.NewServer(db)

	log.Println("Starting server on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enhanced Database System")
	fmt.Println("Commands: create-table, insert, select, update, delete, help, exit")

	for {
		fmt.Print("db> ")
		if !scanner.Scan() {
			break
		}

		cmd := scanner.Text()
		switch cmd {
		case "help":
			printHelp()
		case "delete":
			deleteData(db, scanner)
		case "update":
			updateData(db, scanner)
		case "exit":
			return
		case "create-table":
			createTable(db, scanner)
		case "insert":
			insertData(db, scanner)
		case "select":
			selectData(db, scanner)
		default:
			fmt.Println("Unknown command")
		}
	}
}

func createTable(db *database.Database, scanner *bufio.Scanner) {
	fmt.Print("Table name: ")
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()

	columns := []database.Column{
		{Name: "id", DataType: "int"},
		{Name: "name", DataType: "string"},
		{Name: "active", DataType: "bool"},
	}

	if err := db.CreateTable(name, columns); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Table created successfully")
}

func insertData(db *database.Database, scanner *bufio.Scanner) {
	fmt.Print("Table name: ")
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()

	table, err := db.GetTable(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	row := database.Row{
		"id":     1,
		"name":   "John Doe",
		"active": true,
	}

	if err := table.InsertRow(row); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Row inserted successfully")
}

func selectData(db *database.Database, scanner *bufio.Scanner) {
	fmt.Print("Table name: ")
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()

	table, err := db.GetTable(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	rows, err := table.Select(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, row := range rows {
		fmt.Printf("%v\n", row)
	}
}

func deleteData(db *database.Database, scanner *bufio.Scanner) {
	fmt.Print("Table name: ")
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()

	table, err := db.GetTable(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("WHERE clause (e.g., id=1 AND active=true): ")
	if !scanner.Scan() {
		return
	}
	whereClause := scanner.Text()

	conditions, err := database.ParseWhereClause(whereClause)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	deleted, err := table.Delete(conditions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%d rows deleted\n", deleted)
}

func updateData(db *database.Database, scanner *bufio.Scanner) {
	// Similar to deleteData but with SET clause
	// Implementation details omitted for brevity
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  create-table - Create a new table")
	fmt.Println("  insert      - Insert data into a table")
	fmt.Println("  select      - Query data from a table")
	fmt.Println("  update      - Update existing data")
	fmt.Println("  delete      - Delete data from a table")
	fmt.Println("  help        - Show this help message")
	fmt.Println("  exit        - Exit the program")
}
