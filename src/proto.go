package godis

import (
	"bytes"
	"errors"
	"io"
	"strconv"

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
		return nil, errorNew("result of split is nil")
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
		log.Errorf("encode data: %v", err)
		e.Err = err
		return err
	}

	switch r.Type {
	case TypeStatus, TypeError, TypeInt:
		return e.encodeTextBytes(r.Value)
	case TypeBulk:
		return e.encodeBulkBytes(r.Value)
	case TypeMultiBulk:
		return e.encodeMultiBulkArray(r.Array)
	}
	return nil
}

func (e *Encoder) encodeTextBytes(value []byte) error {
	if n, err := e.ByteWriter.WriteBytes(value); n != len(value) || err != nil {
		return err
	}
	if n, err := e.ByteWriter.WriteBytes([]byte("\r\n")); n != 2 || err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeTextString(str string) error {
	if n, err := e.ByteWriter.WriteString(str); n != len(str) || err != nil {
		return err
	}
	if n, err := e.ByteWriter.WriteString("\r\n"); n != 2 || err != nil {
		return err
	}
	return nil
}

func (e *Encoder) encodeInt(v int64) error {
	return e.encodeTextString(strconv.FormatInt(v, 10))
}

func (e *Encoder) encodeBulkBytes(value []byte) error {
	if err := e.encodeInt(int64(len(value))); err != nil {
		return err
	}

	return e.encodeTextBytes(value)
}

func (e *Encoder) encodeMultiBulkArray(bulks []*EncodeData) error {
	if len(bulks) == 0 {
		return e.encodeInt(-1)
	}
	if err := e.encodeInt(int64(len(bulks))); err != nil {
		return err
	}
	for _, bulk := range bulks {
		if err := e.encodeData(bulk); err != nil {
			return err
		}
	}
	return nil
}
