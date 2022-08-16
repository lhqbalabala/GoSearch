package pagination

import (
	"Search_RPC/models"
	"container/heap"
	"sort"
)

type ScoreSlice []DocItem

func (x ScoreSlice) Len() int {
	return len(x)
}

func (x ScoreSlice) Less(i, j int) bool {
	if x[i].Count != x[j].Count {
		return x[i].Count < x[j].Count
	} else {
		return x[i].Score < x[j].Score
	}
}

func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type DocScore struct {
	// 分词总分数
	score float32

	// 符合分词个数
	count uint32
}

type DocItem struct {
	Id    uint32
	Score float32
	Count uint32
}

type FastSort struct {
	ScoreMap map[uint32]*DocScore
}

// Add :Count the scores of all documents
func (f *FastSort) Add(values []*models.SliceItem) {
	if values == nil {
		return
	}
	for _, item := range values {
		_, ok := f.ScoreMap[item.Id]
		if ok {
			f.ScoreMap[item.Id].score += item.Score
			f.ScoreMap[item.Id].count += 1
		} else {
			f.ScoreMap[item.Id] = &DocScore{score: item.Score, count: 1}
		}
	}
}

// Count 获取数量
func (f *FastSort) Count() int {
	return len(f.ScoreMap)
}

type QueryResHeap []*DocItem

func (h QueryResHeap) Len() int { return len(h) }

func (h QueryResHeap) Less(i, j int) bool {
	return h[i].compare(h[j])
}

func (h QueryResHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *QueryResHeap) Push(x interface{}) {
	*h = append(*h, x.(*DocItem))
}

func (h *QueryResHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (node1 *DocItem) compare(node2 *DocItem) bool {
	if node1.Count == node2.Count {
		return node1.Score < node2.Score
	}
	return node1.Count < node2.Count
}

const MaxResCount = 50000

// GetAll 获取按 score 排序后的结果集
func (f *FastSort) GetAll() (res []DocItem, totlen int) {

	var result = make([]DocItem, len(f.ScoreMap))
	totlen = len(result)
	index := 0
	if len(f.ScoreMap) > MaxResCount {
		minHeap := QueryResHeap{}
		for key, value := range f.ScoreMap {
			if key == 0 {
				continue
			}
			result[index] = DocItem{Id: key, Score: value.score, Count: value.count}
			index++
		}
		for i := 0; i < len(result); i++ {
			if i < MaxResCount {
				minHeap = append(minHeap, &result[i])
				if i == MaxResCount-1 {
					heap.Init(&minHeap)
				}
			} else if minHeap[0].compare(&result[i]) {
				heap.Pop(&minHeap)
				heap.Push(&minHeap, &result[i])
			}

		}
		result = make([]DocItem, MaxResCount)
		for i := 50000 - 1; i >= 0; i-- {
			val := heap.Pop(&minHeap)
			result[i] = *(val.(*DocItem))
		}

	} else {
		for key, value := range f.ScoreMap {
			if key == 0 {
				continue
			}
			result[index] = DocItem{Id: key, Score: value.score, Count: value.count}
			index++
		}

		sort.Sort(sort.Reverse(ScoreSlice(result)))

	}
	return result, totlen
}
