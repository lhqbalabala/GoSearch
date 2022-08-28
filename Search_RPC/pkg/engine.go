package pkg

import (
	"Search_RPC/models"
	"Search_RPC/pkg/bitmap"
	"Search_RPC/pkg/db/badgerStorage"
	"Search_RPC/pkg/db/boltStorage"
	"Search_RPC/pkg/kafka"
	"Search_RPC/pkg/logger"
	"Search_RPC/pkg/pagination"
	"Search_RPC/pkg/redis"
	"Search_RPC/pkg/trie"
	"Search_RPC/pkg/utils"
	"Search_RPC/pkg/utils/colf/doc"
	"Search_RPC/pkg/utils/colf/keyIds"
	"Search_RPC/pkg/utils/colf/results"
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dgraph-io/badger/v3"
	"github.com/nfnt/resize"
	"github.com/wangbin/jiebago"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	//Shard is the default engine shard
	Shard = 100

	// BadgerShard badgerDB的分库数
	BadgerShard = 10

	// BoltBucketSize boltDB中桶的数量
	BoltBucketSize = 100

	// IndexPath 数据文件地址
	IndexPath = "./pkg/data"

	//// GitHubVersion github版本为 50000 条数据
	//GitHubVersion = 50000

	// CloudVersion 完整版本为 1亿条数据
	CloudVersion = 100000000
)

type Option struct {
	PictureName    string
	BoltKeyIdsName string
	BadgerDocName  string
	DictionaryName string
	DocIndexName   string
}
type Engine struct {
	IndexPath string

	Option *Option

	//关键字和Id映射
	KeyIdsStorage   *boltStorage.BoltdbStorage
	boltBucketNames [][]byte

	//文档仓
	DocStorages []*badgerStorage.BadgerStorage

	//慢结果仓
	SlowResultStorages *badgerStorage.BadgerStorage

	//锁
	sync.Mutex
	//等待
	sync.WaitGroup

	//文件分片
	Shard int

	//添加索引的通道
	AddDocumentWorkerChan chan models.IndexDoc

	// index 是data的id索引，可以理解为 doc 的个数
	index uint32

	//是否调试模式
	IsDebug bool

	//分词过滤
	Bitmap *bitmap.Bitmap

	// 版本
	Version int

	//kafka
	SearchResultMQ  kafka.AbsKafkaProducer
	SearchRequestMQ kafka.AbsKafkaProducer

	//kafka通道
	SearchResultMQChan  chan string
	SearchRequestMQChan chan string
}

var seg jiebago.Segmenter
var SearchEngine *Engine
var wordFilterMap map[string]interface{}

func filterMapInit() {
	wordFilterMap = make(map[string]interface{})
	wordFilterMap["了"] = nil
	wordFilterMap["的"] = nil
	wordFilterMap["么"] = nil
	wordFilterMap["呢"] = nil
	wordFilterMap["和"] = nil
	wordFilterMap["与"] = nil
	wordFilterMap["于"] = nil
	wordFilterMap["吗"] = nil
	wordFilterMap["吧"] = nil
	wordFilterMap["呀"] = nil
	wordFilterMap["啊"] = nil
	wordFilterMap["哎"] = nil
	wordFilterMap["是"] = nil
	wordFilterMap["人"] = nil
	wordFilterMap["名"] = nil
	wordFilterMap["在"] = nil
	wordFilterMap["不"] = nil
	wordFilterMap["被"] = nil
	wordFilterMap["有"] = nil
	wordFilterMap["无"] = nil
	wordFilterMap["都"] = nil
	wordFilterMap["也"] = nil
	wordFilterMap["这"] = nil
	wordFilterMap["是"] = nil
	wordFilterMap["好"] = nil
	wordFilterMap["【"] = nil
	wordFilterMap["】"] = nil
	wordFilterMap["《"] = nil
	wordFilterMap["》"] = nil
	wordFilterMap["，"] = nil
	wordFilterMap["。"] = nil
	wordFilterMap["？"] = nil
	wordFilterMap["！"] = nil
	wordFilterMap["、"] = nil
	wordFilterMap["；"] = nil
	wordFilterMap["："] = nil
	wordFilterMap["（"] = nil
	wordFilterMap["）"] = nil
	wordFilterMap["什么"] = nil
}

func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Shard))
}

func Set(e *Engine) {
	SearchEngine = e
}
func DefaultIndexPath() string {
	return IndexPath
}
func (e *Engine) getFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}
func (e *Engine) GetOptions() *Option {
	return &Option{
		PictureName:    "picture/pic",
		DictionaryName: "dictionary.txt",
		DocIndexName:   "dataIndex.txt",
		BoltKeyIdsName: "bolt_keyIds.db",
		BadgerDocName:  "badger_doc",
	}
}
func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()
	//线程数=cpu数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//保持和gin一致
	e.IsDebug = os.Getenv("GIN_MODE") != "release"
	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	bitmap.InitBitmap("pkg/data/keys.txt")
	e.Bitmap = bitmap.GetBitmap()
	//初始化kafka
	e.SearchResultMQ = kafka.NewProducer("120.48.100.204:9092,43.142.141.48:9092,121.196.207.80:9092", "SearchResult", time.Second*5, kafka.Async)
	e.SearchRequestMQ = kafka.NewProducer("120.48.100.204:9092,43.142.141.48:9092,121.196.207.80:9092", "SearchRequest", time.Second*5, kafka.Async)

	err := seg.LoadDictionary(e.getFilePath(e.Option.DictionaryName))
	if err != nil {
		panic("dictionary not find")
	}

	if e.Shard == 0 {
		e.Shard = Shard
	}

	utils.Read(&e.index, e.getFilePath(e.Option.DocIndexName))

	if e.index == 0 {
		e.index = CloudVersion
	}

	//初始化chan
	e.AddDocumentWorkerChan = make(chan models.IndexDoc, 1000)
	e.SearchResultMQChan = make(chan string, 2000)
	e.SearchRequestMQChan = make(chan string, 2000)

	filterMapInit() // 初始化分词过滤

	option := badger.DefaultOptions(e.getFilePath("SlowResult"))
	option.Logger = nil
	e.SlowResultStorages = badgerStorage.Open(option)

	e.KeyIdsStorage, err = boltStorage.Open(e.getFilePath(e.Option.BoltKeyIdsName))
	for i := 0; i < BadgerShard; i++ {
		option := badger.DefaultOptions(e.getFilePath(fmt.Sprintf("%s_%d.db", e.Option.BadgerDocName, i)))
		option.Logger = nil
		badgerdb := badgerStorage.Open(option)
		e.DocStorages = append(e.DocStorages, badgerdb)
	}

	for i := 0; i < BoltBucketSize; i++ {
		bucketName := utils.Uint32ToBytes(uint32(i))
		e.boltBucketNames = append(e.boltBucketNames, bucketName)
		err := e.KeyIdsStorage.CreateBucketIfNotExist(bucketName)
		if err != nil {
			return
		}
	}

	//初始化文件存储
	go e.DocumentWorkerExec()
	//初始化chan通信
	go e.SearchResultMQExec()
	go e.SearchRequestMQExec()

}

//SearchRequestMQ队列
func (e *Engine) SearchRequestMQExec() {
	for {
		mq := <-e.SearchRequestMQChan
		fmt.Println(mq)
		e.SearchRequestMQ.Send([]byte(mq))
	}
}

//SearchResMQ队列
func (e *Engine) SearchResultMQExec() {
	for {
		mq := <-e.SearchResultMQChan
		e.SearchResultMQ.Send([]byte(mq))
	}
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec() {
	for {
		docs := <-e.AddDocumentWorkerChan
		e.AddDocument(&docs)
	}
}
func (e *Engine) WordCutFilter(text string) ([]uint32, map[uint32]float32) {
	//不区分大小写
	text = strings.ToLower(text)

	// wordMap is to save word, keyMap is to save hash(word)
	var keyMap = make(map[uint32]float32)
	var wordMap = make(map[uint32]interface{})
	words := make([]uint32, 0)

	resultChan := seg.CutForSearch(text, true)
	pre := true
	for {
		w, ok := <-resultChan
		if !ok {
			break
		}

		switch len(w) { //  过滤分词
		case 1:
			{
				if (w[0] > 47 && w[0] < 58) || (w[0] > 64 && w[0] < 91) || (w[0] > 96 && w[0] < 123) {
					// 保留数字和字母
				} else if w[0] <= 32 || w[0] == 127 || !pre { // 不要空格和分隔符, 不常用字符和符号只保留首位
					pre = true
					continue
				}
			}
		case 3:
			{
				_, ok := wordFilterMap[w]
				if ok && !pre {
					continue
				}
			}
		}

		value := utils.StringToInt(w)
		words = append(words, value)

		wordMap[value] = nil
	}
	lenWords := float32(len(words))
	for index, val := range words { // 越前面的比重越大
		// math.log10(10 + index) 是为了平衡 f(index=10,len=10) 和 f(index=5000,len=10000)的情况
		// 这种情况下 f(10, 10) > f(5000, 10000), etc...
		// 目的是使得 f(10, 10000) > f(10, 10) > f(5000, 10000)
		keyMap[val] += float32(2*len(words)-index+1) / (float32(math.Log10(float64(10+index))) * lenWords)
	}

	var keysSlice = make([]uint32, len(keyMap))

	index := 0
	for word := range wordMap {
		keysSlice[index] = word
		index += 1
	}

	return keysSlice, keyMap
}

// WordCut 分词，只取长度大于2的词 | 是假的，会取到长度为1的词, 需要取到长度为1的词
// 这里输入的请求长度有限制,最多30,那么分词不会很多,同时考虑对于关键词过滤是否也做限制 10 个
func (e *Engine) WordCut(text string) []string {
	//不区分大小写
	text = strings.ToLower(text)

	// wordMap is to save word, keyMap is to save hash(word)
	var wordMap = make(map[string]int)

	resultChan := seg.CutForSearch(text, true)
	pre := true
	for {
		w, ok := <-resultChan
		if !ok {
			break
		}
		switch len(w) { // 过滤分词
		case 1:
			{
				if (w[0] > 47 && w[0] < 58) || (w[0] > 64 && w[0] < 91) || (w[0] > 96 && w[0] < 123) {
					// 保留数字和字母
				} else if w[0] <= 32 || w[0] == 127 || !pre { // 不要空格和分隔符, 不常用字符和符号只保留首位
					pre = true // 对于这些分隔符可以认为是一句新的话的开头
					continue
				}
			}
		case 3:
			{
				_, ok := wordFilterMap[w]
				if ok && !pre {
					continue
				}
			}
		}
		wordMap[w]++
		pre = false
	}

	var wordsSlice = make([]string, len(wordMap))

	index := 0
	for word := range wordMap {
		if e.Bitmap.Has(int(utils.StringToInt(word))) {
			wordsSlice[index] = word
			index += 1
		}
	}

	return wordsSlice
}

// AddDocument 添加文档
func (e *Engine) AddDocument(index *models.IndexDoc) {
	//等待初始化完成
	e.Wait()

	text := index.Text

	// wordLen 是带重复数据的长度
	keys, wordMap := e.WordCutFilter(text)

	for _, key := range keys {
		data := keyIds.StorageId{
			Id:    index.Id,
			Score: wordMap[key],
		}

		keyBuf := utils.Uint32ToBytes(key)
		bucketName := e.boltBucketNames[key%BoltBucketSize]
		buf, found := e.KeyIdsStorage.Get(keyBuf, bucketName)

		ids := new(keyIds.StorageIds)
		if found {
			ids.UnmarshalBinary(buf)
			ids.StorageIds = append(ids.StorageIds, &data)
		} else {
			ids.StorageIds = append(ids.StorageIds, &data)
		}

		bufs, _ := ids.MarshalBinary()
		e.KeyIdsStorage.Set(keyBuf, bufs, bucketName)
	}
	e.addDoc(index)
}
func (e *Engine) addDoc(index *models.IndexDoc) {
	k := utils.Uint32ToBytes(index.Id)

	value := doc.StorageIndexDoc{Text: index.Text, Url: index.Url}
	buf, err := value.MarshalBinary()
	if err != nil {
		logger.Logger.Infoln("addDoc id = %d, err = %s", index.Id, err)
	}
	err = e.DocStorages[index.Id%BadgerShard].Set(k, buf)
	if err != nil {
		logger.Logger.Infoln("setDoc id = %d, err = %s", index.Id, err)
	}
}

//SimpleSearch key is one of keys, keys -> 当前查询语句的所有分词
func (e *Engine) SimpleSearch(word string, key uint32, call func(ranks []*models.SliceItem)) {
	_time := time.Now()

	s := e.KeyIdsStorage

	kv := utils.Uint32ToBytes(key)

	data, find := s.Get(kv, e.boltBucketNames[key%BoltBucketSize]) // key.ids

	if find {
		array := new(keyIds.StorageIds)
		array.UnmarshalBinary(data)

		results := make([]*models.SliceItem, len(array.StorageIds))
		if e.IsDebug {
			logger.Logger.Infoln("读数据时间: ", time.Since(_time), "--- word: ", word, "--- key: ", key, "--- Ids长度:", len(array.StorageIds))
		}

		for index, id := range array.StorageIds { // 遍历 ids
			rank := &models.SliceItem{}
			rank.Id = id.Id
			rank.Score = float32(math.Log10(float64(e.index)/float64(len(array.StorageIds)+1))) * id.Score
			results[index] = rank
		}
		if e.IsDebug {
			logger.Logger.Infoln("遍历ids时间: ", time.Since(_time), "--- word: ", word, "--- key: ", key, "--- Ids长度:", len(array.StorageIds))
		}
		call(results)
	} else {
		call(nil)
	}
}

func (e *Engine) MultiSearch(request *models.SearchRequest) *models.SearchResult {
	//等待搜索初始化完成
	if e.IsDebug {
		logger.Logger.Infoln("Search start ------------------------")
	}
	e.Wait()
	//搜索开始时间
	StartTime := time.Now()

	ks := request.Query
	//发送查询到mq
	e.SearchRequestMQChan <- ks
	//分页
	pager := new(pagination.Pagination)
	//结果集
	var resultIds []pagination.DocItem
	//结果集长度
	var totlen int
	//返回结果
	var result *models.SearchResult
	//分词
	var words []string
	//是否在redis中
	isfound := redis.Isfound(ks)
	if isfound {
		logger.Logger.Infoln("该结果集从boltdb取出")
		//过滤
		trie.BuildTree(request.FilterWord)

		// 设置fail指针
		trie.SetNodeFailPoint()
		// 处理分页
		GetAndSetDefault(request)
		//分词搜索
		words = e.WordCut(request.Query)
		if e.IsDebug {
			logger.Logger.Infoln("分词数: ", len(words))
		}
		//读取文档
		result = &models.SearchResult{
			Total:     0,
			Time:      0,
			Page:      request.Page,
			Limit:     request.Limit,
			Words:     words,
			Documents: nil,
		}
		bufs, _ := e.SlowResultStorages.Get(utils.Encoder(ks))
		slowresult := results.Result{}
		slowresult.UnmarshalBinary(bufs)
		resultslice := slowresult.Sliceresult
		slowresultlen := int(resultslice[len(resultslice)-1].Id)
		for index := 0; index < len(resultslice)-1; index++ {
			resultIds = append(resultIds, pagination.DocItem{resultslice[index].Id, resultslice[index].Score, 1})
		}
		result.Total = uint32(slowresultlen)

	} else {
		logger.Logger.Infoln("该结果集为正常搜索")
		//分词搜索
		words = e.WordCut(request.Query)
		if e.IsDebug {
			logger.Logger.Infoln("分词数: ", len(words))
		}

		var lock sync.Mutex
		var wg sync.WaitGroup
		wg.Add(len(words))

		var fastSort pagination.FastSort
		_time := utils.ExecTime(func() {
			var allValues = make([]*models.SliceItem, 0)
			_stime := time.Now()
			for _, word := range words {
				go e.SimpleSearch(word, utils.StringToInt(word), func(values []*models.SliceItem) {
					lock.Lock()
					allValues = append(allValues, values...)
					lock.Unlock()
					wg.Done()
				})
			}

			wg.Wait()
			if e.IsDebug {
				logger.Logger.Infoln("获取ids时间", time.Since(_stime), "ms")
			}
			fastSort = pagination.FastSort{ScoreMap: make(map[uint32]*pagination.DocScore, int(float32(len(allValues))*1.5))} // 预分配空间，优化效率
			fastSort.Add(allValues)
			if e.IsDebug {
				logger.Logger.Infoln("初始化fastsort时间", time.Since(_stime), "ms")
			}
		})
		if e.IsDebug {
			logger.Logger.Infoln("搜索时间:", _time, "ms")
		}
		tim := time.Now()
		trie.BuildTree(request.FilterWord)

		// 设置fail指针
		trie.SetNodeFailPoint()

		if e.IsDebug {
			logger.Logger.Infoln("ACTrie init Success! 耗时: ", time.Since(tim))
		}
		_tt := utils.ExecTime(func() {
			resultIds, totlen = fastSort.GetAll()

		})
		if e.IsDebug {
			logger.Logger.Infoln("处理排序耗时", _tt, "ms")
			logger.Logger.Infoln("结果集大小", len(resultIds))
		}
		fmt.Println("搜索后时间", time.Since(StartTime))
		//如果结果集长度超过5w，则默认其为慢查询，将其加入redis与badger中，过期时间设为1小时
		if totlen > 50000 {
			logger.Logger.Infoln("出现慢查询,加入redis与badgerdb时间为", time.Now())
			redis.Set(ks)
			resultIdslen := len(resultIds)
			slowslice := make([]*results.SliceItem, 0, resultIdslen)
			for index := 0; index < resultIdslen; index++ {
				slowslice = append(slowslice, &results.SliceItem{resultIds[index].Id, resultIds[index].Score})
			}
			slowslice = append(slowslice, &results.SliceItem{uint32(totlen), 0})
			Slowresult := results.Result{slowslice}
			bufs, _ := Slowresult.MarshalBinary()
			e.SlowResultStorages.Set(utils.Encoder(ks), bufs)
		}
		// 处理分页
		GetAndSetDefault(request)
		//读取文档
		result = &models.SearchResult{
			Total:     uint32(fastSort.Count()),
			Time:      float32(_time),
			Page:      request.Page,
			Limit:     request.Limit,
			Words:     words,
			Documents: nil,
		}
	}
	_time := utils.ExecTime(func() {

		pager.Init(int(request.Limit), len(resultIds))
		//设置总页数
		result.PageCount = uint32(pager.PageCount)

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(int(request.Page))
			if start == -1 {
				return
			}
			items := resultIds[start : end+1]
			if e.IsDebug {
				logger.Logger.Infoln("Page: ", "start ", start, "end ", end)
			}

			result.Documents = make([]*models.ResponseDoc, len(items))

			var wg sync.WaitGroup
			wg.Add(len(items))
			fmt.Println("lenitems:", len(items))
			//只读取前面 limit 个
			_tt := time.Now()
			for index, item := range items {
				go func(index int, item pagination.DocItem) {
					result.Documents[index] = &models.ResponseDoc{}
					defer wg.Done()

					_cost := time.Now()
					buf := e.GetDocById(item.Id)

					if e.IsDebug {
						logger.Logger.Infoln("Id: ", item.Id, "--- GetDocById: ", time.Since(_cost), "--- 数据长度: ", len(buf))
					}

					storageDoc := new(doc.StorageIndexDoc)
					if buf != nil {
						storageDoc.UnmarshalBinary(buf)
						// 查找 ACTrie, 如果有过滤词, 过滤掉
						if request.FilterWord != nil &&
							len(request.FilterWord) > 0 &&
							trie.AcAutoMatch(storageDoc.Text) {
							return
						}
					}
					result.Documents[index].Score = item.Score

					if buf != nil {
						result.Documents[index].Id = item.Id
						result.Documents[index].Url = storageDoc.Url
						text := storageDoc.Text
						//处理关键词高亮
						highlight := request.Highlight
						if highlight != nil {
							//全部小写
							text = strings.ToLower(text)
							for _, word := range words {
								text = strings.ReplaceAll(
									text, word,
									fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
							}
						}
						result.Documents[index].Text = text
						//判断是否收藏
						result.Documents[index].Islike = false
						if val, ok := request.Likes[item.Id]; ok {
							result.Documents[index].Islike = true
							result.Documents[index].Docsid = val.Docid
							result.Documents[index].Favid = val.Favid
						}
					}
				}(index, item)
			}

			wg.Wait()

			if e.IsDebug {
				logger.Logger.Infoln("分页耗时: ", time.Since(_tt))
			}

		}
	})
	if e.IsDebug {
		logger.Logger.Infoln("处理数据耗时：", _time, "ms")
	}
	SearchType := "搜索文本-正常集合"
	if isfound {
		SearchType = "搜索文本-从boltdb中取出"
	}
	resultText := "当前时间: " + fmt.Sprintf("%v", time.Now().Format("2006/01/02 15:04:05")) + " --- 搜索数据: " + request.Query + " --- 搜索时间: " + fmt.Sprintf("%v", (time.Since(StartTime).Seconds())*1000) + "ms" + " --- 搜索类型: " + SearchType
	e.SearchResultMQChan <- resultText
	result.Time = float32(time.Since(StartTime).Seconds() * 1000)
	return result
}
func GetAndSetDefault(s *models.SearchRequest) {

	if s.Limit == 0 {
		s.Limit = 20
	}
	if s.Page == 0 {
		s.Page = 1
	}

	if s.Highlight == nil {
		s.Highlight = &models.Highlight{
			PreTag:  "<span style=\"color: red;\">",
			PostTag: "</span>",
		}
	}
}
func (e *Engine) MultiSearchPicture(request *models.SearchRequest) *models.SearchPictureResult {
	//等待搜索初始化完成
	if e.IsDebug {
		logger.Logger.Infoln("Search start ------------------------")
	}

	e.Wait()
	//搜索开始时间
	StartTime := time.Now()
	//query的请求
	ks := request.Query
	e.SearchRequestMQChan <- ks
	//分页
	pager := new(pagination.Pagination)
	//结果集
	var resultIds []pagination.DocItem
	//结果集总长度
	var totlen int
	//返回结果
	var result *models.SearchPictureResult
	//分词
	var words []string
	//是否在redis中
	isfound := redis.Isfound(ks)
	if isfound {
		logger.Logger.Infoln("该结果集从boltdb取出")
		//过滤
		trie.BuildTree(request.FilterWord)

		// 设置fail指针
		trie.SetNodeFailPoint()
		// 处理分页
		GetAndSetDefault(request)
		words = e.WordCut(request.Query)
		if e.IsDebug {
			logger.Logger.Infoln("分词数: ", len(words))
		}
		//读取文档
		result = &models.SearchPictureResult{
			Total:     0,
			Time:      0,
			Page:      request.Page,
			Limit:     request.Limit,
			Words:     words,
			Documents: nil,
		}
		bufs, _ := e.SlowResultStorages.Get(utils.Encoder(ks))
		slowresult := results.Result{}
		slowresult.UnmarshalBinary(bufs)
		resultslice := slowresult.Sliceresult
		slowresultlen := int(resultslice[len(resultslice)-1].Id)
		for index := 0; index < len(resultslice)-1; index++ {
			resultIds = append(resultIds, pagination.DocItem{resultslice[index].Id, resultslice[index].Score, 1})
		}
		result.Total = uint32(slowresultlen)
	} else {
		logger.Logger.Infoln("该结果集为正常搜索")

		//分词搜索
		words = e.WordCut(request.Query)
		if e.IsDebug {
			logger.Logger.Infoln("分词数: ", len(words))
		}

		var lock sync.Mutex
		var wg sync.WaitGroup
		wg.Add(len(words))

		var fastSort pagination.FastSort
		_time := utils.ExecTime(func() {
			var allValues = make([]*models.SliceItem, 0)

			for _, word := range words {
				go e.SimpleSearch(word, utils.StringToInt(word), func(values []*models.SliceItem) {
					lock.Lock()
					allValues = append(allValues, values...)
					lock.Unlock()
					wg.Done()
				})
			}

			wg.Wait()
			fastSort = pagination.FastSort{ScoreMap: make(map[uint32]*pagination.DocScore, int(float32(len(allValues))*1.5))} // 预分配空间，优化效率
			fastSort.Add(allValues)
		})
		if e.IsDebug {
			logger.Logger.Infoln("搜索时间:", _time, "ms")
		}
		tim := time.Now()
		trie.BuildTree(request.FilterWord)

		// 设置fail指针
		trie.SetNodeFailPoint()

		if e.IsDebug {
			logger.Logger.Infoln("ACTrie init Success! 耗时: ", time.Since(tim))
		}
		_tt := utils.ExecTime(func() {
			resultIds, totlen = fastSort.GetAll()
		})
		if e.IsDebug {
			logger.Logger.Infoln("处理排序耗时", _tt, "ms")
			logger.Logger.Infoln("结果集大小", len(resultIds))
		}
		if totlen > 50000 {
			logger.Logger.Infoln("出现慢查询,加入redis与badgerdb时间为", time.Now())
			redis.Set(ks)
			resultIdslen := len(resultIds)
			slowslice := make([]*results.SliceItem, 0, resultIdslen)
			for index := 0; index < resultIdslen; index++ {
				slowslice = append(slowslice, &results.SliceItem{resultIds[index].Id, resultIds[index].Score})
			}
			//最后一个为原数据长度
			slowslice = append(slowslice, &results.SliceItem{uint32(totlen), 0})
			Slowresult := results.Result{slowslice}
			bufs, _ := Slowresult.MarshalBinary()
			e.SlowResultStorages.Set(utils.Encoder(ks), bufs)
		}
		// 处理分页
		GetAndSetDefault(request)

		//读取文档
		result = &models.SearchPictureResult{
			Total:     uint32(fastSort.Count()),
			Time:      float32(_time),
			Page:      request.Page,
			Limit:     request.Limit,
			Words:     words,
			Documents: nil,
		}
	}
	_time := utils.ExecTime(func() {

		pager.Init(int(request.Limit), len(resultIds))
		//设置总页数
		result.PageCount = uint32(pager.PageCount)
		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(int(request.Page))
			if start == -1 {
				return
			}

			items := resultIds[start : end+1]
			if e.IsDebug {
				logger.Logger.Infoln("Page: ", "start ", start, "end ", end)
			}

			result.Documents = make([]*models.ResponseUrl, len(items))
			var wg sync.WaitGroup
			wg.Add(len(items)) // 并发上传图片

			//只读取前面 limit 个
			_tt := time.Now()
			for index, item := range items {
				go func(index int, item pagination.DocItem) {
					defer wg.Done()
					_cost := time.Now()
					buf := e.GetDocById(item.Id)

					if e.IsDebug {
						logger.Logger.Infoln("Id: ", item.Id, "--- GetDocById: ", time.Since(_cost), "--- 数据长度: ", len(buf))
					}
					result.Documents[index] = &models.ResponseUrl{}
					result.Documents[index].Score = item.Score

					if buf != nil {
						storageDoc := new(doc.StorageIndexDoc)
						storageDoc.UnmarshalBinary(buf)

						result.Documents[index].Id = item.Id
						result.Documents[index].Url = storageDoc.Url
						e.GetPictureUrl(storageDoc.Url, item.Id, &result.Documents[index].ThumbnailUrl, func() {
						}) // get url
						text := storageDoc.Text
						//处理关键词高亮
						highlight := request.Highlight
						if highlight != nil {
							//全部小写
							text = strings.ToLower(text)
							for _, word := range words {
								text = strings.ReplaceAll(text, word, fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
							}
						}
						result.Documents[index].Text = text
						//判断是否收藏
						result.Documents[index].Islike = false
						if val, ok := request.Likes[item.Id]; ok {
							result.Documents[index].Islike = true
							result.Documents[index].Docsid = val.Docid
							result.Documents[index].Favid = val.Favid
						}
					}
				}(index, item)
			}
			wg.Wait() // 等待并发结束

			if e.IsDebug {
				logger.Logger.Infoln("分页耗时: ", time.Since(_tt))
			}
		}
	})

	if e.IsDebug {
		logger.Logger.Infoln("处理数据耗时：", _time, "ms")
	}
	SearchType := "搜索图片-正常集合"
	if isfound {
		SearchType = "搜索图片-从boltdb中取出"
	}
	resultText := "当前时间: " + fmt.Sprintf("%v", time.Now().Format("2006/01/02 15:04:05")) + " --- 搜索数据: " + request.Query + " --- 搜索时间: " + fmt.Sprintf("%v", (time.Since(StartTime).Seconds())*1000) + "ms" + " --- 搜索类型: " + SearchType
	e.SearchResultMQChan <- resultText
	result.Time = float32(time.Since(StartTime).Seconds() * 1000)
	return result
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id uint32) []byte {
	key := utils.Uint32ToBytes(id)
	buf, found := e.DocStorages[id%BadgerShard].Get(key)
	if found {
		return buf
	}

	return nil
}

// Close will save AVL data
func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	////保存文件
	//e.flushDataIndex()

	for _, db := range e.DocStorages {
		db.Close()
	}
	e.KeyIdsStorage.Close()
}
func pictureSearchHandleError(err error) {
	logger.Logger.Infoln("Error:", err)
	os.Exit(-1)
}

func pictureSearchGetRemote(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		// 如果有错误返回错误内容
		return nil, err
	}
	// 使用完成后要关闭，不然会占用内存
	defer res.Body.Close()
	// 读取字节流
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return byteArray, err
}
func pictureCompress(buf []byte) ([]byte, error) {
	//文件压缩
	decodeBuf, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	// 修改图片的大小
	set := resize.Resize(0, 200, decodeBuf, resize.Lanczos3)
	NewBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&NewBuf, set)
	case "jpeg", "jpg":
		err = jpeg.Encode(&NewBuf, set, &jpeg.Options{Quality: 80})
	default:
		return nil, errors.New("该图片格式不支持压缩")
	}
	if err != nil {
		return nil, err
	}
	if NewBuf.Len() < len(buf) {
		buf = NewBuf.Bytes()
	}
	return buf, nil
}

func (e *Engine) putInOSS(url string, id uint32) bool {
	resByte, err := pictureSearchGetRemote(url)
	if err != nil {
		logger.Logger.Infoln(err)
	}
	resBytes, err := pictureCompress(resByte)
	if err != nil {
		logger.Logger.Infoln(err)
		return false
	}
	//Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := ""
	accessKeySecret := ""
	bucketName := "lgdsearch"
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		pictureSearchHandleError(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		pictureSearchHandleError(err)
	}

	reader := bytes.NewReader(resBytes)
	pictureUrl := fmt.Sprintf("example/%d.jpg", id)
	err = bucket.PutObject(pictureUrl, reader)
	if err != nil {
		pictureSearchHandleError(err)
	}
	redis.Upload(strconv.Itoa(int(id)))
	return true
}

func (e *Engine) GetPictureUrl(url string, id uint32, thumbnailUrl *string, call func()) {
	ok := redis.IsUpload(strconv.Itoa(int(id)))
	if !ok {
		// 压缩图片到服务器
		ok = e.putInOSS(url, id)
	}
	if ok {
		*thumbnailUrl = fmt.Sprintf("%d.jpg", id)
	} else {
		*thumbnailUrl = url
	}
	fmt.Println(&thumbnailUrl)
	call()
}
