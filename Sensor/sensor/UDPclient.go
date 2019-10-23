package sensor

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func (s *Sensor) ClientSensorUPDConnection() {
	hostName := serverIPaddress
	portNum := serverPort

	fmt.Printf("Hostname: %s:%s \n", hostName, strconv.Itoa(s.port))

	service := hostName + ":" + strconv.Itoa(portNum)

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	if err != nil {
		log.Println(err)
		return
	}

	localAddr, err := net.ResolveUDPAddr("udp", s.ipaddress+":"+strconv.Itoa(s.port))

	//LocalAddr := nil
	conn, err := net.DialUDP("udp", localAddr, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)
	}

	/*log.Printf("Established connection to %s \n", service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())*/

	defer conn.Close()

	// write a message to server
	fmt.Printf("Wert: %s \n", s.value)
	message := []byte(string(s.sType) + ":" + s.value)

	_, err = conn.Write(message)

	if err != nil {
		log.Println(err)
	}

	// receive message from server
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP Server : ", addr)
	fmt.Println("Received from UDP server : ", string(buffer[:n]))
}
