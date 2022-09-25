package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"log"
	client "slai.io/takehome/pkg/client"
	"syscall"
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
	// Add in watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	uploadChan := make(chan string, 100)
	go c.Sync(uploadChan)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				log.Printf("Sending: '%s'", event.Name)
				uploadChan <- event.Name
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*folder)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}
