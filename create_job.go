package main

import (
	mycontract2 "app/contracts"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Contract address was not specified")
	}
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws/v3/98e5b995866645f8a253259b536765bb")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("1a944ffbd9a31e25053e2f29fa75ae8e70f7617f8eb908a48744fc8c79238b1a")
	if err != nil {
		log.Fatal(err)
	}

	contract, err := mycontract2.NewMycontract(common.HexToAddress(os.Args[1]), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // this one !

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(100)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	tx, err := contract.CreateJob(auth, "asd", 7777, 100, 0, 1, 2000000000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex())
	fmt.Println(tx.To().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal("Error for waiting for transaction to be mined:", err.Error())
	}
	log.Printf("%v\n", receipt)
}