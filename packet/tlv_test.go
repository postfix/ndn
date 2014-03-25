package packet

import (
	"bytes"
	//"fmt"
	"testing"
)

func TestReadByte(t *testing.T) {
	buf := bytes.NewReader([]byte{0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
	r, o, _ := ReadByte(buf)
	if o != 8 {
		t.Error("not reading the right length")
	}
	if r != 4822678189205111 {
		t.Error("not reading the right value", r)
	}
}

func TestWriteByte(t *testing.T) {
	buf := new(bytes.Buffer)
	WriteByte(buf, uint64(4822678189205111))
	if !EqualBytes(buf.Bytes(), []byte{0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}) {
		t.Error("not writing the right bytes")
	}
}

func TestDecode(t *testing.T) {
	v := new(TLV)
	r, _ := v.Decode([]byte{0xF0, 0x02, 0x01})
	if v.Type != 240 {
		t.Error("type %d, %d", v.Type, 240)
	}

	if len(r) != 1 || r[0] != 1 {
		t.Error("remain %d, %d", len(r), r[0])
	}
	r, _ = v.Decode([]byte{0xF0, 0x4, 0x01, 0x02})
	if v.Value[0] != 1 || v.Value[1] != 2 {
		t.Error("value %d, %d", v.Value[0], v.Value[1])
	}
	if len(r) != 0 {
		t.Error("remain %d, %d", len(r))
	}
}

func TestEncode(t *testing.T) {
	v := new(TLV)
	v.Decode([]byte{0xF0, 0x4, 0x01, 0x02})
	if b, _ := v.Encode(); !EqualBytes(b, []byte{0xF0, 0x4, 0x01, 0x02}) {
		t.Error(v.Encode())
	}
}

func TestName(t *testing.T) {
	s1 := "/ucla/edu/cs/ndn"
	s2 := Uri(Name(s1))
	if s1 != s2 {
		t.Error("expected %v, got %v", s1, s2)
	}
}

func TestIO(t *testing.T) {
	v1 := true
	tlv := new(TLV)
	tlv.Write(v1)
	var v2 bool
	tlv.Read(&v2)
	if !v2 {
		t.Error("should be true %v", tlv.Value)
	}
	f1 := 0.5
	tlv.Write(f1)
	var f2 float64
	tlv.Read(&f2)
	if f1 != f2 {
		t.Error("expected %v, got %v", f1, f2)
	}
	s1 := "hello"
	tlv.Write(s1)
	var s2 string
	tlv.Read(&s2)
	if s1 != s2 {
		t.Error("expected %v, got %v", s1, s2)
	}
}

func TestDecodeSimpleInterest(t *testing.T) {
	name := new(TLV)
	name.Type = NAME
	nonce := new(TLV)
	nonce.Type = NONCE
	// create selector
	selectors := new(TLV)
	selectors.Type = SELECTORS

	max := new(TLV)
	max.Type = MAX_SUFFIX_COMPONENT
	exclude := new(TLV)
	exclude.Type = EXCLUDE
	namecomp := new(TLV)
	namecomp.Type = NAME_COMPONENT
	exclude.Add(namecomp)
	exclude.Add(namecomp)
	exclude.Add(namecomp)
	exclude.Add(namecomp)

	selectors.Add(max)
	selectors.Add(exclude)

	lifetime := new(TLV)
	lifetime.Type = INTEREST_LIFETIME

	interest := new(TLV)
	interest.Type = INTEREST
	interest.Add(name)
	interest.Add(selectors)
	interest.Add(nonce)
	interest.Add(lifetime)

	b, err := interest.Encode()
	if err != nil {
		t.Error(err)
	}
	ip, err := DecodeInterest(b)
	if err != nil {
		t.Error(err)
	}
	if len(ip.Children) != len(interest.Children) {
		t.Error("children count", "expected", len(interest.Children), "actual", len(ip.Children))
	}
	b2, _ := ip.Encode()
	if !EqualBytes(b, b2) {
		t.Error(b, b2)
	}
}

func EqualBytes(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}