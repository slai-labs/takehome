package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

func decodeAndSave(fileData string, filePath string) error {

	decoded, err := base64.StdEncoding.DecodeString(fileData)

	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, decoded, 0644)

	return nil
}

func processSyncRequest(fileChan <-chan common.SyncRequest) {
	for request := range fileChan {
		err := decodeAndSave(request.EncodedFile, request.FileName)

		if err != nil {
			log.Println("Can't process request.s")
		}
	}
}

func HandleSync(msg []byte, client *Client, outputFolder string, fileChan chan common.SyncRequest) error {
	log.Println("Recieved SYNC request.")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid SYNC request.")
	}

	request.FileName = filepath.Join(outputFolder, request.FileName)
	fileChan <- request

	response := common.SyncResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Success: true,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil

}

func HandleEcho(msg []byte, client *Client) error {
	log.Println("Received ECHO request.")

	var request common.EchoRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	response := &common.EchoResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: request.Value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}
