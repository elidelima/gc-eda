package gateway

import "github.com.br/elidelima/go-eda/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
