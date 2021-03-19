package Chain

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/tomassirio/pibcoin/Block"
	transaction "github.com/tomassirio/pibcoin/Transaction"
	"os"
	"strconv"
	"sync"
)

var (
	lock = &sync.Mutex{}
	chainInstance *Chain
)

type Chain []*Block.Block

func GetInstance() *Chain {
	if chainInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if chainInstance == nil {
			fmt.Println("Creating Chain instance now.")
			//// Genesis block

			genesisPK, _ := rsa.GenerateKey(rand.Reader, 2048)
			genesisPublic := &genesisPK.PublicKey

			userPK, _ := rsa.GenerateKey(rand.Reader, 2048)
			userPublic := &userPK.PublicKey

			chainInstance = &Chain{Block.NewBlock("", transaction.NewTransaction(100, genesisPublic, userPublic))}
		} else {
			fmt.Println("Chain instance already created.")
		}
	} else {
		fmt.Println("Chain instance already created.")
	}

	return chainInstance
}

func (c Chain) GetLastBlock() *Block.Block{
	return c[len(c)-1]
}

func mine(nonce int) int{
	sol := 1
	fmt.Println("⛏️  mining...")

	for {
		h := md5.New()
		h.Write([]byte(strconv.Itoa(nonce + sol)))

		att := hex.EncodeToString(h.Sum(nil))

		if att[0:4] == "0000" {
			fmt.Printf("Solved: %d\n", sol)
			return sol
		}
		sol++
	}
}

func (c *Chain) AddBlock(t *transaction.Transaction, spk *rsa.PublicKey, sign []byte) {
	transMarsh, _ := t.ToString()

	_, valid := verify(string(sign), string(transMarsh), *spk)

	if valid {
		nb := Block.NewBlock(c.GetLastBlock().GetHash(), t)
		mine(nb.Nonce)
		*c = append(*c, nb)
	}

}

func verify(signature string, plaintext string, pubkey rsa.PublicKey) (string, bool) {
	sig, _ := base64.StdEncoding.DecodeString(signature)
	hashed := sha256.Sum256([]byte(plaintext))
	err := rsa.VerifyPKCS1v15(&pubkey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return "Error from verification:", false
	}
	return "Signature Verification Passed", true
}