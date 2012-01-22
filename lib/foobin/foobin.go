package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func uuid() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	s := fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	return s, nil
}

type JsonRpcReq struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Result struct {
	Success bool   `json:"success"`
	Code    int64  `json:"code"`
	Note    string `json:"note"`
}

type Person struct {
	Id    int64
	Name  string
	Email string "pattern: \\S+@\\S+.\\S+"
	Title string
}

type ResultStructResponse struct {
	Id     string   `json:"id"`
	Error  RPCError `json:"error"`
	Result Result   `json:"result"`
}

type RPCError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("RPCError code: %d message: %s", e.Code, e.Message)
}

type SampleService interface {
	Create(p Person) (Result, *RPCError)
	Add(a int64, b int64) (int64, *RPCError)
	StoreName(name string) *RPCError
	Say_Hi() (string, *RPCError)
}

type SampleServiceClient struct {
	Url string
}

func (c *SampleServiceClient) call(method string, params interface{}) ([]byte, *RPCError) {

	id, err := uuid()
	if err != nil {
		return nil, &RPCError{-32000, err.Error()}
	}

	req := JsonRpcReq{"2.0", id, method, params}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, &RPCError{-32000, err.Error()}
	}

	buf := bytes.NewBuffer(body)
	resp, err := http.Post(c.Url, "text/plain", buf)
	if err != nil {
		return nil, &RPCError{-32000, err.Error()}
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &RPCError{-32000, err.Error()}
	}

	return body, nil
}

func (c *SampleServiceClient) Add(a int64, b int64) (int64, *RPCError) {
	resp, err := c.call("SampleService_Add", []int64{a, b})
	if err != nil {
		return 0, err
	}

	var f interface{}
	err2 := json.Unmarshal(resp, &f)
	if err2 != nil {
		return 0, &RPCError{-32000, err2.Error()}
	}

	m := f.(map[string]interface{})
	if rpcerr, ok := m["error"]; ok {
		errmap := rpcerr.(map[string]interface{})
		return 0, &RPCError{errmap["code"].(int32), errmap["message"].(string)}
	}

	if result, ok := m["result"]; ok {
		if num, ok := result.(float64); ok {
			return int64(num), nil
		}
	}
	return 0, &RPCError{-32000, fmt.Sprintf("Invalid response: %v", resp)}
}

func (c *SampleServiceClient) Create(p Person) (*Result, *RPCError) {
	resp, err := c.call("SampleService_Create", p)
	if err != nil {
		return nil, err
	}

	var f ResultStructResponse
	err2 := json.Unmarshal(resp, &f)
	if err2 != nil {
		return nil, &RPCError{-32000, err2.Error()}
	}

	if f.Error.Code != 0 || f.Error.Message != "" {
		return nil, &f.Error
	}
	return &f.Response, nil
}

func main() {

	client := SampleServiceClient{"http://localhost:9009"}

	start := time.Now().UnixNano()
	for i := 0; i < 10000; i++ {
		sum, err := client.Add(int64(i), 6)
		if err != nil {
			fmt.Printf("got err: %s\n", err)
		} else if sum != (int64(i) + 6) {
			fmt.Printf("got wrong answer: %d\n", sum)
		}
		//fmt.Printf("Add reply: %d\n", ret)
	}
	elapsed := time.Now().UnixNano() - start
	fmt.Printf("Elapsed: %d\n", elapsed/1e6)

}
