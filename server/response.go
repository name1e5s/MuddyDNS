package server

import (
	"bytes"
	"encoding/binary"
)

// RR type codes
type Type uint16

const (
	TypeA     Type = 1  // a host address
	TypeNS    Type = 2  // an authoritative name server
	TypeMD    Type = 3  // a mail destination (Obsolete - use MX)
	TypeMF    Type = 4  // a mail forwarder (Obsolete - use MX)
	TypeCNAME Type = 5  // the canonical name for an alias
	TypeSOA   Type = 6  // marks the start of a zone of authority
	TypeMB    Type = 7  // a mailbox domain name (EXPERIMENTAL)
	TypeMG    Type = 8  // a mail group member (EXPERIMENTAL)
	TypeMR    Type = 9  // a mail rename domain name (EXPERIMENTAL)
	TypeNULL  Type = 10 // a null RR (EXPERIMENTAL)
	TypeWKS   Type = 11 // a well known service description
	TypePTR   Type = 12 // a domain name pointer
	TypeHINFO Type = 13 // host information
	TypeMINFO Type = 14 // mailbox or mail list information
	TypeMX    Type = 15 // mail exchange
	TypeTXT   Type = 16 // text strings
)

func (typ Type) String() string {
	switch typ {
	case TypeA:
		return "A"
	case TypeNS:
		return "NS"
	case TypeMD:
		return "MD"
	case TypeMF:
		return "MF"
	case TypeCNAME:
		return "CNAME"
	case TypeSOA:
		return "SOA"
	case TypeMB:
		return "MB"
	case TypeMG:
		return "MG"
	case TypeMR:
		return "MR"
	case TypeNULL:
		return "NULL"
	case TypeWKS:
		return "WKS"
	case TypePTR:
		return "PTR"
	case TypeHINFO:
		return "HINFO"
	case TypeMINFO:
		return "MINFO"
	case TypeMX:
		return "MX"
	case TypeTXT:
		return "TXT"
	default:
		return "Unknown"
	}
}

// Class defines the class associated with a request/response.  Different DNS
// classes can be thought of as an array of parallel namespace trees.
type Class uint16

const (
	ClassIN  Class = 1   // Internet
	ClassCS  Class = 2   // the CSNET class (Obsolete used only for examples in some obsolete RFCs)
	ClassCH  Class = 3   // the CHAOS class
	ClassHS  Class = 4   // Hesiod [Dyer 87]
	ClassAny Class = 255 // any class
)

func (class Class) String() string {
	switch class {
	case ClassIN:
		return "IN"
	case ClassCS:
		return "CS"
	case ClassCH:
		return "CH"
	case ClassHS:
		return "HS"
	case ClassAny:
		return "Any"
	default:
		return "Unknown"
	}
}

// Resource record
//                                    1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                                               |
//    /                                               /
//    /                      NAME                     /
//    |                                               |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      TYPE                     |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                     CLASS                     |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      TTL                      |
//    |                                               |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                   RDLENGTH                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//    /                     RDATA                     /
//    /                                               /
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type ResourceRecord struct {
	NAME     []byte
	TYPE     Type
	CLASS    Class
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func (rr ResourceRecord) toBytes() []byte {
	buffer := bytes.NewBuffer(rr.NAME)
	_ = binary.Write(buffer, binary.BigEndian, rr.TYPE)
	_ = binary.Write(buffer, binary.BigEndian, rr.CLASS)
	_ = binary.Write(buffer, binary.BigEndian, rr.TTL)
	_ = binary.Write(buffer, binary.BigEndian, rr.RDLENGTH)
	buffer.Write(rr.RDATA)

	return buffer.Bytes()
}

type Response struct {
	HEADER   Header
	QUESTION Question
	RR       ResourceRecord
}

func (response Response) toBytes() []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(response.HEADER.toBytes())
	buffer.Write(response.QUESTION.toBytes())
	buffer.Write(response.RR.toBytes())
	return buffer.Bytes()
}
