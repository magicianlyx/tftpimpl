package tftp

import (
	"sync"
	"fmt"
)

type QueueHandle struct {
	cmap     sync.Map
}

func NewQueueHandle() *QueueHandle {
	return &QueueHandle{
		sync.Map{},
	}
}

func (q *QueueHandle) Add(host string, port int, f func(article *Article)) {
	if f == nil {
		return
	}
	key := fmt.Sprintf("%s:%d", host, port)
	q.cmap.Store(key, f)
}

func (q *QueueHandle) Del(host string, port int) {
	key := fmt.Sprintf("%s:%d", host, port)
	q.cmap.Delete(key)
}

func (q *QueueHandle) PutArticle(article *Article) {
	q.cmap.Range(func(key, value interface{}) bool {
		k, kok := key.(string)
		if !kok {
			return true
		}
		v, vok := value.(func(article *Article))
		if !vok {
			return true
		}
		
		if fmt.Sprintf("%s:%d", article.Host, article.Port) == k {
			go v(article)
		}
		return true
	})
}
