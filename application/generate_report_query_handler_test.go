//go:build unit
// +build unit

package application_test

import (
	"sku-reader/application"
	"sku-reader/application/mock"
	"testing"
	"time"
)

/*
 As it is required to only use standard library we have to create implementations of the repositories interfaces
otherwise we will use mocks created with GoMock library for example (https://github.com/golang/mock)
*/

func TestGenerateReport(t *testing.T) {
	repo := mock.MessageRepositoryImplementationWithMessages{}
	queryHandler := application.NewGenerateReportQueryHandler(repo)

	expectedReportDTO := application.ReportDTO{
		Received:  6,
		Unique:    1,
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
	if len(report.Skus) != report.Unique {
		t.Fatalf("expected 1 sku, got: %v", len(report.Skus))
	}
}
