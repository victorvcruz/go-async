package main

import "fmt"

func main() {
	var f float64 = 329.030000
	var str string
	str = fmt.Sprintf("%v", f)

	fmt.Println(str)
}
