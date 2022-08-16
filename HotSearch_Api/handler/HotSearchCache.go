package handler

import "HotSearch_Api/models"

var HotSearchCache *models.HotSearchListResponse

func GetHotSearchCache() *models.HotSearchListResponse {
	if HotSearchCache == nil {
		HotSearchCache = &models.HotSearchListResponse{}
	}
	return HotSearchCache
}
func SetHotSearchCache(HotSearch *models.HotSearchListResponse) {
	HotSearchCache = HotSearch
}
