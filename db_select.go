package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
)

func main() {
	db := connectDB()
	defer db.Close()

	query := buildQuery(os.Args)
	sliceOfMaps := execQuery(db, query)
	printFormatted(sliceOfMaps)
}

type Query struct {
	Table   string
	Columns string
	Field   string
	Value   string
}

func buildQuery(attrs []string) Query {
	query := Query{
		Table:   "",
		Columns: "*",
	}
	for i, v := range attrs {
		switch i {
		case 1:
			query.Table = v
		case 2:
			query.Columns = v
		case 3:
			query.Field = v
		case 4:
			query.Value = v
		}
	}

	return query
}

func execQuery(db *sql.DB, query Query) []map[string]interface{} {
	var err error = nil
	var stmt *sql.Stmt
	if query.Field != "" && query.Value != "" {
		stmt, err = db.Prepare("SELECT " + query.Columns + " FROM " + query.Table + " WHERE " + query.Field + " = ? LIMIT 10")
	}
	stmt, err = db.Prepare("SELECT " + query.Columns + " FROM " + query.Table + " LIMIT 10")
	CheckError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	CheckError(err)
	defer rows.Close()

	return scanRowsToMaps(rows)
}

func printFormatted(list []map[string]interface{}) {
	for _, v := range list {
		for key, value := range v {
			fmt.Printf("%v:", key)
			color.Set(color.FgGreen)
			fmt.Printf(" %v\n", value)
			color.Unset()
		}
		fmt.Println("=============================================")
	}
}

func scanRowsToMaps(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var results []map[string]interface{}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("Error scanning rows:", err)
			return results
		}

		dataMap := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			dataMap[colName] = val
		}

		results = append(results, dataMap)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error reading rows:", err)
	}

	return results
}

func connectDB() *sql.DB {
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
