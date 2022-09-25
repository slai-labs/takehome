package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log"
	"path/filepath"
	client "slai.io/takehome/pkg/client"
	"strings"
	"syscall"
)

func main() {
	folder := flag.String("folder", "", "Folder to parse.")
	ignoreDotFiles := flag.Bool("ignore-dot-files", false, "If true ignore dotfiles")
	recursive := flag.Bool("recursive", false, "If passed recursively add subfolders.")
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
				if *ignoreDotFiles && strings.HasPrefix(filepath.Base(event.Name), ".") {
					log.Printf("Ignoring: %s", event.Name)
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("Sending: '%s'", event.Name)
					uploadChan <- event.Name
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*folder)

	if *recursive {
		filepath.WalkDir(*folder, recursiveAdd(watcher))
	}

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}

func recursiveAdd(w *fsnotify.Watcher) fs.WalkDirFunc {
	return func(path string, di fs.DirEntry, err error) error {
		err = w.Add(path)

		if err != nil {
			log.Println("Can't add: ", path)
		}
		return err
	}
}
