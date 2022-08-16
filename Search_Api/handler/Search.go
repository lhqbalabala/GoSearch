package handler

import (
	"Search_Api/models"
	go_micro "Search_Api/pkg/go-micro"
)

func InsertQuery(query string) {
	go callHotSearchApi(go_micro.Getsel(), query)
}
func AddRelatedSearch(query string, rsp *models.SearchResult) {
	rsp.Related, _ = callRelatedSearchApi(go_micro.Getsel(), query)
}
