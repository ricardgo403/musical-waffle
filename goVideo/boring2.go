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

func main() {
	c := boring("abril :( !")
	channel := boring("joel :)")
	for i := 0; i < 10; i++ {
		fmt.Printf("Abril said: %q\n", <-c)
		fmt.Printf("Joel said: %q\n", <-channel)
	}
	fmt.Println("You're both boring. I'm leaving.")
}
