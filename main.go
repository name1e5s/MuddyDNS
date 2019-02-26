package main

import (
	"flag"
	"github.com/name1e5s/socialismDNS/server"
	"log"
	"net"
)

func main() {
	var remote string
	var port int
	flag.StringVar(&remote, "r", "114.114.114.114", "forward DNS server address, default as 114.114.114.114")
	flag.IntVar(&port, "p", 53, "server port, default as 53")
	flag.Parse()

	// Listen to the default port
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	log.Println("Listening: " + listener.LocalAddr().String())

	data := make([]byte, 1024)
	// Forwarding the request to remote server...
	for {
		n, addr, err := listener.ReadFromUDP(data)
		question := server.GetQuestion(data[:n])
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Printf("Forward %s to %s", question.QNAMEToString(), remote)
		_, writeErr := listener.WriteToUDP(server.ForwardRequest(data[:n], remote), addr)
		if writeErr != nil {
			log.Printf("error: %s", writeErr)
		}

	}
}
