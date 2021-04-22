package main

import (
	"app/server/config"
	"app/server/server"
	"fmt"
	"log"
)

func main() {
    fmt.Println("hoba")
    parsedConfig, err := config.ParseConfig("server/config.json")
    if err != nil {
        log.Fatalln("Error parsing config:", err.Error())
    }
    fmt.Printf("%+v\n", parsedConfig)
	server.StartServer()
}
