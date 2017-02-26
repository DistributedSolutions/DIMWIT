package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Helper for jsonrpc calls
type JSONRPCRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	ID      uint32        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewEmptyJSONRPCRequest() *JSONRPCRequest {
	j := new(JSONRPCRequest)
	j.JsonRpc = "2.0"
	j.ID = 0
	j.Method = ""
	j.Params = nil
	return j
}

func NewJSONRPCRequest(method string, params interface{}, id uint32) *JSONRPCRequest {
	paramSet := make([]interface{}, 1)
	paramSet[0] = params
	return NewJSONRPCRequestArray(method, paramSet, id)
}

func NewJSONRPCRequestArray(method string, params []interface{}, id uint32) *JSONRPCRequest {
	j := new(JSONRPCRequest)
	j.JsonRpc = "2.0"
	j.ID = id
	j.Method = method
	j.Params = params
	return j
}

func (j *JSONRPCRequest) CustomMarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JSONRPCRequest) CustomUnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, j)
}

// POST Request takes the return object and unmarshals the reponse data into the obj if there is no error
func (j *JSONRPCRequest) POSTRequest(url string, obj interface{}) (interface{}, *JSONError, error) {
	data, err := j.CustomMarshalJSON()
	if err != nil {
		return nil, nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	jsonResp := new(JSONRPCReponse)
	err = jsonResp.CustomUnmarshalJSON(respData)
	if err != nil {
		return nil, nil, err
	}

	if jsonResp.Error != nil {
		return nil, jsonResp.Error, nil
	}

	result, err := jsonResp.Result.MarshalJSON()
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(result, obj)
	if err != nil {
		return nil, nil, err
	}

	return obj, nil, nil
}

type JSONRPCReponse struct {
	Result *json.RawMessage `json:"result"`
	Error  *JSONError       `json:"error,omitempty"`
	Id     uint32           `json:"id"`
}

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j *JSONRPCReponse) CustomMarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JSONRPCReponse) CustomUnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, j)
}
