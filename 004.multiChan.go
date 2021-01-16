package main

import (
	"fmt"
	"time"
)

/**
FIFO first in , first out 
{
	[], //head
	[],
	[] //tail 
}
**/
var  message = make(chan string, 3)
func test(){
	message <- "hello goroutine1"
	message <- "hello goroutine2"
	message <- "hello goroutine3"
	message <- "hello goroutine3"
}
func test2(){
	time.Sleep(2* time.Second)
	anotherStr := <-message
	anotherStr = anotherStr + " another"
	message<-anotherStr
}

func main() {
	go test()
	go test2()
	time.Sleep(3*time.Second)
	fmt.Println(<-message)
	fmt.Println(<-message)
	fmt.Println(<-message)
	fmt.Println(<-message)
	fmt.Println("test over")
}

// hello goroutine2
// hello goroutine3
// hello goroutine1 another
// test over
