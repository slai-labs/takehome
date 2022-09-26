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

type FileDecodeMsg struct {
	request      common.SyncRequest
	responseChan chan SyncRespWithClient
	client       *Client
}

func decodeAndSave(fileData string, filePath string) error {

	decoded, err := base64.StdEncoding.DecodeString(fileData)

	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, decoded, 0644)

	return nil
}

func processSyncRequest(fileChan <-chan FileDecodeMsg) {
	for msg := range fileChan {

		err := os.MkdirAll(filepath.Dir(msg.request.FileName), os.ModePerm)
		if err != nil {
			log.Fatal("Can't create: ", filepath.Dir(msg.request.FileName))
		}

		errMsg := ""
		success := true

		err = decodeAndSave(msg.request.EncodedFile, msg.request.FileName)
		if err != nil {
			log.Println("Can't process request.")
			errMsg = err.Error()
			success = false
		}

		response := &common.SyncResponse{
			BaseResponse: common.BaseResponse{
				RequestId:   msg.request.RequestId,
				RequestType: msg.request.RequestType,
			},
			Success:  success,
			ErrorMsg: errMsg,
			FileName: msg.request.FileName,
		}

		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Printf("WARN: Can't marshall %s", msg.request.RequestId)

		}

		client := *msg.client
		err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
		if err != nil {
			log.Fatal("Can't send message. Shutting down. IRL this would probably be more graceful.")
		}

	}

}

func HandleSync(msg []byte, client *Client,
	outputFolder string,
	fileChan chan FileDecodeMsg) error {

	log.Println("Received SYNC request.")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Println("Malformed sync request. If I had more type I'd send an error message back to the client.")
	}

	request.FileName = filepath.Join(outputFolder, request.FileName)
	log.Println(request.FileName)
	fileChan <- FileDecodeMsg{
		request: request,
		client:  client,
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
