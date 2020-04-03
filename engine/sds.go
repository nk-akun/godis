package godis

// Sdshdr stores string aims to match the characteristics of godis
type Sdshdr struct {
	len int
	cap int
	buf []byte
}

// SdsNew return a new Sdshdr with default capacity 20
func SdsNew(str *string) *Sdshdr {
	l := len(*str)
	return &Sdshdr{
		len: l,
		cap: l,
		buf: []byte(*str),
	}
}

// SdsNewEmpty return a new Sdshdr with capacity size
func SdsNewEmpty() *Sdshdr {
	return &Sdshdr{
		len: 0,
		cap: 0,
		buf: make([]byte, 0),
	}
}

// SdsGetString return string pointer of buf
func (s *Sdshdr) SdsGetString() *string {
	str := string(s.buf[:s.len])
	return &str
}

func (s *Sdshdr) SdsSetString() {

}
