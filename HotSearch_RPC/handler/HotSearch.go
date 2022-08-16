package handler

import (
	"HotSearch_RPC/models"
	"HotSearch_RPC/pkg/redis"
	"HotSearch_RPC/pkg/trie"
	"context"
	"fmt"
)

type HotSearch struct{}

func (G *HotSearch) GetHotSearchList(ctx context.Context, in *models.HotSearchRequest, res *models.HotSearchListResponse) error {
	size := in.Size
	messages, err := redis.Topk(int64(size))
	if err != nil {
		fmt.Println(err)
	}
	HotSearch := make([]*models.HotSearchModel, 0, len(messages))
	for i := 0; i < len(messages); i++ {
		HotSearch = append(HotSearch, &models.HotSearchModel{Text: messages[i].Member.(string), Num: int32(messages[i].Score)})
	}
	fmt.Println(len(HotSearch))
	res.Data = HotSearch
	return nil
}
func (G *HotSearch) InsertQuery(ctx context.Context, in *models.InsertQueryRequest, res *models.InsertQueryRespense) error {
	Query := in.Query
	trie.SendText(Query)
	return nil
}
