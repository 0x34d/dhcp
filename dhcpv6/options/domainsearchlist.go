package options

// This module defines the OptDomainSearchList structure.
// https://www.ietf.org/rfc/rfc3646.txt

import (
	"encoding/binary"
	"fmt"
)

type OptDomainSearchList struct {
	domainSearchList []string
}

func (op *OptDomainSearchList) Code() OptionCode {
	return DOMAIN_SEARCH_LIST
}

func (op *OptDomainSearchList) ToBytes() []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint16(buf[0:2], uint16(DOMAIN_SEARCH_LIST))
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.Length()))
	buf = append(buf, LabelsToBytes(op.domainSearchList)...)
	return buf
}

func (op *OptDomainSearchList) DomainSearchList() []string {
	return op.domainSearchList
}

func (op *OptDomainSearchList) SetDomainSearchList(dsList []string) {
	op.domainSearchList = dsList
}

func (op *OptDomainSearchList) Length() int {
	var length int
	for _, label := range op.domainSearchList {
		length += len(label) + 2 // add the first and the last length bytes
	}
	return length
}

func (op *OptDomainSearchList) String() string {
	return fmt.Sprintf("OptDomainSearchList{searchlist=%v}", op.domainSearchList)
}

// build an OptDomainSearchList structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptDomainSearchList(data []byte) (*OptDomainSearchList, error) {
	opt := OptDomainSearchList{}
	var err error
	opt.domainSearchList, err = LabelsFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}
