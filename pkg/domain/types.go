package domain

import "strconv"

type UserId int64

func (userId *UserId) String() string {
	return strconv.FormatInt(int64(*userId), 10)
}

type Email string

func (email *Email) String() string {
	return string(*email)
}

type Password string

func (password *Password) String() string {
	return string(*password)
}

type AppId int32

func (appId *AppId) String() string {
	return string(*appId)
}
