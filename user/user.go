package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hashicorp/golang-lru"
	"strings"
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

func Get(id int64) (*User, error) {
	e, ok := userCache.Get(id)

	if !ok {
		return nil, errors.New("user not exist")
	}

	return e.(*User), nil
}

type sign struct {
	Id        int    `json:"id"`
	LoginType string `json:"loginType"`
	Iat       int    `json:"iat"`
	Exp       int64  `json:"exp"`
}

func UpdateToken(id int64, uid, token string) error {
	e, ok := userCache.Get(id)
	if !ok {
		return errors.New("user not exist")
	}

	u := e.(*User)
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

func Add(id int64, email string) {

	u := &User{
		Id:    id,
		Email: email,
	}

	userCache.Add(id, u)
}
