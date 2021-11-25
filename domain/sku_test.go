package domain_test

import (
	"errors"
	"sku-reader/domain"
	"testing"
)

func TestSKUShouldBeCreated(t *testing.T) {
	validSKUString := "VDFR-3467"
	sku, err := domain.NewSKU(validSKUString)
	if nil != err {
		t.Fatalf("Error creating sku: %v", err)
	}

	if sku.Value() != validSKUString {
		t.Fatalf("sku value didn't match, expected: %v, got: %v", validSKUString, sku.Value())
	}
}

func TestSKUShouldThrowAnErrorWhenInvalidStringIsProvided(t *testing.T) {
	type invalidSKUs struct {
		input string
		err   error
	}
	tests := []invalidSKUs{
		{input: "34823-KDID", err: errors.New("sku is invalid")},
		{input: "", err: errors.New("sku is invalid")},
		{input: "KIJDASDDAS", err: errors.New("sku is invalid")},
		{input: "123123123", err: errors.New("sku is invalid")},
	}

	for _, tc := range tests {
		sku, err := domain.NewSKU(tc.input)
		if len(sku.Value()) != 0 {
			t.Fatalf("Expected empty string, got: %v", sku.Value())
		}
		if err.Error() != tc.err.Error() {
			t.Fatalf("expected: %v, got: %v", tc.err.Error(), err.Error())
		}
	}
}
