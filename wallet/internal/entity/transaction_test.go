package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("Client1", "client@one")
	account1 := NewAccount(client1)
	account1.Credit(1000)

	client2, _ := NewClient("Client2", "client@two")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	tranaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)
	assert.NotNil(t, tranaction)
	assert.Equal(t, float64(1100), account2.Balance)
	assert.Equal(t, float64(900), account1.Balance)
}

func TestCreateTransactionWithInsufficientFunds(t *testing.T) {
	client1, _ := NewClient("Client1", "client@one")
	account1 := NewAccount(client1)
	account1.Credit(1000)

	client2, _ := NewClient("Client2", "client@two")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	tranaction, err := NewTransaction(account1, account2, 2000)
	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, tranaction)
	assert.Equal(t, float64(1000), account2.Balance)
	assert.Equal(t, float64(1000), account1.Balance)
}
