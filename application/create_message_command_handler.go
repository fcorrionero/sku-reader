package application

import "sku-reader/domain"

type CreateMessageCommand struct {
	SessionId string
	SKU       string
}

type CreateMessageCommandHandler struct {
	messageRepository domain.MessageRepository
}

func NewCreateMessageCommandHandler(repository domain.MessageRepository) CreateMessageCommandHandler {
	return CreateMessageCommandHandler{
		messageRepository: repository,
	}
}

func (handler *CreateMessageCommandHandler) Handle(command CreateMessageCommand) error {
	message := domain.NewMessage(command.SessionId, command.SKU)
	return handler.messageRepository.Save(&message)
}
