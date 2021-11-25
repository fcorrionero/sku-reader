package application_test

import (
	"sku-reader/application"
	"sku-reader/domain"
	"testing"
	"time"
)

/*
 As it is required to only use standard library we have to create implementations of the repositories interfaces
otherwise we will use mocks created with GoMock library for example (https://github.com/golang/mock)
*/
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

func TestGenerateReport(t *testing.T) {
	repo := MessageRepositoryImplementationWithMessages{}
	queryHandler := application.NewGenerateReportQueryHandler(repo)

	expectedReportDTO := application.ReportDTO{
		Received:  6,
		Unique:    5,
		Discarded: 4,
	}

	id := time.Now().String()
	query := application.GenerateReportQuery{SessionId: id}
	report := queryHandler.Handle(query)

	if report.Unique != expectedReportDTO.Unique {
		t.Fatalf("unique fail, expected: %v, got: %v", expectedReportDTO.Unique, report.Unique)
	}
	if report.Received != expectedReportDTO.Received {
		t.Fatalf("received fail, expected: %v, got: %v", expectedReportDTO.Received, report.Received)
	}
	if report.Discarded != expectedReportDTO.Discarded {
		t.Fatalf("discarded fail, expected: %v, got: %v", expectedReportDTO.Discarded, report.Discarded)
	}
}
