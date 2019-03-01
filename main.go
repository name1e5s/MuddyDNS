package main

import (
	"flag"
	"github.com/name1e5s/MuddyDNS/server"
	"log"
	"net"
)

type receivedData struct {
	addr *net.UDPAddr
	data []byte
}

var remote string
var port int
var harmonyPath string
var harmonyList server.DNSList

func main() {
	flag.StringVar(&remote, "r", "10.3.9.5", "forward DNS server address, default as 10.3.9.5")
	flag.IntVar(&port, "p", 53, "server port, default as 53")
	flag.StringVar(&harmonyPath, "f", "./example", "harmony file path")
	flag.Parse()

	// Show Logo
	log.Println("░█▄█░█░█░█▀▄░█▀▄░█░█░█▀▄░█▀█░█▀▀")
	log.Println("░█░█░█░█░█░█░█░█░░█░░█░█░█░█░▀▀█")
	log.Println("░▀░▀░▀▀▀░▀▀░░▀▀░░░▀░░▀▀░░▀░▀░▀▀▀")

	// Loading file
	harmonyList = server.LoadConfig(harmonyPath)
	log.Println(harmonyList)

	// Listen to the default port
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	log.Println("Listening: " + listener.LocalAddr().String())
	received := readUDP(listener)
	// Forwarding the request to remote server...
	for {
		go writeUDP(listener, <-received)
	}
}

func readUDP(conn *net.UDPConn) chan receivedData {
	ch := make(chan receivedData)
	data := make([]byte, 1024)
	go func() {
		for {
			n, addr, err := conn.ReadFromUDP(data)
			if err != nil {
				log.Printf("error: %s", err)
			}
			ch <- receivedData{addr, data[:n]}
		}
	}()
	return ch
}

func writeUDP(conn *net.UDPConn, data receivedData) {
	_, writeErr := conn.WriteToUDP(server.LocalResolv(data.data, remote, harmonyList), data.addr)
	if writeErr != nil {
		log.Printf("error: %s", writeErr)
	}
}
