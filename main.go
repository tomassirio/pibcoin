package main

import (
	"github.com/tomassirio/pibcoin/Chain"
	"github.com/tomassirio/pibcoin/Wallet"
	"log"
)

func main(){
	satoshi := Wallet.NewWallet()
	bob := Wallet.NewWallet()
	alice := Wallet.NewWallet()

	satoshi.SendMoney(50, &bob.PublicKey)
	bob.SendMoney(23, &alice.PublicKey)
	alice.SendMoney(5, &bob.PublicKey)

	log.Println(Chain.GetInstance())
}