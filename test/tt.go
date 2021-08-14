package main

import "fmt"

func main() {

	j := make(chan int)

	for {

		select {

		case <-j:

			//default:
			//
			//	fmt.Println(1111)

		}

		fmt.Println(1111)
	}

}
