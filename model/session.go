package model

import (
	"time"
)


// Session with built-in json serialization
type Session struct {
	SessionID	string 	 	`json:"sessionId"`
	User		*Account 	`json:"-"`	
	UserType	string		`json:"userType"`
	CreateTime	time.Time	`json:"createTime"`
	UpdateTime	time.Time	`json:"updateTime"`
	Expires		time.Time	`json:"expires"`
}
