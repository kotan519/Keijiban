package repository

import (
	//"time"

	"github.com/kotan519/keijiban/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertData(res models.TokumeiPostData) error
	GetThreadList() ([]models.TokumeiPostData, error)
	GetThreadData(num int) ([]models.TokumeiPostData, error)
	InsertCommentData(as models.TokumeiPostData) error
	GetCommentData(num int) ([]models.TokumeiPostData, error)
	Authenticate(email, testpassword string) (int, string, error)
	InsertUserData(as models.User) error
}