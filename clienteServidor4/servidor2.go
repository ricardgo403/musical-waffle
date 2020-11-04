package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var mostrarProcesos = true
var myProcessesList = list.New()
var myClientsList = list.New()
const MAXCONN int = 5

func f(n int, channel chan int, returnChannel chan uint64) {
	var (
		//goroutine
		i   uint64
		max uint64
	)
	max = 18446744073709551615

	for i = 0; i < max; i++ {
		select {
		case msg1 := <-channel:
			if msg1 == n {
				returnChannel <- i
				return
			} else {
				channel	<- msg1
			}
		default:
			if mostrarProcesos {
				fmt.Println("Channel: ", channel, " id: ", n, ":", i)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}



type Proceso struct {
	Id int
	Valor uint64
	Channel chan int
	returnChannel chan uint64
}

func servidor(){
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		for {
			c, err := s.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				// Create a new list and put some numbers in it.
				if myClientsList.Len() < MAXCONN {
					myClientsList.PushBack(c)
					go handleCliente(c)
				}
			}
		}
	}
}

func handleCliente(c net.Conn){
	//err := gob.NewDecoder(c).Decode(&proceso)
	if myProcessesList.Len() > 0 {
		var proceso = myProcessesList.Front().Value.(Proceso) //proceso["channel"] <- proceso["id"]
		proceso.Channel <- proceso.Id
		reply := <-proceso.returnChannel
		fmt.Println("Respuesta: ", reply)
		if reply > 0 {
			proceso.Valor = reply
			fmt.Println("Proceso: ", proceso)
			err := gob.NewEncoder(c).Encode(proceso)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Mensaje:", proceso)
			}
			myProcessesList.Remove(myProcessesList.Front())
		}
	}
}

func main() {
	c := make(chan int)
	returnChannel := make(chan uint64)

	for i := 0; i < 5; i++ {
		var aux = Proceso{
			Id:      i,
			Valor:   0,
			Channel: c,
			returnChannel: returnChannel,
		}
		myProcessesList.PushBack(aux)
		go f(aux.Id, aux.Channel, aux.returnChannel)
	}

	go servidor()
	var input string
	fmt.Scanln(&input)
}