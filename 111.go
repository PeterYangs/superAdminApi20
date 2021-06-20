package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Downloader struct {
	io.Reader
	Total   int64
	Current int64
}

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.Current += int64(n)
	fmt.Printf("\r正在下载，下载进度：%.2f%%", float64(d.Current*10000/d.Total)/100)
	if d.Current == d.Total {
		fmt.Printf("\r下载完成，下载进度：%.2f%%", float64(d.Current*10000/d.Total)/100)
	}
	return
}

func downloadFile(url, filePath string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	file, err := os.Create(filePath)
	defer func() {
		_ = file.Close()
	}()
	downloader := &Downloader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
	}
	if _, err := io.Copy(file, downloader); err != nil {
		log.Fatalln(err)
	}
}

var wg sync.WaitGroup

func main() {
	task := make(map[string]string)
	task["https://dl.softmgr.qq.com/original/game/WeGameSetup3.32.4.6183_gjwegame_0_0.exe"] = "d:/wegame.exe"
	for k, v := range task {
		wg.Add(1)
		downloadFile(k, v)
	}
	wg.Wait()
}
