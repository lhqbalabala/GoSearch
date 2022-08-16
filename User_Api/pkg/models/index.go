package models

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	Text string `json:"text,omitempty"`
	Url  string `json:"url,omitempty"`
}

// StorageId leveldb中的Ids存储对象
type StorageId struct {
	Id    uint32
	Score float32
}

type KeyIndex struct {
	KeyValue uint32
	Id       uint32
}

type WordMap struct {
	Map map[uint32]float32
	Len int
}

type ResponseUrl struct {
	ThumbnailUrl string  `json:"thumbnailUrl,omitempty"`
	Url          string  `json:"url,omitempty"`
	Id           uint32  `json:"id,omitempty"`
	Text         string  `json:"text,omitempty"`
	Score        float32 `json:"score,omitempty"`
	Islike       bool    `json:"islike"`           //是否被收藏
	Docsid       uint    `json:"docsid,omitempty"` //收藏id
	Favid        uint    `json:"favid,omitempty"`  //收藏夹id
}

type RemoveIndexModel struct {
	Id uint32 `json:"id,omitempty"`
}
