package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"./process"
	"github.com/google/uuid"
)

var mostrarProcesos = true
var myProcessesList = list.New()
var myClientsList = list.New()
var myReturnChannel = make(chan process.Process)
var myIdChannel = make(chan int)

const MAXCONN int = 5

func f(myProcess process.Process) {
	var (
		//goroutine
		i   uint64
		max uint64
	)
	max = 18446744073709551615
	// fmt.Println(myProcess)
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
				fmt.Println("id: ", myProcess.Id, ":", i)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func isClientConnected(id uuid.UUID) bool {
	isConnected := false
	for e := myClientsList.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		if e.Value == id {
			isConnected = true
		}
	}
	return isConnected
}

func server() {
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
				var clientId uuid.UUID
				err = gob.NewDecoder(c).Decode(&clientId)
				fmt.Println("Received a client id:", clientId)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if isClientConnected(clientId) {
					// fmt.Println("This client id is already known:", clientId)
					// var aux process.Process
					// err := gob.NewDecoder(c).Decode(&aux)
					// fmt.Println(aux)
					// if err != nil {
					// 	fmt.Println("Error**\n", err)
					// 	continue
					// } else {
					// 	fmt.Println("Mensaje recibido:", aux)
					// 	myProcessesList.PushBack(aux)
					// 	go f(aux)

					// 	// c.Close()
					// 	fmt.Println("Disconnected...")
					// }
					retrieveProcess2(c, clientId)
				} else if myClientsList.Len() < MAXCONN {
					fmt.Println("This client id", clientId, "is not registered...")
					myClientsList.PushBack(clientId)
					go handleClient(c)
					fmt.Println("Finish handling client")
				}
			}
			// c.Close()
		}
	}
}

func handleClient(c net.Conn) {
	//err := gob.NewDecoder(c).Decode(&proceso)
	if myProcessesList.Len() > 0 {
		var proceso = myProcessesList.Front().Value.(process.Process) //proceso["channel"] <- proceso["id"]
		myIdChannel <- proceso.Id
		reply := <-myReturnChannel
		// fmt.Println("Respuesta: ", reply)
		// if reply > 0 {
		// 	proceso.Valor = reply
		// 	fmt.Println("Proceso: ", proceso)
		// 	err := gob.NewEncoder(c).Encode(proceso)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	} else {
		// 		fmt.Println("Mensaje:", proceso)
		// 	}
		// 	myProcessesList.Remove(myProcessesList.Front())
		// }
		nuevoProceso := reply
		fmt.Println("Proceso: ", proceso)
		fmt.Println("Proceso nuevo: ", nuevoProceso)
		err := gob.NewEncoder(c).Encode(nuevoProceso)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Mensaje:", nuevoProceso)
		}
		myProcessesList.Remove(myProcessesList.Front())
	}
}

func retrieveProcess(c net.Conn) {
	var aux process.Process
	err := gob.NewDecoder(c).Decode(&aux)
	fmt.Println(aux)
	if err != nil {
		fmt.Println("Error**\n", err)
		return
	}

	fmt.Println("Mensaje recibido:", aux)
	myProcessesList.PushBack(aux)
	go f(aux)

	c.Close()
	fmt.Println("Disconnected...")
}

func deleteClient(clientId uuid.UUID) bool {
	wasDeleted := false
	for e := myClientsList.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		if e.Value == clientId {
			myClientsList.Remove(e)
			wasDeleted = true
			fmt.Println("Was deleted:", wasDeleted)
		}
	}
	return wasDeleted
}

func retrieveProcess2(c net.Conn, clientId uuid.UUID) {
	fmt.Println("This client id is already known:", clientId)
	var aux process.Process
	err := gob.NewDecoder(c).Decode(&aux)
	fmt.Println(aux)
	if err != nil {
		fmt.Println("Error**\n", err)
	} else {
		fmt.Println("Mensaje recibido:", aux)
		myProcessesList.PushBack(aux)
		fmt.Println("Len Processes:", myProcessesList.Len())
		deleteClient(clientId)
		go f(aux)

		// c.Close()
		fmt.Println("Disconnected...")
	}
	// err = gob.NewDecoder(c).Decode(&aux)
	// fmt.Println(aux)
	// if err != nil {
	// 	fmt.Println("Error**\n", err)
	// 	return
	// } else {
	// 	fmt.Println("Mensaje recibido:", aux)
	// 	myProcessesList.PushBack(aux)
	// 	go f(aux)

	// 	// c.Close()
	// 	fmt.Println("Disconnected...")
	// }
}

func main() {
	for i := 0; i < 5; i++ {
		var aux = process.Process{
			Id:    i,
			Value: 0,
		}
		myProcessesList.PushBack(aux)
		go f(aux)
	}

	go server()
	var input string
	fmt.Scanln(&input)
}
