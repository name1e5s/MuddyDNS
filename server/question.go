package server

import (
	"bytes"
	"strings"
)

// DNS Question section
//                                    1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                                               |
//    /                     QNAME                     /
//    /                                               /
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                     QTYPE                     |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                     QCLASS                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Question struct {
	QNAME  []byte // Domain Address
	QTYPE  []byte // Question Type(won't use in our server)
	QCLASS []byte // Question Class.(won't use)
}

func GetQuestion(data []byte) Question {
	// Here the data contains a Header.
	length := 12
	for data[length] != 0 {
		length++
	}
	return Question{
		QNAME:  data[12 : length+1],
		QTYPE:  data[length+1 : length+3],
		QCLASS: data[length+3 : length+5],
	}
}

func (question Question) toBytes() []byte {
	// Here the output doesn't contain a Header.
	buffer := bytes.NewBuffer(question.QNAME)
	buffer.Write(question.QTYPE)
	buffer.Write(question.QCLASS)
	return buffer.Bytes()
}

// each label consists of a length octet followed by that
// number of octets. The domain name terminates with the
// zero length octet for the null label of the root.  Note
// QNAMEToString converts the QNAME to a string.
//
// QNAME is a domain name represented as a sequence of labels, where
// that this field may be an odd number of octets; no
// padding is used.
func (question Question) QNAMEToString() string {
	parts := make([]string, 0)
	for i := 0; i < len(question.QNAME); {
		length := int(question.QNAME[i])

		if length == 0 {
			break
		}

		offset := i + 1
		parts = append(parts, string(question.QNAME[offset:offset+length]))
		i = offset + length
	}
	return strings.Join(parts, ".")
}

