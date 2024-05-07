package database_config

import (
	"database/sql"
	"time"
)

func SeedData(db *sql.DB) error {
	// Define initial data to insert into the 'users' table
	clients := [][]string{
		{
			"0eddd18b-1f78-4932-a13f-3b04c083056d",
			"José Ninguém",
			"jn@email.com",
		},
		{
			"dc2fcaf6-ff57-41c8-b963-a6880d9fea57",
			"Maria Ninguém",
			"mn@email.com",
		},
	}

	accountsIds := []string{
		"2161e2cf-27ba-46f2-aab7-950c00dacabf",
		"b06936c1-89d4-49dc-a480-c7381e25b582",
	}

	// Prepare SQL statement for inserting data
	clientInsertStmt := "INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)"
	accountInsertStmt := "INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)"

	// Iterate over initial data and execute SQL insert statements
	for index, user := range clients {
		_, err := db.Exec(clientInsertStmt, user[0], user[1], user[2], time.Now())
		if err != nil {
			return err
		}

		_, err = db.Exec(accountInsertStmt, accountsIds[index], user[0], 1000, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
