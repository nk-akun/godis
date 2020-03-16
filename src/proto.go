package godis

import (
	"bytes"
	"errors"
	"io"
	"log"

	"github.com/nk-akun/godis/src/util/bufio2"
)

const (
	TypeStatus    = '+'
	TypeError     = '-'
	TypeInt       = ':'
	TypeBulk      = '$'
	TypeMultiBulk = '*'
)

// EncodeData is the struct to store encoding data of cmd
type EncodeData struct {
	Err   error
	Value []byte
	Type  byte
	Array []*EncodeData
}

// EncodeCmd is func to encode client cmd
func EncodeCmd(content string) (b []byte, err error) {
	return EncodeBytes([]byte(content))
}

// EncodeBytes ...
func EncodeBytes(bytesContent []byte) (re []byte, err error) {
	chunks := bytes.Split(bytesContent, []byte(" "))
	if chunks == nil {
		return nil, errorTrace(errorNew("result of split is nil"))
	}

	r := NewMultiBulk(nil)
	for _, chunk := range chunks {
		if len(chunk) > 0 {
			r.Array = append(r.Array, NewBulk(chunk))
		}
	}
	return EncodeMultiBulk(r)
}

// EncodeMultiBulk ...
func EncodeMultiBulk(r *EncodeData) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := Encode(b, r); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Encode ...
func Encode(w io.Writer, r *EncodeData) error {
	return NewEncoder(w).Encode(r, true) //true means it must flush the buffer even though the buf is not full yet.
}

// NewEncoder ...
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{ByteWriter: bufio2.NewWriterSize(w, 10240)}
}

// NewEncoderSize ...
func NewEncoderSize(w io.Writer, size int) *Encoder {
	return &Encoder{ByteWriter: bufio2.NewWriterSize(w, size)}
}

// NewMultiBulk generates a new struct whose type is TypeMultiBulk
func NewMultiBulk(array []*EncodeData) *EncodeData {
	ans := &EncodeData{}
	ans.Type = TypeMultiBulk
	ans.Array = array
	return ans
}

// NewBulk generate a new struct whose type is TypeBulk
func NewBulk(data []byte) *EncodeData {
	ans := &EncodeData{}
	ans.Type = TypeBulk
	ans.Value = data
	return ans
}

func errorTrace(err error) error {
	if err != nil {
		log.Println("error tracing: ", err.Error())
	}
	return err
}

func errorNew(errorMsg string) error {
	return errors.New("error: " + errorMsg)
}

// Encoder ...
type Encoder struct {
	Err        error
	ByteWriter *bufio2.Writer
}

// Encode ...
func (e *Encoder) Encode(r *EncodeData, flush bool) error {
	if err := e.encodeData(r); err != nil {
		e.Err = err
	} else if flush {
		e.Err = e.ByteWriter.Flush()
	}
	return e.Err
}

func (e *Encoder) encodeData(r *EncodeData) error {
	if err := e.ByteWriter.WriteByte(byte(r.Type)); err != nil {
		e.Err = err
		return err
	}

	// switch r.Type {
	// case :

	// }
	return nil
}
