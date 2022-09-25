package main

import (
	"flag"
	"log"
	"syscall"
	"time"

	client "slai.io/takehome/pkg/client"
)

func main() {
	folder := flag.String("folder", "", "Folder to parse.")
	flag.Parse()

	if *folder == "" {
		log.Println("-folder required.")
		log.Println(*folder)
		syscall.Exit(1)
	}

	log.Printf("Starting client. Monitoring folder: %q", *folder)

	c, err := client.NewClient("./")
	if err != nil {
		log.Fatal(err)
	}

	someMessage := "hello there"
	for {

		log.Printf("Sending: '%s'", someMessage)

		value, err := c.Echo(someMessage)
		if err != nil {
			log.Fatal("Unable to send request.")
		}

		log.Printf("Received: '%s'", value)

		time.Sleep(time.Second)
	}

}
