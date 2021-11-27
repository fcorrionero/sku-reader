//go:build unit
// +build unit

package application_test

import (
	"errors"
	"sku-reader/application"
	"sku-reader/domain"
	"testing"
	"time"
)

/*
 As it is required to only use standard library we have to create implementations of the repositories interfaces
otherwise we will use mocks created with GoMock library for example (https://github.com/golang/mock)
*/
type MessageRepositoryImplementationWithoutErrors struct{}

func (r MessageRepositoryImplementationWithoutErrors) Save(message *domain.Message) error {
	return nil
}
func (r MessageRepositoryImplementationWithoutErrors) FindAll(id string) []*domain.Message {
	messages := []*domain.Message{}
	return messages
}

type MessageRepositoryImplementationWithErrors struct{}

func (r MessageRepositoryImplementationWithErrors) Save(message *domain.Message) error {
	return errors.New("")
}
func (r MessageRepositoryImplementationWithErrors) FindAll(id string) []*domain.Message {
	messages := []*domain.Message{}
	return messages
}

func TestCreateMessageCommandHandlerWithoutErrors(t *testing.T) {
	command := application.CreateMessageCommand{
		SessionId: time.Now().String(),
		SKU:       "VDFR-4455",
	}

	repository := MessageRepositoryImplementationWithoutErrors{}
	handler := application.NewCreateMessageCommandHandler(repository)

	err := handler.Handle(command)

	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err.Error())
	}
}

func TestCreateMessageCommandHandlerWithErrors(t *testing.T) {
	command := application.CreateMessageCommand{
		SessionId: time.Now().String(),
		SKU:       "VDFR-4455",
	}

	repository := MessageRepositoryImplementationWithErrors{}
	handler := application.NewCreateMessageCommandHandler(repository)

	err := handler.Handle(command)

	if err == nil {
		t.Fatalf("Expected error not nil, got a nil error")
	}

}
