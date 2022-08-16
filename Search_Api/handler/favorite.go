package handler

import (
	"Search_Api/models"
	"Search_Api/pkg/db"
	"github.com/gin-gonic/gin"
)

func QueryFavoritesAndDocs(userId uint) ([]models.Favorite, error) {
	favorites := make([]models.Favorite, 0, 10)
	result := db.Engine.Preload("Docs").Find(&favorites, "user_id = ?", userId)
	return favorites, result.Error
}
func GetAllFavoritesDocs(c *gin.Context) (likes map[uint32]*models.Docs) {
	userid := GetUserIdByToken(c)
	if userid == ErrUserId {
		return nil
	}
	likes = make(map[uint32]*models.Docs, 0)
	favorites, _ := QueryFavoritesAndDocs(uint(userid))
	for _, fav := range favorites {
		for _, doc := range fav.Docs {
			likes[uint32(doc.DocIndex)] = &models.Docs{Favid: uint32(fav.ID), Docid: uint32(doc.ID)}
		}
	}
	return likes
}
