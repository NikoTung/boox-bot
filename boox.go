package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	endpoint = "https://send2boox.com/api/1/%s"
)

type Response interface {
	isSuccess() bool
}

//BooxResponse contain the raw response from boox
//  {
//      "data":
//      {
//          "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTEwMDU4LCJsb2dpblR5cGUiOiJlbWFpbCIsImlhdCI6MTY0ODIyMDc0OSwiZXhwIjoxNjYzNzcyNzQ5fQ.jjZojH9_gjT-N2BFtzTqLeykPrygfpsziIM1dDc4_mc"
//      },
//      "message": "SUCCESS",
//      "result_code": 0
//  }
type BooxResponse struct {
	Data       json.RawMessage `json:"data"`
	Message    string          `json:"message"`
	ResultCode int64           `json:"result_code"`
}

func (b *BooxResponse) isSuccess() bool {
	return b.ResultCode == 0
}

type token struct {
	Token string `json:"token"`
}

//Requestable interface for request
type Requestable interface {
	body() (io.Reader, error)
	uri() string
}

//SendCode contains the information to send login code
type SendCode struct {
	Email string `json:"mobi"`
}

func (sc SendCode) body() (io.Reader, error) {
	b, err := json.Marshal(sc)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (sc SendCode) uri() string {
	return "users/sendMobileCode"
}

//SignUp contain the information for user to login boox
type SignUp struct {
	Email string `json:"mobi"`
	Code  string `json:"code"`
}

func (s SignUp) body() (io.Reader, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (s SignUp) uri() string {
	return "users/signupByPhoneOrEmail"
}

func Send(email string) error {
	body := SendCode{Email: email}

	err, _ := request(endpoint, body)
	if err != nil {
		return err
	}

	return nil
}

func LoginBoox(email string, code string) error {

	signUp := SignUp{Email: email, Code: code}

	err, b := request(endpoint, signUp)
	if err != nil {
		return err
	}
	var t token
	err = json.Unmarshal(b.Data, &t)
	log.Printf("Login with token %s\n", t)

	//TODO save token

	return nil
}

func request(endpoint string, param Requestable) (error, *BooxResponse) {

	b, err := param.body()
	if err != nil {
		return err, nil
	}
	fullEndpoint := fmt.Sprintf(endpoint, param.uri())

	request, err := http.NewRequest("POST", fullEndpoint, b)
	if err != nil {
		return err, nil
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	request.Header.Set("Sec-GPC", "1")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	var br BooxResponse
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&br)

	if err != nil {
		return err, nil
	}
	return nil, &br
}
