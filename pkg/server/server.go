package server

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

type Client struct {
	host string
	ws   *websocket.Conn
}

type SyncRespWithClient struct {
	common.SyncResponse
	client *Client
}

const addr = "localhost:5555"

var upgrader = websocket.Upgrader{}
var wg sync.WaitGroup

func handleMessage(w http.ResponseWriter, r *http.Request,
	outputFolder string,
	fileChan chan FileDecodeMsg) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}

	client := Client{
		ws:   c,
		host: r.Host,
	}

	log.Println("Connected to host: ", client.host)

	defer c.Close()
	defer wg.Wait()

	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var request common.BaseRequest
		err = json.Unmarshal(msg, &request)
		if err != nil {
			log.Println("Invalid request:", err)
			break
		}

		switch request.RequestType {
		case string(common.Echo):
			go HandleEcho(msg, &client)
		case string(common.Sync):
			go HandleSync(msg, &client, outputFolder, fileChan)

		}
	}
}

func handlerWrapper(outputFolder string, fileChan chan FileDecodeMsg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleMessage(w, r, outputFolder, fileChan)
	}
}

func StartServer() {

	folder := flag.String("folder", "", "Folder to parse.")
	flag.Parse()

	if *folder == "" {
		log.Println("-folder required.")
		log.Println(*folder)
		syscall.Exit(1)
	}
	flag.Parse()

	fileChan := make(chan FileDecodeMsg, 100)
	go processSyncRequest(fileChan)
	http.HandleFunc("/", handlerWrapper(*folder, fileChan))
	log.Println("Starting server @", addr)
	log.Fatal((http.ListenAndServe(addr, nil)))
}
