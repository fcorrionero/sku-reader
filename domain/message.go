package domain

type Message struct {
	sessionId string
	sku       string
	discard   bool
}

func NewMessage(id string, skuValue string) Message {
	_, err := NewSKU(skuValue)
	if err != nil {
		return Message{
			sessionId: id,
			sku:       skuValue,
			discard:   true,
		}
	}

	return Message{
		sessionId: id,
		sku:       skuValue,
		discard:   false,
	}
}

func (m *Message) SessionId() string {
	return m.sessionId
}

func (m *Message) SKU() string {
	return m.sku
}

func (m *Message) Discard() bool {
	return m.discard
}
