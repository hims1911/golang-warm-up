package main

import (
	"fmt"
	"sync"
	"time"
)

var n = 5
var mu sync.Mutex

func main() {
	// creating a channel make(chan int)
	ch := make(chan int)
	go square(2, ch)
	result := <-ch
	fmt.Println(result)

	time.Sleep(1 * time.Second) // reason to add sleep is thread will run concurrently but main execution can stop
	// fmt.Println("execution done")

	a := []int{1, 2, 3}
	ch = make(chan int, len(a))
	// go square with array will send channel chan int
	go squareWithArray(a, ch)
	for i := 0; i < len(a); i++ {
		fmt.Printf("Result: %v \n", <-ch)
	}
	time.Sleep(1 * time.Second)

	// for c := range ch {
	// 	fmt.Println(c) // deadlock example deadlock will happen because we are not closing the channel close(ch)
	// }

	ch1 := make(chan int, len(a))
	ch2 := make(chan int, len(a))

	go timesThree(a, ch1)
	go minusThree(a, ch2)

	for i := 0; i < 2*len(a); i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("ch1 got the value", msg1)
		case msg2 := <-ch2:
			fmt.Println("ch2 got the value", msg2)
		default:
			fmt.Println("got the default use case")
		}
	}

	time.Sleep(2 * time.Second)

	// Mutual Exclusion Example
	for i := range 10 {
		fmt.Println(i) // should run sequentially
		go squareForMutualExclusion(i)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Mutual Exclusion Resolution")

	n = 5
	// solving mutual exclusion is all about adding sync.mutex package
	for i := range 10 {
		fmt.Println(i)
		go squareForMutualExclusionResolution(i)
	}

	time.Sleep(1 * time.Second)
}

// square will calculate the square
func square(data int, ch chan int) {
	fmt.Println(data * data)
	ch <- data * data * data

}

func squareWithArray(a []int, ch chan int) {
	minusCh := make(chan int, len(a))
	for _, v := range a {
		val := v * 3
		if val%3 == 0 {
			go minusChFromSquareArray(val, minusCh)
			val = <-minusCh
		}
		ch <- val
	}
}

func minusChFromSquareArray(val int, ch chan int) {
	ch <- val - 3
	fmt.Println("The functions continues after returning the result")
}

// squareForMutualExclusion() will help to perform mutual exclusion
// it is the critical section of the program
func squareForMutualExclusion(i int) {
	fmt.Println("Got executed", i)
	n = n * 2
	fmt.Println(n)
}

// squareForMutualExclusion() will help to perform mutual exclusion
// it is the critical section of the program
func squareForMutualExclusionResolution(i int) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Got executed", i)
	n = n * 2
	fmt.Println(n)
}

func timesThree(arr []int, ch chan int) {
	for _, elem := range arr {
		ch <- elem * 3
	}
}

func minusThree(arr []int, ch chan int) {
	for _, elem := range arr {
		ch <- elem - 3
	}
}
