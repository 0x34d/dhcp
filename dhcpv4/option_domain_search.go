package dhcpv4

// This module defines the OptDomainSearch structure.
// https://tools.ietf.org/html/rfc3397

import (
	"fmt"
)

// OptDomainSearch represents an option encapsulating a domain search list.
type OptDomainSearch struct {
	DomainSearch []string
}

func (op *OptDomainSearch) Code() OptionCode {
	return OptionDNSDomainSearchList
}

func (op *OptDomainSearch) ToBytes() []byte {
	buf := []byte{byte(op.Code()), byte(op.Length())}
	buf = append(buf, LabelsToBytes(op.DomainSearch)...)
	return buf
}

func (op *OptDomainSearch) Length() int {
	var length int
	for _, label := range op.DomainSearch {
		length += len(label) + 2 // add the first and the last length bytes
	}
	return length
}

func (op *OptDomainSearch) String() string {
	return fmt.Sprintf("DNS Domain Search List ->", op.DomainSearch)
}

// build an OptDomainSearch structure from a sequence of bytes.
func ParseOptDomainSearch(data []byte) (*OptDomainSearch, error) {
	if len(data) < 2 {
		return nil, ErrShortByteStream
	}
	code := OptionCode(data[0])
	if code != OptionDNSDomainSearchList {
		return nil, fmt.Errorf("expected code %v, got %v", OptionDNSDomainSearchList, code)
	}
	length := int(data[1])
	if len(data) < 2+length {
		return nil, ErrShortByteStream
	}
	domainSearch, err := LabelsFromBytes(data[2:length+2])
	if err != nil {
		return nil, err
	}
	return &OptDomainSearch{DomainSearch: domainSearch}, nil
}
