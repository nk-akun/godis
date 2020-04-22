package godis

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/nk-akun/godis/engine/util/bufio2"
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

// NewError ...
func NewError(data []byte) *EncodeData {
	ans := &EncodeData{}
	ans.Type = TypeError
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

// Decoder ...
type Decoder struct {
	ByteReader *bufio2.Reader
	Err        error
}

// NewDecoder generate a decoder
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{ByteReader: bufio2.NewReaderSize(reader, 2048)}
}

// DecodeMultiBulks decode multibulks into several parts
// for example,*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n5\r\n ---> set num 5
func (d *Decoder) DecodeMultiBulks() ([]*EncodeData, error) {
	result, err := d.decodeMultiBulks()
	if err != nil {
		d.Err = err
		return nil, err
	}
	return result, nil
}

func (d *Decoder) decodeMultiBulks() ([]*EncodeData, error) {
	t, err := d.ByteReader.GlanceByte()
	if err != nil {
		return nil, err
	}
	if byte(t) != TypeMultiBulk {
		return nil, errorNew("command format error!")
	}

	if _, err := d.ByteReader.ReadByte(); err != nil {
		return nil, err
	}

	n, err := d.decodeInt()
	if err != nil {
		log.Errorf("decode int err:%v", err)
		return nil, err
	}

	bulks := make([]*EncodeData, n)
	for i := range bulks {
		bulk, err := d.decodeData()
		if err != nil {
			return nil, err
		}
		if bulk.Type != TypeBulk {
			log.Errorf("bad bulk content,should be bulkBytes")
			return nil, errorNew("bad bulk content")
		}
		bulks[i] = bulk
	}
	return bulks, nil
}

func (d *Decoder) decodeData() (*EncodeData, error) {
	t, err := d.ByteReader.ReadByte()
	if err != nil {
		return nil, err
	}
	e := &EncodeData{}
	e.Type = byte(t)

	switch e.Type {
	case TypeStatus, TypeError, TypeInt:
		e.Value, err = d.decodeTextBytes()
	case TypeBulk:
		e.Value, err = d.decodeBulkBytes()
	case TypeMultiBulk:
		e.Array, err = d.decodeMultiBulkArray()
	}
	return e, err
}

func (d *Decoder) decodeTextBytes() ([]byte, error) {
	b, err := d.ByteReader.ReadBytesDelim('\n')
	if err != nil {
		return nil, err
	}

	return b[:len(b)-2], err
}

func (d *Decoder) decodeBulkBytes() ([]byte, error) {
	n, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	b, err := d.ByteReader.ReadBytesLen(int(n) + 2)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}

func (d *Decoder) decodeMultiBulkArray() ([]*EncodeData, error) {
	n, err := d.decodeInt()
	if err != nil {
		return nil, err
	}
	bulks := make([]*EncodeData, n)
	for i := range bulks {
		bulk, err := d.decodeData()
		if err != nil {
			return nil, err
		}
		bulks[i] = bulk
	}
	return bulks, nil
}

func (d *Decoder) decodeInt() (int64, error) {
	num, err := d.ByteReader.ReadBytesDelim('\n')
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(num[:(len(num)-2)]), 10, 64)
}
