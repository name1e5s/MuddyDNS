package server

import (
	"bufio"
	"github.com/name1e5s/MuddyDNS/utils"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

type DNSList map[string]string

// LoadConfig loads the hosts-like config file in to a map.
func LoadConfig(path string) (harmonyList DNSList) {
	log.Trace("Loading config file: ", path)
	harmonyList = make(DNSList)
	harmonyFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open Config File Error: %s", err)
		return nil
	}

	defer func() {
		err := harmonyFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fileScanner := bufio.NewScanner(harmonyFile)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		strarr := strings.Fields(line)
		if len(strarr) == 2 {
			if utils.IsDomainName(strarr[1]) && utils.IsIP(strarr[0]) {
				harmonyList[strings.ToLower(strarr[1])] = strarr[0]
			} else {
				log.Warnf("Wrong config format: %s", line)
			}
		}
	}
	log.Trace("Load config file ", path, " done.")
	return
}

// LocalResolv Resolves the domain name and protect the user from some
// "illegal" websites.
func LocalResolv(addr *net.UDPAddr, rawdata []byte, remote string, harmonyList DNSList) []byte {
	header := GetHeader(rawdata)
	question := GetQuestion(rawdata)
	qnameStr := strings.ToLower(question.QNAMEToString())
	log.Debugf("Get new package with ID 0x%04x", header.ID)
	log.Trace("Get new query package for " + qnameStr + " from " + addr.String())
	if utils.BytesToUInt16(question.QTYPE) != uint16(1) || harmonyList == nil || harmonyList[qnameStr] == "" {
		log.Debugf("Forward request for %s to %s.", qnameStr, remote)
		return ForwardRequest(rawdata, remote) // A "harmonic" domain name.
	}
	if harmonyList[qnameStr] == "0.0.0.0" {
		log.Warn("Illegal domain name: ", qnameStr)
		header := Header{
			ID:           header.ID,
			QR:           true,
			Opcode:       OpcodeQuery,
			AA:           false,
			TC:           false,
			RD:           false,
			RA:           false,
			Z:            0,
			ResponseCode: ResponseCodeNXDomain,
			QDCOUNT:      1,
			ANCOUNT:      0,
			NSCOUNT:      0,
			ARCOUNT:      0,
		}

		return append(header.toBytes(), question.toBytes()...)
	} else {
		log.Debugf("Redirect request for %s to %s.", qnameStr, harmonyList[qnameStr])
		response := Response{
			HEADER: Header{
				ID:           header.ID,
				QR:           true,
				Opcode:       OpcodeQuery,
				AA:           false,
				TC:           false,
				RD:           false,
				RA:           false,
				Z:            0,
				ResponseCode: ResponseCodeNoErr,
				QDCOUNT:      1,
				ANCOUNT:      1,
				NSCOUNT:      0,
				ARCOUNT:      0,
			},
			QUESTION: question,
			RR: ResourceRecord{
				NAME: []byte{0xc0, 0x0c}, // Points to QNAME in Question Section.
				// Ref: https://tools.ietf.org/html/rfc1035#section-4.1.4
				TYPE:     TypeA,
				CLASS:    ClassIN,
				TTL:      114514,
				RDLENGTH: 4,
				RDATA:    utils.IpToBytes(harmonyList[qnameStr]),
			},
		}

		return response.toBytes()
	}
}
