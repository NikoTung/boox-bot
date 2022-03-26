package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	endpoint = "https://send2boox.com/api/1/users/sendMobileCode"
)

type sendCode struct {
	Email string `json:"mobi"`
}

func SendCode(email string) error {
	body := sendCode{Email: email}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return request(endpoint, b)
}

func request(endpoint string, param []byte) error {

	request, err := http.NewRequest("POST", endpoint, bytes.NewReader(param))
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Boox Endpoint %s response %s\n", endpoint, string(body))

	return nil
}
