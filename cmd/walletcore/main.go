package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com.br/elidelima/go-eda/database_config"
	"github.com.br/elidelima/go-eda/internal/database"
	"github.com.br/elidelima/go-eda/internal/event"
	"github.com.br/elidelima/go-eda/internal/event/handler"
	"github.com.br/elidelima/go-eda/internal/usecase/create_account"
	"github.com.br/elidelima/go-eda/internal/usecase/create_client"
	"github.com.br/elidelima/go-eda/internal/usecase/create_transaction"
	"github.com.br/elidelima/go-eda/internal/web"
	"github.com.br/elidelima/go-eda/internal/web/webserver"
	"github.com.br/elidelima/go-eda/pkg/events"
	"github.com.br/elidelima/go-eda/pkg/kafka"
	"github.com.br/elidelima/go-eda/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func createAndSeedTables(db *sql.DB) {
	err := database_config.CreateTables(db)
	if err != nil {
		panic(err)
	}

	// Seed initial data
	if err := database_config.SeedData(db); err != nil {
		panic(err)
	}

	fmt.Println("Initialization and seeding completed successfully!")
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "wallet-db", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createAndSeedTables(db)

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

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
		uow,
		eventDispatcher,
		transactionCreatedEvent,
		balanceUpdatedEvent,
	)

	port := ":8080"
	webserver := webserver.NewWebServer(port)

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Printf("Server started on port %s", port)
	webserver.Start()
}
