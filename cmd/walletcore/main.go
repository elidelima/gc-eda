package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com.br/elidelima/go-eda/internal/database"
	"github.com.br/elidelima/go-eda/internal/event"
	"github.com.br/elidelima/go-eda/internal/usecase/create_account"
	"github.com.br/elidelima/go-eda/internal/usecase/create_client"
	"github.com.br/elidelima/go-eda/internal/usecase/create_transaction"
	"github.com.br/elidelima/go-eda/internal/web"
	"github.com.br/elidelima/go-eda/internal/web/webserver"
	"github.com.br/elidelima/go-eda/pkg/events"
	"github.com.br/elidelima/go-eda/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		uow, eventDispatcher, transactionCreatedEvent,
	)

	port := ":3000"
	webserver := webserver.NewWebServer(port)

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}
