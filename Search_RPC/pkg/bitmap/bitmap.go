package bitmap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Bitmap struct {
	words []byte
}

var bitmap *Bitmap = &Bitmap{}
var mu sync.Mutex
var change bool

func InitBitmap(filepath string) {
	Load(filepath)
	go bitmap.AutomaticFlush(filepath) // 刷新到磁盘
}
func GetBitmap() *Bitmap {
	return bitmap
}
func Getlen() int {
	return len(bitmap.words)
}
func (bitmap *Bitmap) Has(num int) bool {
	word, bit := num>>3, uint32(num&(7))
	return word < len(bitmap.words) && (bitmap.words[word]&(1<<bit)) != 0
}

func (bitmap *Bitmap) Add(num int) {
	word, bit := num>>3, uint(num&(7))
	for word >= len(bitmap.words) {
		bitmap.words = append(bitmap.words, 0)
	}
	// 判断num是否已经存在bitmap中
	if bitmap.words[word]&(1<<bit) == 0 {
		bitmap.words[word] |= 1 << bit
		change = true
	}
}

func (bitmap *Bitmap) String() string {
	var buf bytes.Buffer
	var length int
	buf.WriteByte('{')
	for i, v := range bitmap.words {
		if v == 0 {
			continue
		}
		for j := uint(0); j < 8; j++ {
			if v&(1<<j) != 0 {
				length++
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", uint(i)<<3+j)
			}
		}
	}
	buf.WriteByte('}')
	fmt.Fprintf(&buf, "\nLength: %d", length)
	return buf.String()
}

// 自动保存修改后的bitmap，10min检测一次
func (bitmap *Bitmap) AutomaticFlush(filepath string) {
	ticker := time.NewTicker(time.Minute * 10)

	for {
		if change {
			<-ticker.C
			Flush(filepath)
			change = !change
		}
	}

}
func Load(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	bitmap.words = buf
}

func Flush(filepath string) {
	file, _ := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600) // 清空
	file.Close()

	err := ioutil.WriteFile(filepath, bitmap.words, 0600)
	if err != nil {
		panic(err)
	}
}
