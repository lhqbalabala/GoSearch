package handler

import (
	"Search_Api/models"
	"context"
	"github.com/gin-gonic/gin"
	httpServer "github.com/go-micro/plugins/v4/client/http"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/selector"
)

func callUserAPI(c *gin.Context, s selector.Selector) (code int32, userid int32, err error) {
	myClient := httpServer.NewClient(
		client.Selector(s),
		client.ContentType("application/json"),
	)
	// 请求对象
	request := myClient.NewRequest("User_Api", "/GetUserIdByToken"+"?mytoken=\"lgdSearch\"", models.TokenRequest{Token: c.Request.Header.Get("Authorization")})
	// 响应对象
	var resp models.TokenRes
	err = myClient.Call(context.Background(), request, &resp)
	if err != nil {
		return 500, -1, err
	}
	return resp.Code, resp.UserId, nil
}
func callRelatedSearchApi(s selector.Selector, text string) (res []string, err error) {
	myClient := httpServer.NewClient(
		client.Selector(s),
		client.ContentType("application/json"),
	)
	// 请求对象
	request := myClient.NewRequest("RelatedSearch_Client", "/RelatedSearch/:"+text+"?mytoken=\"lgdSearch\"", nil)
	var resp models.RelatedSearchResponse
	err = myClient.Call(context.Background(), request, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Related, nil
}
func callHotSearchApi(s selector.Selector, text string) (err error) {
	myClient := httpServer.NewClient(
		client.Selector(s),
		client.ContentType("application/json"),
	)
	// 请求对象
	request := myClient.NewRequest("HotSearch_Client", "/HotSearch/InsertQuery"+"?mytoken=\"lgdSearch\"", models.InsertQueryRequest{Query: text})

	var resp models.RelatedSearchResponse
	err = myClient.Call(context.Background(), request, &resp)
	if err != nil {
		return err
	}
	return nil

}
