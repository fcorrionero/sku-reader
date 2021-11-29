package application

import (
	"log"
	"os"
	"sku-reader/domain"
)

type GenerateReportQuery struct {
	SessionId string
}

type GenerateReportQueryHandler struct {
	messageRepository domain.MessageRepository
}

type ReportDTO struct {
	Received  int
	Unique    int
	Discarded int
}

func NewGenerateReportQueryHandler(messageRepository domain.MessageRepository) GenerateReportQueryHandler {
	return GenerateReportQueryHandler{
		messageRepository: messageRepository,
	}
}

func (handler *GenerateReportQueryHandler) Handle(query GenerateReportQuery) ReportDTO {
	report := ReportDTO{
		Received:  0,
		Unique:    0,
		Discarded: 0,
	}

	var skus []string

	messages := handler.messageRepository.FindAll(query.SessionId)

	for _, m := range messages {
		unique := handler.skuUnique(m.Sku, skus)
		if unique && !m.Discard {
			report.Unique++
			skus = append(skus, m.Sku)
		}
		if m.Discard {
			report.Discarded++
		}
		report.Received++
	}

	err := handler.generateLogFile(skus)
	if err != nil {
		log.Println("error generating skus log file: " + err.Error())
	}

	return report
}

func (handler *GenerateReportQueryHandler) skuUnique(sku string, skus []string) bool {
	for _, s := range skus {
		if s == sku {
			return false
		}
	}

	return true
}

func (handler *GenerateReportQueryHandler) generateLogFile(skus []string) error {
	file, err := os.OpenFile("skus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	log.SetOutput(file)
	for _, sku := range skus {
		log.Println(sku)
	}

	return nil
}
