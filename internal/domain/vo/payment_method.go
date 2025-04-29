package vo

import (
	"errors"
	"strings"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
)

func (p PaymentMethod) Valid() bool {
	switch p {
	case PaymentMethodCreditCard:
		return true
	default:
		return false
	}
}

func NewPaymentMethod(method string) (PaymentMethod, error) {
	pm := PaymentMethod(strings.ToLower(method))
	if !pm.Valid() {
		return "", errors.New("invalid payment method")
	}
	return pm, nil
}

func (p PaymentMethod) String() string {
	return string(p)
}
