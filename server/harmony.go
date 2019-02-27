package server

import (
	"bufio"
	"github.com/name1e5s/MuddyDNS/utils"
	"log"
	"os"
	"strings"
)

type DNSList map[string]string

// LoadConfig loads the hosts-like config file in to a map.
func LoadConfig(path string) (harmonyList DNSList) {
	harmonyList = make(DNSList)
	harmonyFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open Config File Error: %s", err)
		return nil
	}
	defer harmonyFile.Close()
	fileScanner := bufio.NewScanner(harmonyFile)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		strarr := strings.Split(line, " ")
		if len(strarr) == 2 {
			if utils.IsDomainName(strarr[1]) && utils.IsIP(strarr[0]) {
				harmonyList[strarr[1]] = strarr[0]
			}
		}
	}
	return
}

// LocalResolv Resolves the domain name and protect the user from some
// "illegal" websites.
func LocalResolv(rawdata []byte, remote string, harmonyList DNSList) []byte {
	header := GetHeader(rawdata)
	question := GetQuestion(rawdata)
	qnameStr := question.QNAMEToString()
	if harmonyList == nil || harmonyList[qnameStr] == "" {
		return ForwardRequest(rawdata, remote) // A "harmonic" domain name.
	}

	if harmonyList[qnameStr] == "0.0.0.0" {
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
