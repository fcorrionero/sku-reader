package controller

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sku-reader/application"
	"time"
)

type SkuController struct {
	createMessageCommandHandler application.CreateMessageCommandHandler
	generateReportQueryHandler  application.GenerateReportQueryHandler
	listener                    net.Listener
	ctx                         context.Context
}

func NewSkuController(
	createMessageCommandHandler application.CreateMessageCommandHandler,
	generateReportQueryHandler application.GenerateReportQueryHandler,
	listener net.Listener,
	ctx context.Context,
) SkuController {
	return SkuController{
		createMessageCommandHandler: createMessageCommandHandler,
		generateReportQueryHandler:  generateReportQueryHandler,
		listener:                    listener,
		ctx:                         ctx,
	}
}

func (controller *SkuController) HandleConnections(sessionId string, endSequence string, finishReading chan bool, errorStream chan interface{}) {
	conn, err := controller.listener.Accept()
	if err != nil {
		return
	}
	log.Println("Client " + conn.RemoteAddr().String() + " connected.")

	for {
		select {
		case <-finishReading:
			return
		default:
			buffer, err := bufio.NewReader(conn).ReadString('\n')

			if err != nil {
				err := conn.Close()
				if err != nil {
					errorStream <- err
					return
				}

				controller.HandleConnections(sessionId, endSequence, finishReading, errorStream)
				return
			}

			log.Println("client message:", buffer[:len(buffer)-1])

			if buffer[:len(buffer)-1] == endSequence {
				close(finishReading)
				return
			}

			createMessageCommand := application.CreateMessageCommand{
				SessionId: sessionId,
				SKU:       buffer[:len(buffer)-1],
			}
			err = controller.createMessageCommandHandler.Handle(createMessageCommand)
			if err != nil {
				errorStream <- err
				return
			}
		}
	}
}

func (controller *SkuController) GenerateReport(sessionId string) string {
	query := application.GenerateReportQuery{SessionId: sessionId}
	reportDto := controller.generateReportQueryHandler.Handle(query)
	err := controller.generateLogFile(reportDto.Skus)
	if err != nil {
		log.Println("error writting skus to log file : " + err.Error())
	}

	return fmt.Sprintf("Received %d unique product skus, %d duplicates, %d discard values", reportDto.Unique, reportDto.Received-reportDto.Unique-reportDto.Discarded, reportDto.Discarded)
}

func (controller *SkuController) generateLogFile(skus []string) error {
	file, err := os.OpenFile("skus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, sku := range skus {
		_, err = file.WriteString(fmt.Sprintf("%s\t%s\n", time.Now().Format(time.RFC3339), sku))
		if err != nil {
			log.Println("error writting skus to log file : " + err.Error())
		}
	}

	return nil
}
