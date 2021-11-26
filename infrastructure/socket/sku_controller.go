package socket

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sku-reader/application"
)

type SkuController struct {
	createMessageCommandHandler application.CreateMessageCommandHandler
	generateReportQueryHandler  application.GenerateReportQueryHandler
	listener                    net.Listener
	ctx                         context.Context
	cancel                      context.CancelFunc
}

func NewSkuController(
	createMessageCommandHandler application.CreateMessageCommandHandler,
	generateReportQueryHandler application.GenerateReportQueryHandler,
	listener net.Listener,
	ctx context.Context,
	cancel context.CancelFunc,
) SkuController {
	return SkuController{
		createMessageCommandHandler: createMessageCommandHandler,
		generateReportQueryHandler:  generateReportQueryHandler,
		listener:                    listener,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (controller SkuController) HandleConnections(sessionId string, endSequence string) {
	conn, err := controller.listener.Accept()
	if err != nil {
		return
	}
	log.Println("Client " + conn.RemoteAddr().String() + " connected.")

	for {
		select {
		case <-controller.ctx.Done():
			return
		default:
			buffer, err := bufio.NewReader(conn).ReadString('\n')

			if err != nil {
				err := conn.Close()
				if err != nil {
					log.Fatal(err)
					return
				}

				controller.HandleConnections(sessionId, endSequence)
				return
			}

			log.Println("client message:", buffer[:len(buffer)-1])
			createMessageCommand := application.CreateMessageCommand{
				SessionId: sessionId,
				SKU:       buffer[:len(buffer)-1],
			}
			err = controller.createMessageCommandHandler.Handle(createMessageCommand)
			if err != nil {
				log.Fatalf("error creating message: %v", err.Error())
				return
			}

			if buffer[:len(buffer)-1] == endSequence {
				controller.cancel()
			}
		}
	}
}

func (controller SkuController) GenerateReport(sessionId string) {
	query := application.GenerateReportQuery{SessionId: sessionId}
	reportDto := controller.generateReportQueryHandler.Handle(query)

	log.Println(fmt.Sprintf("Received %d uniqueproduct skus, %d duplicates, %d discard values", reportDto.Received, reportDto.Unique, reportDto.Discarded))
}
