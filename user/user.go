package user

import "errors"

type User struct {
	Id    int64
	Email string
	Token string
}

var users = make(map[int64]*User)

func Get(id int64) *User {
	u := &User{Id: id, Email: "niko.tung@protonmail.com", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTEwMDU4LCJsb2dpblR5cGUiOiJlbWFpbCIsImlhdCI6MTY0ODIyMDc0OSwiZXhwIjoxNjYzNzcyNzQ5fQ.jjZojH9_gjT-N2BFtzTqLeykPrygfpsziIM1dDc4_mc"}
	users[id] = u

	return u
}

func UpdateToken(id int64, token string) error {
	u := users[id]
	if u == nil {
		return errors.New("user not exist")
	}
	users[id].Token = token

	return nil
}

func Add(id int64, email string) {
	u := users[id]
	if u != nil {
		return
	}

	users[id] = &User{
		Id:    id,
		Email: email,
	}
}
