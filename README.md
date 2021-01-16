# go-basic

## editor

- goLand
- vscode

## gotoutine

> 如何开启？

- 1. go + 匿名函数

```go
message := make(chan string)
```

- 2. go + 函数名

```go
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
}s
```

> 如何输入输出？

- 声明 chan
- chan <- 往 chan 里写
- <- chan 从 chan 里面输出
- 多个的话 指定 chan 长度
