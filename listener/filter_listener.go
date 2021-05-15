package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Contract address was not specified")
	}
	conn, err := ethclient.Dial("wss://ropsten.infura.io/ws/v3/98e5b995866645f8a253259b536765bb"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err.Error())
	}
	log.Println("Client created successfully")

	Topics := make([][]common.Hash, 1)
	Topics[0] = make([]common.Hash, 1)
	Topics[0][0] = common.BytesToHash([]byte("0x41ac6f01834131cd6fc6ef5bc2a06a19178750f8894f01ab8db6a51251a4ac72")) // topic
	myContractAddress := common.HexToAddress(os.Args[1])

	query := ethereum.FilterQuery{
		FromBlock: nil,
		ToBlock:   nil,
		//Topics:    Topics,
		Addresses: []common.Address{myContractAddress},
	}

	logs := make(chan types.Log, 2)
	s, err := conn.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalln("Error subscribing to filter logs: ", err.Error())
	}
	errChan := s.Err()
	for {
		select {
		case err := <-errChan:
			log.Fatal("Logs subscription error", err.Error())
		case l := <-logs:
			log.Printf("Hoba %v\n", l)
		}
	}
}
