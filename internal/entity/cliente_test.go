package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("Zé Ninguém", "z@n.com")
	_ = err
	assert.NotNil(t, client)
	assert.Equal(t, "Zé Ninguém", client.Name)
	assert.Equal(t, "z@n.com", client.Email)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	_ = err
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("Zé Ninguém", "z@n.com")
	err := client.Update("Zé Ninguém Updated", "z@n.com")
	assert.Nil(t, err)
	assert.Equal(t, "Zé Ninguém Updated", client.Name)
	assert.Equal(t, "z@n.com", client.Email)
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client, _ := NewClient("Zé Ninguém", "z@n.com")
	err := client.Update("", "z@n.com")
	assert.Error(t, err, "name is required")
}

func TestAddAccount(t *testing.T) {
	client, _ := NewClient("Zé Ninguém", "z@n.com")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}
