package handler

import (
	"Search_RPC/models"
	"Search_RPC/pkg"
	"Search_RPC/pkg/utils/colf/doc"
	"context"
)

type Search struct {
	models.UnimplementedSearchServer
}

func (*Search) GetSearchResult(ctx context.Context, req *models.SearchRequest) (res *models.SearchResult, err error) {
	res = pkg.SearchEngine.MultiSearch(req)
	return res, nil
}
func (*Search) GetSearchPictureResult(ctx context.Context, req *models.SearchRequest) (res *models.SearchPictureResult, err error) {
	res = pkg.SearchEngine.MultiSearchPicture(req)
	return res, nil
}
func (*Search) GetDocById(ctx context.Context, req *models.DocRequest) (res *models.DocResult, err error) {
	buf := pkg.SearchEngine.GetDocById(req.GetId())
	storageDoc := new(doc.StorageIndexDoc)
	storageDoc.UnmarshalBinary(buf)
	res = &models.DocResult{}
	res.Text = storageDoc.Text
	res.Url = storageDoc.Url
	return res, nil
}
