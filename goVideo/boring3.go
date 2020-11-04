package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func fanIn(input <-chan string, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

func main() {
	c := fanIn(boring("joel"), boring("abril"))
	for i := 0; i < 10; i++ {
		fmt.Printf("Somebody said: %q\n", <-c)
	}
	fmt.Println("You're both boring. I'm leaving.")
}
