package main

import (
	"fmt"
	"time"
)

func main() {
	sum := 1 + 1
	fmt.Println(sum)

	var array []string
	array = append(array, "ARROZ", "FEIJÃO", "MACARRONADA")

	fmt.Println(array)
	for _, value := range array {
		fmt.Println(value)
	}

	aux := "AUX"
	fmt.Println(aux)

	sumTwo, name := SumTwoNumberAndReturnName(1, 2, "Veronica")
	fmt.Println(sumTwo)
	fmt.Println(name)

	go start("Victor")
	go finish("Veronica")

	time.Sleep(4 * time.Second)
}

func SumTwoNumberAndReturnName(numberOne, numberTwo int, name string) (int, string) {
	return numberOne + numberTwo, name
}

func start(str string) {
	fmt.Println("START - " + str)
	time.Sleep(3 * time.Second)
	fmt.Println("START")

}

func finish(str string) {
	fmt.Println("FINISH -  " + str)
	time.Sleep(4 * time.Second)
	fmt.Println("FINISH")
}
