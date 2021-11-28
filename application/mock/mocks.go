package mock

import (
	"errors"
	"sku-reader/domain"
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
	return []*domain.Message{}
}

type MessageRepositoryImplementationWithErrors struct{}

func (r MessageRepositoryImplementationWithErrors) Save(message *domain.Message) error {
	return errors.New("")
}
func (r MessageRepositoryImplementationWithErrors) FindAll(id string) []*domain.Message {
	messages := []*domain.Message{}
	return messages
}

type MessageRepositoryImplementationWithMessages struct{}

func (r MessageRepositoryImplementationWithMessages) Save(message *domain.Message) error {
	return nil
}
func (r MessageRepositoryImplementationWithMessages) FindAll(id string) []*domain.Message {
	messages := []*domain.Message{}

	skus := []string{
		"VDFR-3467",
		"VDFR-3467",
		"34823-KDID",
		"",
		"KIJDASDDAS",
		"123123123",
	}

	for _, sku := range skus {
		message := domain.NewMessage(id, sku)
		messages = append(messages, &message)
	}

	return messages
}
