package create_account

import (
	"testing"

	"github.com.br/elidelima/go-eda/internal/entity"
	"github.com.br/elidelima/go-eda/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("Ze Ninguem", "z@n")
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountMock, clientMock)
	inputDto := CreateAccountInputDTO{
		ClientID: client.ID,
	}
	outputDto, err := uc.Execute(inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, outputDto)
	assert.NotEmpty(t, outputDto, outputDto.ID)
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
