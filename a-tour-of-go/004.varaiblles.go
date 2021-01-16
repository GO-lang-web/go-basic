package main

import "fmt"

var c, python, java bool


// Variables with initializers
var  a, b int = 1, 2


func test(){
	var i, j int = 1, 2
	k := 3
	c, python, java := true, false, "no!"

	fmt.Println(i, j, k, c, python, java)
}

func main() {
	var i int
	fmt.Println(i, c, python, java)	//0 false false falses
	fmt.Println(a,b) // 1 2 


	//Inside a function, the := short assignment statement can be used 
	//in place of a var declaration with implicit type.
	test()//1 2 3 true false no!
}
