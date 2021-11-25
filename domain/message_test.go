package domain_test

import (
	"sku-reader/domain"
	"testing"
	"time"
)

func TestMessagesShouldBeCreated(t *testing.T) {
	id := time.Now().String()
	type SKUs struct {
		sku     string
		discard bool
	}
	tests := []SKUs{
		{sku: "VDFR-3467", discard: false},
		{sku: "34823-KDID", discard: true},
		{sku: "", discard: true},
		{sku: "KIJDASDDAS", discard: true},
		{sku: "123123123", discard: true},
	}

	for _, tc := range tests {
		m := domain.NewMessage(id, tc.sku)
		if m.SessionId() != id {
			t.Fatalf("sessionId invalid, expected: %v, got: %v", id, m.SessionId())
		}
		if m.SKU() != tc.sku {
			t.Fatalf("sku invalid, expected: %v, got: %v", tc.sku, m.SKU())
		}
		if m.Discard() != tc.discard {
			t.Fatalf("discard invalid, expected: %v, got: %v", tc.discard, m.Discard())
		}
	}
}
