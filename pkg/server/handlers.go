package server

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"path/filepath"
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

type FileDecodeMsg struct {
	request      common.SyncRequest
	responseChan chan SyncRespWithClient
	client       *Client
}

func processSyncRequest(fileChan <-chan FileDecodeMsg) {
	for msg := range fileChan {
		err := decodeAndSave(msg.request.EncodedFile, msg.request.FileName)

		if err != nil {
			log.Println("Can't process request.")
		}

		response := SyncRespWithClient{
			common.SyncResponse{
				BaseResponse: common.BaseResponse{
					RequestId:   msg.request.RequestId,
					RequestType: msg.request.RequestType,
				},
				Success: true,
			},
			msg.client,
		}

		msg.responseChan <- response

	}

}

func returnSyncStatus(syncChan <-chan SyncRespWithClient) {
	for response := range syncChan {
		log.Println("Responding")
		responsePayload, err := json.Marshal(response.SyncResponse)
		if err != nil {
			log.Println("Can't marshall")

		}
		client := *response.client

		err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
		if err != nil {
			log.Println("Can't send")
		}

	}

}

func HandleSync(msg []byte, client *Client,
	outputFolder string,
	fileChan chan FileDecodeMsg,
	responseChan chan SyncRespWithClient) error {

	log.Println("Received SYNC request.")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid SYNC request.")
	}

	request.FileName = filepath.Join(outputFolder, request.FileName)
	fileChan <- FileDecodeMsg{
		request:      request,
		responseChan: responseChan,
		client:       client,
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
