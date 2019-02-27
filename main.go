package main

import (
	"flag"
	"github.com/name1e5s/MuddyDNS/server"
	"log"
	"net"
)

func main() {
	var remote string
	var port int
	var harmonyPath string
	flag.StringVar(&remote, "r", "114.114.114.114", "forward DNS server address, default as 114.114.114.114")
	flag.IntVar(&port, "p", 53, "server port, default as 53")
	flag.StringVar(&harmonyPath, "f", "./example", "harmony file path")
	flag.Parse()

	// Show Logo
	log.Println("░█▄█░█░█░█▀▄░█▀▄░█░█░█▀▄░█▀█░█▀▀")
	log.Println("░█░█░█░█░█░█░█░█░░█░░█░█░█░█░▀▀█")
	log.Println("░▀░▀░▀▀▀░▀▀░░▀▀░░░▀░░▀▀░░▀░▀░▀▀▀")

	// Loading file
	harmonyList := server.LoadConfig(harmonyPath)
	log.Println(harmonyList)

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
		if err != nil {
			log.Printf("error: %s", err)
		}
		_, writeErr := listener.WriteToUDP(server.LocalResolv(data[:n], remote, harmonyList), addr)
		if writeErr != nil {
			log.Printf("error: %s", writeErr)
		}

	}
}
