// server is a simple DNS server implementation.
//
// For more details, go to:
// https://www.ietf.org/rfc/rfc1035.txt
package server

import (
	"github.com/name1e5s/socialismDNS/utils"
)

// OpCode defines some different operations types
type OpCode uint8 // Why there is no uint4????

const (
	OpcodeQuery  OpCode = 0 // Query
	OpcodeIQuery OpCode = 1 // Inverse Query
	OpcodeStatus OpCode = 2 // Statue Request
	// Reversed for the future.
)

func (opcode OpCode) String() string {
	switch opcode {
	case OpcodeQuery:
		return "QUERY"
	case OpcodeIQuery:
		return "INVERSE QUERY"
	case OpcodeStatus:
		return "STATUS"
	default:
		return "UNKNOWN"
	}
}

// ResponseCode provides response codes for question answers.
type ResponseCode uint8

const (
	ResponseCodeNoErr    ResponseCode = 0 // No error
	ResponseCodeFormErr  ResponseCode = 1 // Format Error
	ResponseCodeServFail ResponseCode = 2 // Server Failure
	ResponseCodeNXDomain ResponseCode = 3 // Non-Existent Domain
	ResponseCodeNotImp   ResponseCode = 4 // Not Implemented
	ResponseCodeRefused  ResponseCode = 5 // Query Refused
)

func (res ResponseCode) String() string {
	switch res {
	case ResponseCodeNoErr:
		return "No Error"
	case ResponseCodeFormErr:
		return "Format Error"
	case ResponseCodeServFail:
		return "Server Failure"
	case ResponseCodeNXDomain:
		return "Non-Existent Domain"
	case ResponseCodeNotImp:
		return "Not Implemented"
	case ResponseCodeRefused:
		return "Query Refused"
	default:
		return "Unknown"
	}
}

//  DNS Header
//  0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                      ID                       |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    QDCOUNT                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    ANCOUNT                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    NSCOUNT                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//  |                    ARCOUNT                    |
//  +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Header struct {
	ID     uint16
	QR     bool
	Opcode OpCode

	AA           bool  // Authoritative answer
	TC           bool  // Truncated
	RD           bool  // Recursion desired
	RA           bool  // Recursion available
	Z            uint8 // Reserved
	ResponseCode ResponseCode

	QDCOUNT uint16 // number of entries in the question section.
	ANCOUNT uint16 // number of resource records in the answer section.
	NSCOUNT uint16 // number of name server resource records in the authority records section
	ARCOUNT uint16 // number of resource records in the additional records section
}

func GetHeader(data []byte) Header {
	return Header{
		ID:           utils.BytesToUInt16(data[0:2]),
		QR:           data[2]&0x80 != 0,
		Opcode:       OpCode(data[2]>>3) & 0xF,
		AA:           data[2]&0x04 != 0,
		TC:           data[2]&0x02 != 0,
		RD:           data[2]&0x01 != 0,
		RA:           data[3]&0x80 != 0,
		Z:            uint8(data[3]>>4) & 0x7,
		ResponseCode: ResponseCode(data[3] & 0xF),
		QDCOUNT:      utils.BytesToUInt16(data[4:6]),
		ANCOUNT:      utils.BytesToUInt16(data[6:8]),
		NSCOUNT:      utils.BytesToUInt16(data[8:10]),
		ARCOUNT:      utils.BytesToUInt16(data[10:12]),
	}
}

func (header Header) toBytes() []byte {
	buffer := make([]byte, 12)
	buffer[0], buffer[1] = utils.UInt16ToBytes(header.ID)
	buffer[2] = byte(utils.BoolToInt(header.QR)<<7 |
		int(header.Opcode)<<3 | utils.BoolToInt(header.AA)<<2 |
		utils.BoolToInt(header.TC)<<1 | utils.BoolToInt(header.RD))
	buffer[3] = byte(utils.BoolToInt(header.RA)<<7 | int(header.Z)<<4 |
		int(header.ResponseCode&0xF))
	buffer[4], buffer[5] = utils.UInt16ToBytes(header.QDCOUNT)
	buffer[6], buffer[7] = utils.UInt16ToBytes(header.ANCOUNT)
	buffer[8], buffer[9] = utils.UInt16ToBytes(header.NSCOUNT)
	buffer[10], buffer[11] = utils.UInt16ToBytes(header.ARCOUNT)
	return buffer
}
