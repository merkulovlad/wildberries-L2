package main

import (
	"log"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
