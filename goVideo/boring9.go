package main

import (
	"fmt"
	"math/rand"
)

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <-fmt.Sprintf("%s: %d", msg, i):

			case <-quit:
				quit <- "see ya!"
				return
			}

		}
	}()
	return c
}

func main() {
	quit := make(chan string)
	c := boring("Joel", quit)
	for i := rand.Intn(10); i >= 0 ; i-- {
		fmt.Println(<-c)
	}
	quit <- "bye!"
	fmt.Printf("Somebody says: %q\n", <-quit)
}