package server

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

func HandleSync(msg []byte, client *Client) error {
	log.Println("Recieved SYNC request.")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid SYNC request.")
	}
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
