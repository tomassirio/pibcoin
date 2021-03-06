package Wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/tomassirio/pibcoin/Chain"
	transaction "github.com/tomassirio/pibcoin/Transaction"
	"log"
	"math/big"
	"os"
	"time"
)

type Wallet struct {
	PublicKey rsa.PublicKey
	PrivateKey rsa.PrivateKey
}

func NewWallet() *Wallet {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		panic("Cannot generate RSA key\n")
		os.Exit(1)
	}

	pubKey := &privKey.PublicKey

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ORGANIZATION_NAME"},
			Country:       []string{"COUNTRY_CODE"},
			Province:      []string{"PROVINCE"},
			Locality:      []string{"CITY"},
			StreetAddress: []string{"ADDRESS"},
			PostalCode:    []string{"POSTAL_CODE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pubKey, privKey)

	if err != nil {
		log.Println("Create PEM Certificate Failed", err)
		os.Exit(1)
	}

	// Public key
	certOut, err := os.Create("ca.crt")
	pem.Encode(certOut, &pem.Block{Type: "SPKI", Bytes: ca_b})
	certOut.Close()
	log.Print("Written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "PKCS8", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	keyOut.Close()
	log.Print("Written key.pem\n")

	return &Wallet{PublicKey: *pubKey, PrivateKey: *privKey}
}

func (w *Wallet) SendMoney(amount int, payeePublicKey *rsa.PublicKey) {
	trans := transaction.NewTransaction(amount, &w.PublicKey, payeePublicKey)
	marsh, _ := trans.ToString()
	h := sha256.New()
	h.Write(marsh)

	bodyHash := sign(string(marsh), w.PrivateKey)

	Chain.GetInstance().AddBlock(trans, &w.PublicKey, bodyHash)
}

func sign(plaintext string, privKey rsa.PrivateKey)  []byte {
	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader
	hashed := sha256.Sum256([]byte(plaintext))
	signature, err := rsa.SignPKCS1v15(rng, &privKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return []byte("Error from signing")
	}
	return []byte(base64.StdEncoding.EncodeToString(signature))
}