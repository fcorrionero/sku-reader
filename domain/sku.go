package domain

import (
	"errors"
	"regexp"
)

const patternValidator = "^[a-zA-Z]{4}-[0-9]{4}$"

type SKU struct {
	value string
}

func NewSKU(value string) (SKU, error) {
	valid, _ := regexp.MatchString(patternValidator, value)
	if !valid {
		return SKU{}, errors.New("sku is invalid")
	}
	sku := SKU{value: value}
	return sku, nil
}

func (sku *SKU) Value() string {
	return sku.value
}
