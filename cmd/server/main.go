package main

import (
	"log"

	"github.com/skillcoder/pow-protect/internal/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	log.Println("Bye!")
}
