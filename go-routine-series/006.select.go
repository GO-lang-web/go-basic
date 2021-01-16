package main 

import(
	"strconv"
	"time"
	"fmt"
)


func sample(ch chan string){
	for i:= 0; i< 5; i++ {
		ch <- "I'm sample num: " + strconv.Itoa(i)
		time.Sleep( 1* time.Second)
	}
}

func sample2(ch chan int){
	for i:= 0; i< 3; i++ {
		ch <- i
		time.Sleep( 2* time.Second)
	}

}

func main(){
	chan1 := make( chan string , 3)
	chan2 := make( chan int , 5)
	
	for i:=0; i<10; i++ {
		go sample(chan1)
		go sample2(chan2)
	}

	//类似switch 
	for i:=0; i<100; i++ {
		select {
		case str , isOk1  := <- chan1 :
			if !isOk1 {
				fmt.Println("ch1 failed")
			}
			fmt.Println( str )
		case str2 , isOk2  := <- chan1 :
			if !isOk2 {
				fmt.Println("ch2 failed")
			}
			fmt.Println( str2 )
		}
	}
	
	time.Sleep(60 * time.Second)
}