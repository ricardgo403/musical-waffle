package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string, quit chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			select {
			case c <-fmt.Sprintf("%s: %d", msg, i):

			case <-quit:
				return

			default:
				time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			}

		}
	}()
	return c
}

func main() {
	quit := make(chan bool)
	c := boring("Joel", quit)
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-time.After(1 * time.Second):
			quit <- true
			return
		}
	}
}