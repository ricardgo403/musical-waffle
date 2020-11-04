package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str string
	wait chan bool
}

func boring(msg string) <-chan Message {
	waitForIt := make(chan bool)
	c := make(chan Message)
	go func() {
		for i := 0; ; i++ {
			c <- Message{
				str:  fmt.Sprintf("%s: %d", msg, i),
				wait: waitForIt,
			}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}

func fanIn(input <-chan Message, input2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			select {
				case s := <-input:
					c <- s
				case s := <-input2:
					c <- s
			}
		}
	}()
	return c
}

func main() {
	c := fanIn(boring("joel"), boring("abril"))

	for i := 0; i < 10; i++ {
		msg1 := <-c; fmt.Println(msg1.str)
		msg2 := <-c; fmt.Println(msg2.str)
		msg1.wait <- true
		msg2.wait <- true
		//fmt.Printf("Somebody said: %q\n", <-c)
	}
	fmt.Println("You're both boring. I'm leaving.")
}
