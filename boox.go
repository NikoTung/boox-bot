package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/niko/boox-bot/user"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	endpoint        = "https://send2boox.com/api/1/%s"
	bucket          = "onyx-cloud"
	aliyun_endpoint = "https://oss-cn-shenzhen.aliyuncs.com/"
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

func (b *Boox) uid() string {
	return b.User.BooxUid
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

type saveResp struct {
	Category          []interface{} `json:"category"`
	Tags              []interface{} `json:"tags"`
	Formats           []string      `json:"formats"`
	ViewCount         int           `json:"viewCount"`
	DownloadCount     int           `json:"downloadCount"`
	CommentsCount     int           `json:"commentsCount"`
	Size              int           `json:"size"`
	SourceType        int           `json:"sourceType"`
	ChildCount        int           `json:"childCount"`
	PushNum           int           `json:"pushNum"`
	IsFolder          string        `json:"isFolder"`
	Parent            interface{}   `json:"parent"`
	Id                string        `json:"_id"`
	UserId            string        `json:"userId"`
	Name              string        `json:"name"`
	OwnerId           string        `json:"ownerId"`
	Title             string        `json:"title"`
	DistributeChannel string        `json:"distributeChannel"`
	Guid              string        `json:"guid"`
	Mac               string        `json:"mac"`
	DeviceModel       string        `json:"deviceModel"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

type save struct {
	Data saveData `json:"data"`
}

type saveData struct {
	Bucket              string      `json:"bucket"`
	Name                string      `json:"name"`
	Parent              interface{} `json:"parent"`
	ResourceDisplayName string      `json:"resourceDisplayName"`
	ResourceKey         string      `json:"resourceKey"`
	ResourceType        string      `json:"resourceType"`
	Title               string      `json:"title"`
}

func (s save) uri() string {
	return "push/saveAndPush"
}

func (s save) body() (io.Reader, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

type push struct {
	Ids []string `json:"ids"`
}

func (p push) body() (io.Reader, error) {
	return nil, nil
}

func (p push) uri() string {
	return "push/rePush/bat"
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

type meResp struct {
	AreaCode      string      `json:"area_code"`
	Avatar        interface{} `json:"avatar"`
	AvatarUrl     interface{} `json:"avatarUrl"`
	DeviceLimit   int         `json:"device_limit"`
	Email         string      `json:"email"`
	GiftCount     int         `json:"giftCount"`
	GoogleId      interface{} `json:"google_id"`
	Id            int         `json:"id"`
	LoginType     string      `json:"login_type"`
	Nickname      interface{} `json:"nickname"`
	OauthId       interface{} `json:"oauth_id"`
	Phone         interface{} `json:"phone"`
	RoleValue     int         `json:"roleValue"`
	Sex           interface{} `json:"sex"`
	StorageLimit  int64       `json:"storage_limit"`
	StorageUsed   int         `json:"storage_used"`
	Uid           string      `json:"uid"`
	VipCloud      int         `json:"vip_cloud"`
	VipCloudEnd   interface{} `json:"vip_cloud_end"`
	VipCloudStart interface{} `json:"vip_cloud_start"`
	WechatId      interface{} `json:"wechat_id"`
}

type me struct {
}

func (m me) uri() string {
	return "users/me"
}

func (m me) body() (io.Reader, error) {
	return nil, nil
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

//Send code to email
func (b *Boox) Send(email string) error {
	body := SendCode{Email: email}

	err, bd := b.post(body)
	if err != nil {
		return err
	}

	if !bd.isSuccess() {
		return errors.New(bd.Message)
	}

	return nil
}

//LoginBoox login to boox with email and code
//
// error
// string token
// string boox uid
func (b *Boox) LoginBoox(email string, code string) (error, string, string) {

	signUp := SignUp{Email: email, Code: code}

	err, r := b.post(signUp)
	if err != nil {
		return err, "", ""
	}

	err, mi := b.meInfo(me{})
	if err != nil {
		return err, "", ""
	}

	var t token
	err = json.Unmarshal(r.Data, &t)
	log.Printf("Login uid %s with token %s\n", mi.Uid, t)

	return err, t.Token, mi.Uid
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

	if !br.isSuccess() {
		log.Printf("[POST] Request to %s with data %s, failed with result %s.", param.uri(), param, br.Message)
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
	if !br.isSuccess() {
		log.Printf("[GET] Request to %s wit data %s, failed with result %s.", param.uri(), param, br.Message)
	}
	return nil, &br
}

func (b *Boox) saveAndPush(s save) (error, *saveResp) {
	err, sr := b.post(s)
	if err != nil {
		return err, nil
	}

	if !sr.isSuccess() {
		return errors.New(sr.Message), nil
	}

	var svr saveResp
	err = json.Unmarshal(sr.Data, &svr)

	return nil, &svr
}

func (b *Boox) rePush(p push) error {
	err, rp := b.post(p)
	if err != nil {
		return err
	}

	if !rp.isSuccess() {
		return errors.New(rp.Message)
	}

	return nil
}

func (b *Boox) Upload(url, name string) error {
	err, a := b.aliyunSts()
	if err != nil {
		return err
	}

	client, err := oss.New(aliyun_endpoint, a.AccessKeyId, a.AccessKeySecret, oss.SecurityToken(a.SecurityToken))
	if err != nil {
		log.Println("Aliyun oss client create error, ", err)
		return err
	}

	h, err := http.Get(url)
	if err != nil {
		log.Println("Get document error, ", err)
		return err
	}

	bk, err := client.Bucket(bucket)
	if err != nil {
		return err
	}

	key, t := resourceKey(b.uid(), name)
	err = bk.PutObject(key, h.Body)
	if err != nil {
		log.Printf("Put object ,resource key %s error %s", key, err)
		return err
	}

	s := save{saveData{Name: name, Bucket: bucket, ResourceDisplayName: name, ResourceType: t, Title: name, ResourceKey: key}}
	err, m := b.saveAndPush(s)
	if err != nil {
		log.Printf("Save and push to boox error,%s,%s", s, err)
		return err
	}

	err = b.rePush(push{Ids: []string{m.Guid}})
	if err != nil {
		return err
	}

	return nil
}

func (b *Boox) meInfo(m Requestable) (error, meResp) {
	err, mi := b.get(m)
	if err != nil {
		return err, meResp{}
	}

	if !mi.isSuccess() {
		return errors.New(mi.Message), meResp{}
	}

	var mr meResp
	err = json.Unmarshal(mi.Data, &mr)

	return err, mr
}

func resourceKey(uid, name string) (string, string) {
	f, _ := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 54)
	s := fmt.Sprintf("%x", md5.Sum([]byte(f)))
	t := ""
	if len(name) > 0 {
		if i := strings.LastIndex(name, "."); i != -1 {
			t = name[i+1:]
		}
	}

	return strings.Join([]string{uid, "push", s}, "/") + "." + t, t
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
