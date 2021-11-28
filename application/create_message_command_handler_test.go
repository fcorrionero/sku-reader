//go:build unit
// +build unit

package application_test

import (
	"sku-reader/application"
	"sku-reader/application/mock"
	"testing"
	"time"
)

func TestCreateMessageCommandHandlerWithoutErrors(t *testing.T) {
	command := application.CreateMessageCommand{
		SessionId: time.Now().String(),
		SKU:       "VDFR-4455",
	}

	repository := mock.MessageRepositoryImplementationWithoutErrors{}
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

	repository := mock.MessageRepositoryImplementationWithErrors{}
	handler := application.NewCreateMessageCommandHandler(repository)

	err := handler.Handle(command)

	if err == nil {
		t.Fatalf("Expected error not nil, got a nil error")
	}

}
