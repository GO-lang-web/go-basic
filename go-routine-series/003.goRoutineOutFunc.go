package main

import (
	"fmt"
	"time"
)

var  message = make(chan string)
func test(){
	message <- "hello goroutine"
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
	fmt.Println("test over")
}

// hello goroutine another
// test over
