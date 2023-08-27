package core

const GeoAlphabet string = "0123456789bcdefghjkmnpqrstuvwxyz"

type encoding struct {
	encode string
	decode [256]byte
}

func (e *encoding) Decode(s string) uint64 {
	var x uint64

	decode := [255]byte{}
	for i := 0; i < len(GeoAlphabet); i++ {
		decode[GeoAlphabet[i]] = byte(i)
	}
	for i := 0; i < 10; i++ {
		x = (x << 5) | uint64(decode[s[i]])
	}
	return x
}

/*
break x into 5-bit blocks and map each block to a character in GeoAlphabet.
If x is 52-bit long, the 2 last bits are encoded as 0. Example:

	  0b10010 11010 10010 10110 10100 10101 10101 00101 01101 01001 01
		    v     u     q     q     q     p     p     5     e     9  0
*/
func (e *encoding) Encode(x uint64) string {
	b := [11]byte{}
	for i := 0; i < 11; i++ {
		shift := 52 - (i+1)*5
		if shift <= 0 {
			b[i] = GeoAlphabet[0]
			break
		}
		idx := (x >> shift) & 0b11111
		b[i] = GeoAlphabet[idx]
	}
	return string(b[:])
}

func newBase32Encoding() *encoding {
	e := &encoding{
		encode: GeoAlphabet,
		decode: [256]byte{},
	}

	for i := 0; i < len(e.decode); i++ {
		e.decode[i] = 0xff
	}
	for i := 0; i < len(e.encode); i++ {
		e.decode[e.encode[i]] = byte(i)
	}
	return e
}

var Base32encoding = newBase32Encoding()
