package user

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("db"))
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

}

type User struct {
	Id      int64
	Email   string
	Token   string
	BooxUid string
	Expire  int64
}

//Get or create a user
func Get(id int64) *User {

	e := &User{
		Id: id,
	}

	rows, err := db.Query("SELECT id ,email,token,boox_uid as booxUid,expire_at as expire FROM user WHERE id = ?", id)

	if err != nil {
		return e
	}

	if rows.Next() {
		var id int64
		var email string
		var token string
		var booxUid string
		var expire int64
		if err := rows.Scan(&id, &email, &token, &booxUid, &expire); err != nil {
			return e
		}
		e = &User{
			Id:      id,
			Email:   email,
			Token:   token,
			BooxUid: booxUid,
			Expire:  expire,
		}
	} else {
		_ = add(id)
	}

	return e

}

func add(id int64) error {
	_, err := db.Exec("INSERT INTO user (id) VALUES (?)", id)
	return err
}

type sign struct {
	Id        int    `json:"id"`
	LoginType string `json:"loginType"`
	Iat       int    `json:"iat"`
	Exp       int64  `json:"exp"`
}

func (u *User) IsLogin() bool {

	return u.Expire > 0 && u.Expire > time.Now().Unix()
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

	_, err := db.Exec("UPdate user set token=?,boox_uid=? where id=? ", token, uid, u.Id)
	if err != nil {
		log.Println("update token error", err)
	}

	u.Token = token
	u.BooxUid = uid

	return nil
}

func (u *User) UpdateEmail(email string) {

	_, err := db.Exec("UPdate user set email=? where id=? ", email, u.Id)
	if err != nil {
		log.Println("update email error", err)
	}

	u.Email = email
}
