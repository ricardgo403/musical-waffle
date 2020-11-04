package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

var quitChannel = make(chan uint64)
var mostrarProcesos = true

type Persona struct {
	Nombre string
	Email []string
}

type Proceso struct {
	Id int
	Valor uint64
	Channel chan string
}

func f2(n int, valor uint64, channel chan string, returnChannel chan uint64) {
	var (
		//goroutine
		i   uint64
		max uint64
	)
	max = 18446744073709551615

	for i = valor; i < max; i++ {
		select {
		case msg1 := <-channel:
			if msg1 == strconv.Itoa(n) {
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



func cliente(){
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Enviando: ", persona)
	//err = gob.NewEncoder(c).Encode(persona)
	var proceso Proceso
	err = gob.NewDecoder(c).Decode(&proceso)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Mensaje recibido:", proceso)
		go f2(proceso.Id, proceso.Valor, proceso.Channel, quitChannel)
	}
	c.Close()
}

func main() {
	/*persona := Persona{
		"Joel",
		[]string{
		"joel.gv@gmail.com",
		"joel.gv@udg.mx",
		},
	}*/
	go cliente()
	var input string
	fmt.Scanln(&input)
}