package handler

import (
	go_micro "Search_Api/pkg/go-micro"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ErrUserId = -1

func GetUserIdByToken(c *gin.Context) (userid int32) {
	userid = -ErrUserId
	code, userid, err := callUserAPI(c, go_micro.Getsel())
	if code != http.StatusOK || err != nil {
		return ErrUserId
	}
	return userid
}
