package domain

type Message struct {
	SessionId string
	Sku       string
	Discard   bool
}

func NewMessage(id string, skuValue string) Message {
	_, err := NewSKU(skuValue)
	if err != nil {
		return Message{
			SessionId: id,
			Sku:       skuValue,
			Discard:   true,
		}
	}

	return Message{
		SessionId: id,
		Sku:       skuValue,
		Discard:   false,
	}
}
