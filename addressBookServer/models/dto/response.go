package dto

import (
	"encoding/json"
	"httpserver/pkg"
	"log"
)

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}


func (resp *Response) GetJson() (byteResp []byte, err error) {
	myErr := pkg.NewMyError("package pkg: func GetJson()")
	if resp.Data == nil {
		resp.Data = json.RawMessage(`{}`)
	}
	byteResp, err = json.Marshal(resp)
	if err != nil {
		e := myErr.Wrap(err, "")
		log.Println(e.Error())
		return byteResp, e
	}
	return byteResp, nil
}
