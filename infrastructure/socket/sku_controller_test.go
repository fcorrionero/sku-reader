//go:build integration
// +build integration

package socket_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sku-reader/application"
	"sku-reader/application/mock"
	"sku-reader/infrastructure/socket"
	"testing"
	"time"
)

const (
	PROTOCOL     = "tcp"
	HOST         = "localhost"
	PORT         = "5000"
	END_SEQUENCE = "terminate"
)

func TestShouldHandleConnections(t *testing.T) {

	listener, err := net.Listen(PROTOCOL, HOST+":"+PORT)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	repository := mock.MessageRepositoryImplementationWithoutErrors{}
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	finishReading := make(chan bool)
	errorStream := make(chan interface{})
	defer func() { close(finishReading) }()

	createMessageCommandHandler := application.NewCreateMessageCommandHandler(repository)
	generateReportQueryHandler := application.NewGenerateReportQueryHandler(repository)

	controller := socket.NewSkuController(createMessageCommandHandler, generateReportQueryHandler, listener, ctx)

	sessionId := "test"
	go controller.HandleConnections(sessionId, END_SEQUENCE, finishReading, errorStream)

	conn, err := net.Dial(PROTOCOL, HOST+":"+PORT)
	if err != nil {
		t.Fatalf("expected no error connecting to %v,%v:%v ", PROTOCOL, HOST, PORT)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()
	b := []byte("SKUT-1111\n")
	_, err = conn.Write(b)
	if err != nil {
		t.Fatalf("expected no error writing to client")
	}

}

func TestShouldEndWhenEndSequenceIsPresent(t *testing.T) {
	listener, err := net.Listen(PROTOCOL, HOST+":"+PORT)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	repository := mock.MessageRepositoryImplementationWithoutErrors{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second) // To avoid blocking execution
	defer cancelCtx()
	finishReading := make(chan bool)
	errorStream := make(chan interface{})

	createMessageCommandHandler := application.NewCreateMessageCommandHandler(repository)
	generateReportQueryHandler := application.NewGenerateReportQueryHandler(repository)

	controller := socket.NewSkuController(createMessageCommandHandler, generateReportQueryHandler, listener, ctx)

	sessionId := "test"
	go controller.HandleConnections(sessionId, END_SEQUENCE, finishReading, errorStream)

	conn, err := net.Dial(PROTOCOL, HOST+":"+PORT)
	if err != nil {
		t.Fatalf("expected no error connecting to %v,%v:%v ", PROTOCOL, HOST, PORT)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	b := []byte(END_SEQUENCE + "\n")
	_, err = conn.Write(b)
	for {
		select {
		case <-finishReading:
			return
		case <-ctx.Done():
			close(finishReading)
			t.Fatalf("expected to end test when end sequence: %v is written in host", END_SEQUENCE)
		}

	}
}

func TestShouldEndWhenErrorIsPresent(t *testing.T) {
	listener, err := net.Listen(PROTOCOL, HOST+":"+PORT)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	repository := mock.MessageRepositoryImplementationWithErrors{}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second) // To avoid blocking execution
	defer cancelCtx()
	finishReading := make(chan bool)
	errorStream := make(chan interface{})

	createMessageCommandHandler := application.NewCreateMessageCommandHandler(repository)
	generateReportQueryHandler := application.NewGenerateReportQueryHandler(repository)

	controller := socket.NewSkuController(createMessageCommandHandler, generateReportQueryHandler, listener, ctx)

	sessionId := "test"
	go controller.HandleConnections(sessionId, END_SEQUENCE, finishReading, errorStream)

	conn, err := net.Dial(PROTOCOL, HOST+":"+PORT)
	if err != nil {
		t.Fatalf("expected no error connecting to %v,%v:%v ", PROTOCOL, HOST, PORT)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	b := []byte("SKUT-1111" + "\n")
	_, err = conn.Write(b)
	for {
		select {
		case <-errorStream:
			close(finishReading)
			return
		case <-ctx.Done():
			t.Fatalf("expected to end test when error is present")
		}

	}
}

func TestShouldGenerateReport(t *testing.T) {
	listener, err := net.Listen(PROTOCOL, HOST+":"+PORT)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v", err.Error())
		}
	}(listener)
	if err != nil {
		log.Fatalf("error listening: %v", err.Error())
	}

	repository := mock.MessageRepositoryImplementationWithMessages{}
	ctx, cancelCtx := context.WithCancel(context.Background()) // To avoid blocking execution
	defer cancelCtx()

	createMessageCommandHandler := application.NewCreateMessageCommandHandler(repository)
	generateReportQueryHandler := application.NewGenerateReportQueryHandler(repository)

	controller := socket.NewSkuController(createMessageCommandHandler, generateReportQueryHandler, listener, ctx)

	actualReport := controller.GenerateReport("id")
	expectedReport := fmt.Sprintf("Received %d unique product skus, %d duplicates, %d discard values", 1, 1, 4)

	if actualReport != expectedReport {
		t.Fatalf("expected: %v \n got: %v", expectedReport, actualReport)
	}

	err = os.Remove("skus.log")
	if err != nil {
		t.Fatalf("skus.log file should be present")
	}

}
