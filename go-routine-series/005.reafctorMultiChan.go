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

func test( message chan string ){
	message <- "hello goroutine1"
	message <- "hello goroutine2"
	message <- "hello goroutine3"
	message <- "hello goroutine4"
}
func test2( message chan string ){
	time.Sleep(2* time.Second)
	anotherStr := <-message
	anotherStr = anotherStr + " another"
	message<-anotherStr

	//需要关闭channel
	close(message)
}

func main() {
	var  message = make(chan string, 3)

	go test(message)
	go test2(message)
	time.Sleep(3*time.Second)
	//usage  of range 


	// fmt.Println(<-message)
	// fmt.Println(<-message)
	// fmt.Println(<-message)
	// fmt.Println(<-message)
	//message 为 nil 是否为继续执行？ 会的 all goroutines are asleep - deadlock!
	for str := range message {
		fmt.Println( str )
	}
	fmt.Println("test over")
}



