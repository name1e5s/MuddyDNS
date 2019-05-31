package main

import (
	"flag"
	"github.com/name1e5s/MuddyDNS/server"
	"github.com/name1e5s/MuddyDNS/utils"
	log "github.com/sirupsen/logrus"
	"net"
)

type receivedData struct {
	addr *net.UDPAddr
	data []byte
}

var remote string
var port int
var harmonyPath string
var debug bool
var debugDetailed bool
var harmonyList server.DNSList

func main() {
	flag.StringVar(&remote, "r", "10.3.9.5", "forward DNS server address, default as 10.3.9.5")
	flag.IntVar(&port, "p", 53, "server port, default as 53")
	flag.StringVar(&harmonyPath, "f", "./example", "harmony file path")
	flag.BoolVar(&debug, "d", false, "print debug info")
	flag.BoolVar(&debugDetailed, "dd", false, "print (more detailed) debug info")
	flag.Parse()

	// Set logger info
	log.SetFormatter(&utils.Formatter{})
	if debugDetailed {
		log.SetLevel(log.TraceLevel)
	} else if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Show Logo
	log.Info("░█▄█░█░█░█▀▄░█▀▄░█░█░█▀▄░█▀█░█▀▀")
	log.Info("░█░█░█░█░█░█░█░█░░█░░█░█░█░█░▀▀█")
	log.Info("░▀░▀░▀▀▀░▀▀░░▀▀░░░▀░░▀▀░░▀░▀░▀▀▀")

	// Loading file
	harmonyList = server.LoadConfig(harmonyPath)
	log.Println(harmonyList)

	// Listen to the default port
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Listening: " + listener.LocalAddr().String())
	received := readUDP(listener)
	defer close(received)
	// Forwarding the request to remote server...
	for {
		go writeUDP(listener, <-received)
	}
}

func readUDP(conn *net.UDPConn) chan receivedData {
	ch := make(chan receivedData)
	go func() {
		for {
			data := make([]byte, 65536)
			n, addr, err := conn.ReadFromUDP(data)
			if err != nil {
				log.Debugf("error: %s", err)
			}
			qry := receivedData{addr, data[:n]}
			ch <- qry
		}
	}()
	return ch
}

func writeUDP(conn *net.UDPConn, data receivedData) {
	_, writeErr := conn.WriteToUDP(server.LocalResolv(data.addr, data.data, remote, harmonyList), data.addr)
	if writeErr != nil {
		log.Debugf("error: %s", writeErr)
	}
}
