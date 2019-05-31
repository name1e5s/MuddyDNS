package server

import (
	log "github.com/sirupsen/logrus"
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
		log.Error("Connection failed!")
		return nil // Failed
	}

	// Error? What error?
	_ = socket.SetDeadline(time.Now().Add(time.Duration(time.Second * 2)))

	defer func() {
		err := socket.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = socket.Write(data)
	if err != nil {
		log.Error("Send data failed!")
		return nil // Failed
	}

	receive := make([]byte, 4*1024)

	num, addr, err := socket.ReadFromUDP(receive)
	if err != nil || num < 0 {
		log.Debug("Read data from ", addr, " failed: ", err)
		return nil
	}
	return receive[:num]
}
