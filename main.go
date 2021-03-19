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

	log.Println(Chain.GetInstance().GetLastBlock().Transaction)

	satoshi.SendMoney(50, &bob.PublicKey)
	log.Println(Chain.GetInstance().GetLastBlock().Transaction)

	bob.SendMoney(23, &alice.PublicKey)
	log.Println(Chain.GetInstance().GetLastBlock().Transaction)

	alice.SendMoney(5, &bob.PublicKey)
	log.Println(Chain.GetInstance().GetLastBlock().Transaction)

}