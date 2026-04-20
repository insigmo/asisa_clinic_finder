package main

import (
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFileName = "db/local.db"
	dbName     = "sqlite3"
)

func main() {
	saveCities()
	saveDirections()
}

func saveCities() {
	data := readJSON[map[string][]int]("cmd/data_migrator/data/cities.json")

	insertInTx(
		`INSERT INTO city_postal_codes (city, postal_code) VALUES (?, ?)`,
		func(exec func(args ...any)) {
			for city, codes := range data {
				for _, code := range codes {
					exec(city, code)
				}
			}
		},
	)
}

func saveDirections() {
	data := readJSON[map[string]string]("cmd/data_migrator/data/directions.json")

	insertInTx(
		`INSERT INTO medical_direction (reference_name, name) VALUES (?, ?)`,
		func(exec func(args ...any)) {
			for reference_name, name := range data {
				exec(reference_name, name)
			}
		},
	)
}

// readJSON читает JSON-файл в значение типа T.
func readJSON[T any](path string) T {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data T
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		panic(err)
	}

	return data
}

func insertInTx(query string, feed func(exec func(args ...any))) {
	client, err := sql.Open(dbName, dbFileName)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	tx, err := client.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	exec := func(args ...any) {
		if _, err := stmt.Exec(args...); err != nil {
			panic(err)
		}
	}

	feed(exec)

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
