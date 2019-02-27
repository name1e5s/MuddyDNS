package server

import (
	"log"
	"net"
	"time"
)

func ForwardRequest(data []byte, server string) []byte {
	socket, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   net.ParseIP(server),
			Port: 53, // Default DNS port
		})

	if err != nil {
		log.Println("Connection failed!")
		return nil // Failed
	}

	// Error? What error?
	_ = socket.SetDeadline(time.Now().Add(time.Duration(time.Millisecond * 300)))
	defer socket.Close()

	_, err = socket.Write(data)
	if err != nil {
		log.Println("Send data failed!")
		return nil // Failed
	}

	receive := make([]byte, 4*1024)

	num, addr, err := socket.ReadFromUDP(receive)
	if err != nil {
		log.Print("Read data from ", addr, " failed: ",err)
		return nil // Failed
	}
	return receive[:num]
}
