package main

import "fmt"

func main() {
	// fmt.Println("Hello, Shakil !")
	// a:=10
	// fmt.Println(a)
	/*
	int 
	float32
	string
	bool
	*/
	// var x int = 20
	// fmt.Println(x)

	//var x int = 10
	// a :=10
	// a := "Hello, shakil
	// 
	age := 21
	sex := "male"
	if age > 20 && sex == "male"{
		fmt.Println("You are eligible to vote")
	} else{
		fmt.Println("You are not eligible to vote")
	}

	switch age {
	case 18:
		fmt.Println("You are eligible to vote")
	default:

		fmt.Println("You are not eligible to vote")
	}
}
