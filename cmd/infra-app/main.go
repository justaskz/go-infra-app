package main

import (
	"log"

	"github.com/justaskz/infra-app/internal/routes"
)

func main() {
	log.Println("Starting App")
	routes.Init().Run()
}
