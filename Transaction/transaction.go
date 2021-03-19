package transaction

import (
	"crypto/rsa"
	"encoding/json"
)

type Transaction struct {
	Amount int
	Payer, Payee *rsa.PublicKey
}

func NewTransaction(a int, payer, payee *rsa.PublicKey) *Transaction{
	return &Transaction{Amount: a, Payer: payer, Payee: payee}
}

func (t *Transaction) ToString() ([]byte, error) {
	return json.Marshal(t)
}