package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"./process"
	"github.com/google/uuid"
)

var clientId = uuid.New()
var quitChannel = make(chan uint64)
var mostrarProcesos = true
var myReturnChannel = make(chan process.Process)
var myIdChannel = make(chan int)
var myProcess process.Process

func f() {
	var (
		//goroutine
		i   uint64
		max uint64
	)
	max = 18446744073709551615

	for i = myProcess.Value; i < max; i++ {
		select {
		case msg1 := <-myIdChannel:
			if msg1 == myProcess.Id {
				myProcess.Value = i
				myReturnChannel <- myProcess
				return
			} else {
				myIdChannel <- msg1
			}
		default:
			if mostrarProcesos {
				fmt.Println("Cliente", "id: ", myProcess.Id, ":", i)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func client() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Enviando: ", persona)
	//err = gob.NewEncoder(c).Encode(persona)
	err = gob.NewEncoder(c).Encode(clientId)

	err = gob.NewDecoder(c).Decode(&myProcess)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Mensaje recibido:", myProcess)
		go f()
	}
	c.Close()
	fmt.Println("Disconnected...")
	// return proceso
}

func sendProcess() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Enviando: ", persona)
	//err = gob.NewEncoder(c).Encode(persona)
	myIdChannel <- myProcess.Id
	thisProcess := <-myReturnChannel
	fmt.Println(thisProcess)
	err = gob.NewEncoder(c).Encode(clientId)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(myProcess)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Mensaje enviado:", myProcess)
		// go f(myProcess)
	}

	defer c.Close()
	fmt.Println("Disconnected...")
	// return proceso
}

func main() {
	go client()
	var input string
	fmt.Scanln(&input)
	sendProcess()
	fmt.Scanln(&input)
	// myIdChannel <- proceso.Id
	// reply := <-myReturnChannel
}
