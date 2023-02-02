package models

import "time"

type User struct {
	ID int
	UserName string
	Email string
	Password string
	AccressLevel int
	CreatedAt time.Time 
	UpdatedAt time.Time
}

type TokumeiPostData struct {
	ID      int
	ThreadID int
	Title 	string //name
	Text 	string //body
	CreatedAt time.Time 
	UpdatedAt time.Time
}

type TokumeiPostDataNumber struct {
	ThreadNumber int
}

type TokumeiPostComment struct {
	ID int
	ThreadID int
	Title	string
	Text	string
	CreatedAt	time.Time
	UpdatedAt	time.Time
}
