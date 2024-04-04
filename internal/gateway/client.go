package gateway

import "github.com.br/elidelima/go-eda/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
