package socket

import (
	"fmt"
	"log"
	"net"
)

func HandleUDPConnection(conn *net.UDPConn, reqQueue chan string) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", string(buffer[:n]))

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	reqQueue <- string(buffer[:n])

	// write message back to client
	message := []byte("Connection with Server successful")
	_, err = conn.WriteToUDP(message, addr)

	if err != nil {
		log.Println(err.Error())
	}

}
