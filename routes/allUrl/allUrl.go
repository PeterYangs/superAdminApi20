package allUrl

import (
	"strings"
	"sync"
)

type allUrl struct {
	urls []string
	lock sync.Mutex
}

var AllUrl *allUrl

var once = sync.Once{}

// NewAllUrl 单例
func NewAllUrl() *allUrl {

	once.Do(func() {

		AllUrl = &allUrl{
			urls: []string{},
			lock: sync.Mutex{},
		}

	})

	return AllUrl
}

func (a *allUrl) Add(url string) {

	a.lock.Lock()

	defer a.lock.Unlock()

	a.urls = append(a.urls, url)

}

func (a *allUrl) Search(keyword string) []string {

	a.lock.Lock()

	defer a.lock.Unlock()

	var list []string

	for _, url := range a.urls {

		if strings.Contains(url, keyword) {

			list = append(list, url)
		}

		if len(list) >= 10 {

			return list
		}

	}

	return list

}
