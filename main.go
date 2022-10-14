package main

import (
	"fmt"
	"sync"
)

type Object struct {
	name string
}

type Arr struct {
	arr []*Object
}

func main() {
	c := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		select {
		case a := <-c:
			fmt.Println(a)
			wg.Done()
		}
	}()
	go func() {
		select {
		case a := <-c:
			fmt.Println(a)
			wg.Done()
		}
	}()

	c <- 4
	//c <- 3
	wg.Wait()
}
