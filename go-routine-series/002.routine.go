package main

import (
	"fmt"
	"time"
)

func main() {
	//main 里面的执行顺序是从上到下的 协程之间不确定
	message := make(chan string)
	//need channel 
	go func(){
		//字符串写入message
		message <- "hello goroutine"
		// fmt.Println("hello goroutine")
	}()

	go func(){
			time.Sleep(2* time.Second)
			anotherStr := <-message
			anotherStr = anotherStr + " another"

			message<-anotherStr
	}()
	// str := <-message
	// fmt.Println(str)
	time.Sleep(3*time.Second)
	fmt.Println(<-message)
	//输出message 里面的内容
	// fmt.Println( <-message)
	// fmt.Println("Hello, 世界")
}

//chan <- 往里面写
// <- 从chan里面往外输出

// hello goroutine
// Hello, 世界
