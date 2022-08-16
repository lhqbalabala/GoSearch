package handler

import (
	"RelatedSearch_RPC/models"
	"RelatedSearch_RPC/pkg/trie"
	"context"
)

type RelatedSearch struct{}

func (*RelatedSearch) GetRelatedSearch(ctx context.Context, req *models.RelatedSearchRequest, res *models.RelatedSearchResponse) (err error) {
	res.Related = trie.Tree.Search([]rune(req.Query))
	return nil
}
