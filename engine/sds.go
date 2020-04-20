package godis

import "github.com/nk-akun/godis/engine/util"

// Sdshdr stores string aims to match the characteristics of godis
type Sdshdr struct {
	len int
	cap int
	buf []byte
}

// SdsNewString return a new Sdshdr with default capacity 20
func SdsNewString(str *string) *Sdshdr {
	l := len(*str)
	return &Sdshdr{
		len: l,
		cap: l,
		buf: []byte(*str),
	}
}

// SdsNewBuf return a new Sdshdr
func SdsNewBuf(buf []byte) *Sdshdr {
	l := len(buf)
	sds := &Sdshdr{
		len: l,
		cap: l,
		buf: make([]byte, l),
	}
	copy(sds.buf, buf)
	return sds
}

// SdsNewEmpty return a new Sdshdr with 0 size
func SdsNewEmpty() *Sdshdr {
	return &Sdshdr{
		len: 0,
		cap: 0,
		buf: make([]byte, 0),
	}
}

// SdsGetString return string pointer of buf
func (sds *Sdshdr) SdsGetString() *string {
	str := string(sds.buf[:sds.len])
	return &str
}

// SdsGetBuf return the buf of sds
func (sds *Sdshdr) SdsGetBuf() []byte {
	buf := make([]byte, sds.len)
	copy(buf, sds.buf)
	return buf
}

// SdsCopy uses str to replace previously existing content
func (sds *Sdshdr) SdsCopy(str *string) {
	l := len(*str)
	if sds.cap >= l {
		copy(sds.buf[0:], *str)
	} else {
		sds.buf = []byte(*str)
		sds.cap = l
	}
	sds.len = l
}

// SdsLen return the length of sds
func (sds *Sdshdr) SdsLen() int {
	return sds.len
}

// Sdsavail return the size of available space
func (sds *Sdshdr) Sdsavail() int {
	return sds.cap - sds.len
}

// SdsClear clear the content in sds
func (sds *Sdshdr) SdsClear() {
	sds.len = 0
}

//TODO: Optimize join performance

// SdsCat appends str to end of buf
func (sds *Sdshdr) SdsCat(str *string) {
	l := len(*str)
	if sds.Sdsavail() >= l {
		copy(sds.buf[sds.len:], *str)
	} else {
		buf := make([]byte, sds.len+l)
		copy(buf, sds.buf)
		copy(buf[sds.len:], *str)
		sds.buf = buf
		sds.cap = sds.len + l
	}
	sds.len += l
}

// SdsCatSds append s.buf to sds
func (sds *Sdshdr) SdsCatSds(s *Sdshdr) {
	l := s.len
	if sds.Sdsavail() >= l {
		copy(sds.buf[sds.len:], s.buf)
	} else {
		buf := make([]byte, sds.len+l)
		copy(buf, sds.buf)
		copy(buf[sds.len:], s.buf)
		sds.buf = buf
		sds.cap = sds.len + s.len
	}
	sds.len += l
}

// SdsCmp return 1 if the lexicographical order of s1 is larger than s2,-1 smaller,0 equal
func SdsCmp(s1 *Sdshdr, s2 *Sdshdr) int {
	limit := util.Min(s1.len, s2.len)
	for i := 0; i < limit; i++ {
		if s1.buf[i] < s2.buf[i] {
			return -1
		} else if s1.buf[i] > s2.buf[i] {
			return 1
		}
	}
	if s1.len > s2.len {
		return 1
	} else if s1.len < s2.len {
		return -1
	} else {
		return 0
	}
}
