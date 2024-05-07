package database_config

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	// Define SQL statements for dropping and creating tables
	tables := []string{
		`
			DROP TABLE IF EXISTS clients;
		`,
		`
			CREATE TABLE clients (
				id VARCHAR(255), 
				name VARCHAR(255), 
				email VARCHAR(255), 
				created_at DATE
			);
		`,
		`
			DROP TABLE IF EXISTS accounts;
		`,
		`
			CREATE TABLE accounts (
				id VARCHAR(255),
				client_id VARCHAR(255),
				balance INT,
				created_at DATE
			)
		`,
		`
			DROP TABLE IF EXISTS transactions;
		`,
		`
			CREATE TABLE transactions (
				id VARCHAR(255),
				account_id_from VARCHAR(255),
				account_id_to VARCHAR(255),
				amount INT, 
				created_at DATE
			)
		`,
	}

	for _, createTableStmt := range tables {
		_, err := db.Exec(createTableStmt)
		if err != nil {
			return err
		}
	}
	return nil
}
