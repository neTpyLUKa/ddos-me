package payload

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("1a944ffbd9a31e25053e2f29fa75ae8e70f7617f8eb908a48744fc8c79238b1a")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(address.Hex()) // this one !

	payload, err := CreatePayload(address, privateKey)
	if err != nil {
		log.Fatal("Error creating payload: ", err.Error())
	}

	signature, err := CreateSignature(payload, privateKey)
	if err != nil {
		log.Fatal("Error creating signature")
	}

	err = Verify(payload, signature, publicKeyECDSA)
	if err != nil {
		log.Fatal("Error matching signatures: ", err.Error())
	}

	signature[3] = 'a'
	err = Verify(payload, signature, publicKeyECDSA)
	if err == nil {
		log.Fatal("Verify supposed to fail, got err == nil")
	}
}