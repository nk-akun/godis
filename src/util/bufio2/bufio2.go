package bufio2

import "io"

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

func (w *Writer) free() int {
	return len(w.buf) - w.rpos
}
