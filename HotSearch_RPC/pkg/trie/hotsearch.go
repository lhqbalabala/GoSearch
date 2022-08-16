package trie

import (
	"HotSearch_RPC/pkg/redis"
	"HotSearch_RPC/pkg/utils"
	"fmt"

	"os"
	"runtime"
	"sync"
	"time"
)

const HotSearchFileName = "HotSearch.txt"

// hotSearch use a map and a queue to get top 10 hotSearch message from today's search message
type queueNode struct {
	TimeMessage time.Time
	Text        string
	Next        *queueNode
}

type queue struct {
	size int
	head *queueNode
	end  *queueNode
	sync.Mutex
}

func (q *queue) Size() int {
	return q.size
}

func (q *queue) Push(node *queueNode) {
	q.Lock()
	defer q.Unlock()
	q.size += 1
	if q.size == 1 {
		q.head = node
	} else {
		q.end.Next = node
	}
	q.end = node
}

func (q *queue) Pop() *queueNode {
	q.Lock()
	defer q.Unlock()
	if q.size == 0 {
		return q.head
	}
	top := q.head
	q.head = q.head.Next
	q.size -= 1
	if q.size == 0 {
		q.end = nil
	}
	return top
}

func (q *queue) Top() *queueNode {
	return q.head
}

var MessageChan chan string

type HotSearch struct {
	searchQueue *queue // queue 维护map中的数据
	sync.Mutex
}

var MyHotSearch *HotSearch

func InitHotSearch(filepath string) {
	GetHotSearch()
	GetHotSearch().Load(filepath)
	go GetHotSearch().InQueueExec()
	go GetHotSearch().OutQueueExec()
	go GetHotSearch().AutoReGetArray(filepath)
}

func GetHotSearch() *HotSearch {
	if MyHotSearch == nil {
		MyHotSearch = &HotSearch{
			searchQueue: &queue{},
		}
		MessageChan = make(chan string, 1000)
	}
	return MyHotSearch
}

func (hot *HotSearch) Queue() *queue {
	return hot.searchQueue
}

func (hot *HotSearch) showQueueElement() {
	fmt.Println("Queue start -----------")
	index := 0
	for head := hot.searchQueue.head; head != nil; head = head.Next {
		fmt.Println("index: ", index, "Text: ", head.Text, "Time: ", head.TimeMessage)
		index++
	}
	fmt.Println("Queue end -------------")
}

type HotSearchMessage struct {
	Text string `json:"text,omitempty"`
	Num  int    `json:"num,omitempty"`
}

func (node1 *HotSearchMessage) compare(node2 *HotSearchMessage) bool {
	return node1.Num < node2.Num
}

func (hot *HotSearch) OutQueueExec() {
	for {
		// 保留 24 小时内数据
		if hot.searchQueue.Size() != 0 &&
			time.Now().Sub(hot.searchQueue.Top().TimeMessage).Hours() > 24. {
			node := hot.searchQueue.Pop()
			redis.HotSearchZreduce(node.Text)
		}
	}
}

func SendText(text string) {
	MessageChan <- text
}

func (hot *HotSearch) InQueueExec() {
	for {
		text := <-MessageChan
		redis.HotSearchZadd(text)
		hot.searchQueue.Push(&queueNode{TimeMessage: time.Now(), Text: text})
	}
}

// AutoReGetArray 60S 自动更新一次,flush持久化
func (hot *HotSearch) AutoReGetArray(filepath string) {
	ticker := time.NewTicker(time.Second * 60)
	head := hot.searchQueue.head
	size := hot.searchQueue.Size()

	for {
		<-ticker.C
		if head == hot.searchQueue.head && size == hot.searchQueue.Size() {
			continue
		}
		head = hot.searchQueue.head
		size = hot.searchQueue.Size()
		hot.Flush(filepath)

		runtime.GC()
	}
}

func (hot *HotSearch) Flush(filepath string) {
	data := make([]queueNode, 0)
	for head := hot.searchQueue.head; head != nil; head = head.Next {
		data = append(data, queueNode{Text: head.Text, TimeMessage: head.TimeMessage, Next: nil})
	}
	file, _ := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600) // 清空
	file.Close()
	fmt.Println("更新一次", filepath)
	utils.Write(&data, filepath)
}

func (hot *HotSearch) Load(filepath string) {
	if filepath == "" {
		filepath = "./pkg/data/HotSearch.txt"
	}
	data := make([]queueNode, 0)
	utils.Read(&data, filepath)
	for _, val := range data { // 这里的val是单个变量，是不可以直接插入的
		hot.searchQueue.Push(&queueNode{Text: val.Text, TimeMessage: val.TimeMessage})
	}
}
