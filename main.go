package main

import (
	"app/client/client"
	"app/client/client_config"
	"app/server/server"
	"app/server/server_config"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run main.go client/server")
	}
	if os.Args[1] == "server" {
		parsedConfig, err := server_config.ParseConfig("server/config.json")
		if err != nil {
			log.Fatalln("Error parsing server_config:", err.Error())
		}
		fmt.Printf("%+v\n", parsedConfig)
		server.StartServer(parsedConfig)
	} else if os.Args[1] == "client" {
		parsedConfig, err := client_config.ParseConfig("client/config.json")
		if err != nil {
			log.Fatalln("Error parsing client_config: ", err.Error())
		}

		cl, err := client.CreateClient(parsedConfig)
		ch := make(chan struct{}, 10)
		for i := 1; i <= 10; i++ {
			go func(i int) {
				if err != nil {
					log.Fatalln("Error creating client:", err.Error())
				}
				err = cl.Create(i)
				if err != nil {
					log.Fatalln("Error creating payload:", err.Error())
				}
				err = cl.Increment(i)
				if err != nil {
					log.Fatalln("Error incrementing payload:", err.Error())
				}
				err = cl.Increment(i)
				if err != nil {
					log.Fatalln("Error incrementing payload:", err.Error())
				}
				err = cl.Increment(i)
				if err != nil {
					log.Fatalln("Error incrementing payload:", err.Error())
				}
				ch <- struct{}{}
			}(i)
		}
		for i := 0; i < 10; i++ {
			<-ch
		}
		err = cl.Merge()
		if err != nil {
			log.Fatalln("Error merging payloads:", err.Error())
		}
		log.Println(cl.Decoded[0], cl.Signatures[0])
	} else {
		log.Fatalln("Unexpected start mode")
	}
}
