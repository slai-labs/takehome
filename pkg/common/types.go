package common

type RequestType string

const (
	Echo RequestType = "ECHO"
	Sync RequestType = "SYNC"
)

type BaseRequest struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type BaseResponse struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type EchoRequest struct {
	BaseRequest
	Value string
}

type EchoResponse struct {
	BaseResponse
	Value string
}

type SyncRequest struct {
	BaseRequest
	EncodedFile string
	FileName    string
}

type SyncResponse struct {
	BaseResponse
	Success  bool
	FileName string
	ErrorMsg string
}
