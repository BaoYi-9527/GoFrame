package main

import "fmt"

func main() {
	//fmt.Println("before panic")
	//panic("crash")
	//fmt.Println("after panic")

	//arr := []int{1, 2, 3}
	//fmt.Println(arr[4])

	testRecover()
	fmt.Println("after recover")
}

func testRecover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}
