package bufio2

import (
	"bufio"
	"bytes"
	"io"
)

// Writer is the struct contains the io.Writer and a buffer data
type Writer struct {
	err    error
	writer io.Writer
	buf    []byte // buf stores buffer data before writing to writer whose aim is reducing the number of writes
	rpos   int    // the right edge of buf
}

// NewWriterSize ...
func NewWriterSize(w io.Writer, size int) *Writer {
	return &Writer{
		writer: w,
		rpos:   0,
		buf:    make([]byte, size),
	}
}

// NewWriter ...
func NewWriter(w io.Writer) *Writer {
	return NewWriterSize(w, 2048)
}

// Flush to flush the buffer data stored at
func (w *Writer) Flush() error {
	if w.rpos == 0 {
		return nil
	}

	n, err := w.writer.Write(w.buf)
	if err != nil {
		w.err = err
	} else if n < w.rpos {
		w.err = io.ErrShortWrite
	} else {
		w.rpos = 0
	}
	return w.err
}

// WriteByte is the func to write a byte b into buf or io.writer
func (w *Writer) WriteByte(b byte) error {
	if w.free() == 0 && w.Flush() != nil {
		return w.err
	}
	w.buf[w.rpos] = b
	w.rpos++
	return nil
}

// WriteBytes is the func to write []byte into buf or io.writer
func (w *Writer) WriteBytes(value []byte) (wn int, err error) {
	var n int

	for w.err == nil && w.free() < len(value) {
		if w.rpos == 0 {
			n, w.err = w.writer.Write(value)
		} else { //make sure that the bytes stored in buf could be writen first
			n = copy(w.buf[w.rpos:], value)
			w.rpos += n
			w.err = w.Flush()
		}
		value = value[n:]
		wn += n
	}

	if len(value) == 0 || w.err != nil {
		return wn, w.err
	}

	n = copy(w.buf[w.rpos:], value)
	w.rpos += n
	return wn + n, nil
}

// WriteString is the func to write string into buf or io.writer
func (w *Writer) WriteString(str string) (wn int, err error) {
	var n int
	for w.err == nil && w.free() < len(str) {
		if w.rpos == 0 {
			n, w.err = w.writer.Write([]byte(str))
		} else {
			n = copy(w.buf[w.rpos:], str)
			w.rpos += n
			w.err = w.Flush()
		}
		str = str[n:]
		wn += n
	}

	if len(str) == 0 || w.err != nil {
		return wn, err
	}

	n = copy(w.buf[w.rpos:], str)
	w.rpos += n
	return wn + n, nil
}

func (w *Writer) free() int {
	return len(w.buf) - w.rpos
}

// Reader contains io.Reader and buf ect.
type Reader struct {
	reader io.Reader
	buf    []byte
	lpos   int
	rpos   int
	err    error
}

// NewReaderSize return Reader pointer
func NewReaderSize(reader io.Reader, size int) *Reader {
	if size < 0 {
		size = 2014
	}
	return &Reader{
		reader: reader,
		buf:    make([]byte, size),
		lpos:   0,
		rpos:   0,
	}
}

// GlanceByte return first byte of buf
func (r *Reader) GlanceByte() (byte, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.bufferLen() == 0 && r.fillBuf() != nil {
		return 0, r.err
	}
	return r.buf[r.lpos], nil
}

// ReadByte return first byte of buf and move the lpos
func (r *Reader) ReadByte() (byte, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.bufferLen() == 0 && r.fillBuf() != nil {
		return 0, r.err
	}
	b := r.buf[r.lpos]
	r.lpos++
	return b, nil
}

// ReadBytesDelim return bytes from r.buf[r.lpos] to delim in r.buf
func (r *Reader) ReadBytesDelim(delim byte) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	for {
		pos := bytes.IndexByte(r.buf[r.lpos:r.rpos], delim)
		if pos > 0 {
			newLpos := r.lpos + pos + 1
			ans := r.buf[r.lpos:newLpos]
			r.lpos = newLpos
			return ans, nil
		}
		if r.bufferLen() == len(r.buf) { // can't find the delim though r.buf is full
			r.lpos = r.rpos
			return r.buf, bufio.ErrBufferFull
		}
		if r.fillBuf() != nil {
			return nil, r.err
		}
	}
}

// ReadBytesLen read bytes which length is n from r.lpos
func (r *Reader) ReadBytesLen(n int) ([]byte, error) {
	re := make([]byte, n)
	for l := 0; l < n; {
		if length := r.bufferLen(); length >= (n-l) || length == len(r.buf) {
			tmp := copy(re[l:], r.buf[r.lpos:])
			l += tmp
			r.lpos += tmp
		}
		if err := r.fillBuf(); err != nil {
			return nil, err
		}
	}
	return re, nil
}

// bufferLen return valid length of buf
func (r *Reader) bufferLen() int {
	return r.rpos - r.lpos
}

func (r *Reader) fillBuf() error {
	if r.err != nil {
		return r.err
	}

	//move the existing data to the front
	n := copy(r.buf, r.buf[r.lpos:r.rpos])
	r.lpos = 0
	r.rpos = n

	n, err := r.reader.Read(r.buf[r.rpos:])
	if err != nil {
		r.err = err
	} else if n == 0 {
		r.err = io.ErrNoProgress
	} else {
		r.rpos += n
	}
	return r.err
}
