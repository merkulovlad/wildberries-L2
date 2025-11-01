package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	log.SetOutput(os.Stderr)
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
	currentTime := time.Now().Add(response.ClockOffset)
	fmt.Printf("Current time: %02d:%02d:%02d\n", currentTime.Hour(), currentTime.Minute(), currentTime.Second())
}
