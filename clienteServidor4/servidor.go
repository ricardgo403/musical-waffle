package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

var mostrarProcesos = true
var myProcessesList = list.New()
var myClientsList = list.New()

func f(n int, channel chan string) {
	var (
		//goroutine
		i   uint64
		max uint64
	)
	max = 18446744073709551615

	for i = 0; i < max; i++ {
		select {
		case msg1 := <-channel:
			if msg1 ==  strconv.Itoa(n) {
				channel <- string(i)
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

const MAXCONN int = 5

type Proceso struct {
	Id int
	Valor uint64
	Channel chan string
	ChannelExit chan string
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
		for {
			var proceso = myProcessesList.Front().Value.(Proceso) //proceso["channel"] <- proceso["id"]
			proceso.Channel <- strconv.Itoa(proceso.Id)
			reply := <-proceso.Channel
			fmt.Println("Respuesta: ", reply)
			entero, _ := strconv.Atoi(reply)
			if entero > 0 {
				proceso.Valor = uint64(entero)
				fmt.Println("Proceso: ", proceso)
				err := gob.NewEncoder(c).Encode(proceso)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Mensaje:", proceso)
				}
				myProcessesList.Remove(myProcessesList.Front())
				break
			} else {
				continue
			}
		}
	}
}

func main() {
	c := make(chan string)
	/*a := func (n int, channel chan string) {
		var (
			//goroutine
			i   uint64
			max uint64
		)
		max = 18446744073709551615

		for i = 0; i < max; i++ {
			select {
			case msg1 := <-channel:
				if msg1 == strconv.Itoa(n) {
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
	}*/
	/*var b = a
	go b(1, c)
	go a(2, c)*/

	for i := 0; i < 5; i++ {
		var aux = Proceso{
			Id:      i,
			Valor:   0,
			Channel: c,
		}
		myProcessesList.PushBack(aux)
		go f(aux.Id, aux.Channel)
	}
	/*for e := myProcessesList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}*/
	go servidor()
	var input string
	fmt.Scanln(&input)
}