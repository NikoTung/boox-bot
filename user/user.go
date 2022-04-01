package user

import (
	"encoding/base64"
	"encoding/json"
	"github.com/hashicorp/golang-lru"
	"strings"
	"time"
)

//TODO store user

type User struct {
	Id      int64
	Email   string
	Token   string
	BooxUid string
	Expire  int64
}

var userCache, _ = lru.New(100)

//Get or create a user
func Get(id int64) *User {
	e, ok := userCache.Get(id)

	if !ok {
		e = &User{
			Id: id,
		}
		userCache.Add(id, e)
	}

	return e.(*User)
}

type sign struct {
	Id        int    `json:"id"`
	LoginType string `json:"loginType"`
	Iat       int    `json:"iat"`
	Exp       int64  `json:"exp"`
}

func (u *User) IsLogin() bool {

	return u.Expire > 0 && time.Now().Unix() > u.Expire
}

func (u *User) UpdateToken(uid, token string) error {

	li := strings.LastIndex(token, ".")
	fi := strings.Index(token, ".")

	if li != -1 && fi != -1 && li != fi {
		m := token[fi+1 : li-1]
		d, err := base64.StdEncoding.DecodeString(m)
		if err != nil {
			return err

		}

		var s sign
		err = json.Unmarshal(d, &s)
		if err != nil {
			return err
		}
		u.Expire = s.Exp
	}

	u.Token = token
	u.BooxUid = uid

	return nil
}

func (u *User) UpdateEmail(email string) {

	u.Email = email
}
