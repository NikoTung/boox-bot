package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/niko/boox-bot/user"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	endpoint = "https://send2boox.com/api/1/%s"
	bucket   = "onyx-cloud"
)

type Response interface {
	isSuccess() bool
}

type Boox struct {
	User *user.User
}

func NewBoox(u *user.User) *Boox {
	return &Boox{
		User: u,
	}
}

func (b *Boox) token() string {
	if len(b.User.Token) > 0 {
		return b.User.Token
	}

	return ""
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

//Requestable interface for post
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

//AliyunConfig get aliyun oss upload configuration
type AliyunConfig struct {
}

func (ac AliyunConfig) uri() string {

	return "config/stss"
}

func (ac AliyunConfig) body() (io.Reader, error) {
	return nil, nil
}

type aliyunSts struct {
	AccessKeyId     string    `json:"AccessKeyId"`
	AccessKeySecret string    `json:"AccessKeySecret"`
	Expiration      time.Time `json:"Expiration"`
	SecurityToken   string    `json:"SecurityToken"`
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

//Send send code to email
func (b *Boox) Send(email string) error {
	body := SendCode{Email: email}

	err, _ := b.post(body)
	if err != nil {
		return err
	}

	return nil
}

//LoginBoox login to boox with email and code
func (b *Boox) LoginBoox(email string, code string) (error, string) {

	signUp := SignUp{Email: email, Code: code}

	err, r := b.post(signUp)
	if err != nil {
		return err, ""
	}
	var t token
	err = json.Unmarshal(r.Data, &t)
	log.Printf("Login with token %s\n", t)

	return err, t.Token
}

func (b *Boox) aliyunSts() (error, aliyunSts) {
	a := AliyunConfig{}
	err, sts := b.get(a)

	var t aliyunSts
	err = json.Unmarshal(sts.Data, &t)

	return err, t
}

func (bx *Boox) post(param Requestable) (error, *BooxResponse) {

	b, err := param.body()
	if err != nil {
		return err, nil
	}
	fullEndpoint := fmt.Sprintf(endpoint, param.uri())

	r, err := http.NewRequest("POST", fullEndpoint, b)

	resp, err := request(r, bx.token())

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

func (b *Boox) get(param Requestable) (error, *BooxResponse) {
	fullEndpoint := fmt.Sprintf(endpoint, param.uri())

	r, err := http.NewRequest("GET", fullEndpoint, nil)
	if err != nil {
		return err, nil
	}

	resp, err := request(r, b.token())

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

func request(r *http.Request, t string) (*http.Response, error) {

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json;charset=UTF-8")
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36")
	r.Header.Set("Sec-GPC", "1")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	r.Header.Set("Sec-Fetch-Mode", "cors")
	r.Header.Set("Sec-Fetch-Dest", "empty")
	if len(t) > 0 {
		r.Header.Set("Authorization", "Bearer "+t)
	}

	return http.DefaultClient.Do(r)
}
