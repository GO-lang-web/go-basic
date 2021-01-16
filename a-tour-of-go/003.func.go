package main 

import "fmt"

//define a func 
func add(x int, y int) int {
	return x + y 
}

//When two or more consecutive named function parameters share a type, you can omit the type from all but the last.
func add2(x,y int) int{
	return x + y 
}


//A function can return any number of results.
func swap(x,y string) (string, string){
	return y, x
}

//Go's return values may be named. 
//If so, they are treated as variables defined at the top of the function.
func split(sum int) (x,y int){
	x = sum *4/9
	y = sum - x 
	return 
}

func main(){
	fmt.Println(add(1,2))//3
	fmt.Println(add(1,4))//5

	a ,b := swap("hello", "world")//world hello
	fmt.Println(a,b)
	fmt.Println(split(17))// 7 10 
}



