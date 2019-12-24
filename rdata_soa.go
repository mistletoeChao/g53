package g53

import (
	"errors"
	"strings"

	"github.com/mistletoeChao/g53/util"
)

type SOA struct {
	MName   *Name
	RName   *Name
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
}

func (soa *SOA) Rend(r *MsgRender) {
	rendField(RDF_C_NAME, soa.MName, r)
	rendField(RDF_C_NAME, soa.RName, r)
	rendField(RDF_C_UINT32, soa.Serial, r)
	rendField(RDF_C_UINT32, soa.Refresh, r)
	rendField(RDF_C_UINT32, soa.Retry, r)
	rendField(RDF_C_UINT32, soa.Expire, r)
	rendField(RDF_C_UINT32, soa.Minimum, r)
}

func (soa *SOA) ToWire(buffer *util.OutputBuffer) {
	fieldToWire(RDF_C_NAME, soa.MName, buffer)
	fieldToWire(RDF_C_NAME, soa.RName, buffer)
	fieldToWire(RDF_C_UINT32, soa.Serial, buffer)
	fieldToWire(RDF_C_UINT32, soa.Refresh, buffer)
	fieldToWire(RDF_C_UINT32, soa.Retry, buffer)
	fieldToWire(RDF_C_UINT32, soa.Expire, buffer)
	fieldToWire(RDF_C_UINT32, soa.Minimum, buffer)
}

func (soa *SOA) String() string {
	var ss []string
	ss = append(ss, fieldToStr(RDF_D_NAME, soa.MName))
	ss = append(ss, fieldToStr(RDF_D_NAME, soa.RName))
	ss = append(ss, fieldToStr(RDF_D_INT, soa.Serial))
	ss = append(ss, fieldToStr(RDF_D_INT, soa.Refresh))
	ss = append(ss, fieldToStr(RDF_D_INT, soa.Retry))
	ss = append(ss, fieldToStr(RDF_D_INT, soa.Expire))
	ss = append(ss, fieldToStr(RDF_D_INT, soa.Minimum))
	return strings.Join(ss, " ")
}

func SOAFromWire(buffer *util.InputBuffer, ll uint16) (*SOA, error) {
	name, ll, err := fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}
	mname, _ := name.(*Name)

	name, ll, err = fieldFromWire(RDF_C_NAME, buffer, ll)
	if err != nil {
		return nil, err
	}
	rname, _ := name.(*Name)

	i, ll, err := fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}
	serial, _ := i.(uint32)

	i, ll, err = fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}
	refresh, _ := i.(uint32)

	i, ll, err = fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}
	retry, _ := i.(uint32)

	i, ll, err = fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}
	expire, _ := i.(uint32)

	i, ll, err = fieldFromWire(RDF_C_UINT32, buffer, ll)
	if err != nil {
		return nil, err
	}
	minimum, _ := i.(uint32)

	if ll != 0 {
		return nil, errors.New("extra data in rdata part")
	}

	return &SOA{mname, rname, serial, refresh, retry, expire, minimum}, nil
}

func SOAFromString(s string) (*SOA, error) {
	fields := strings.Split(s, " ")
	if len(fields) != 6 {
		return nil, errors.New("short of fields for soa")
	}

	name, err := fieldFromStr(RDF_D_NAME, fields[0])
	if err != nil {
		return nil, err
	}
	mname, _ := name.(*Name)

	name, err = fieldFromStr(RDF_D_NAME, fields[1])
	if err != nil {
		return nil, err
	}
	rname, _ := name.(*Name)

	i, err := fieldFromStr(RDF_D_INT, fields[2])
	if err != nil {
		return nil, err
	}
	serial, _ := i.(int)

	i, err = fieldFromStr(RDF_D_INT, fields[3])
	if err != nil {
		return nil, err
	}
	refresh, _ := i.(int)

	i, err = fieldFromStr(RDF_D_INT, fields[4])
	if err != nil {
		return nil, err
	}
	retry, _ := i.(int)

	i, err = fieldFromStr(RDF_D_INT, fields[4])
	if err != nil {
		return nil, err
	}
	expire, _ := i.(int)

	i, err = fieldFromStr(RDF_D_INT, fields[5])
	if err != nil {
		return nil, err
	}
	minimum, _ := i.(int)

	return &SOA{mname, rname, uint32(serial), uint32(refresh), uint32(retry), uint32(expire), uint32(minimum)}, nil
}
