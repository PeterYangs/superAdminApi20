package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	for i := 0; i < 10; i++ {

		wait.Add(1)

		go func(index int) {

			defer wait.Done()

			<-cxt.Done()

			if index == 0 {

				time.Sleep(10 * time.Second)
			}

			fmt.Println("11111")

		}(i)

	}

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		//time.Sleep(3 * time.Second)
		fmt.Println("结束完毕")

		cancel()
	}()

	wait.Wait()

}
