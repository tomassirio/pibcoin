package Block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	transaction "github.com/tomassirio/pibcoin/Transaction"
	"math"
	"math/rand"
	"time"
)

type Block struct {
	PrevHash string
	Transaction transaction.Transaction
	Date time.Time
	Nonce int
}

func NewBlock(ph string, t *transaction.Transaction) *Block {
	return &Block{PrevHash: ph, Transaction: *t, Date: time.Now(), Nonce: int(math.Round(rand.Float64() * 999999999))}
}

func (b *Block) GetHash() string{
	str, _ := json.Marshal(b)
	h := sha256.New()
	h.Write(str)
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	return sha256Hash
}
