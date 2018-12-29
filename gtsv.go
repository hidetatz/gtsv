package gtsv

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unsafe"
)

// Reader reads TSV.
type Reader struct {
	reader   io.Reader
	readBuff []byte // temporary buffer which stores line
	colBuff  []byte // buffer which stores current column
	readErr  error
	col      int
	row      int
	err      error

	buff [6 << 10]byte // large enough
}

// New returnds new TSV reader from io.Reader
func New(r io.Reader) *Reader {
	return &Reader{reader: r, err: nil}
}

// Error returns error
func (gr *Reader) Error() error {
	return gr.err
}

// Next returns true when next row exists.
func (gr *Reader) Next() bool {
	if gr.err != nil {
		return false
	}

	gr.row++
	for {
		if len(gr.tmpBuff) <= 0 {
			fmt.Println("aaa")
			if gr.rErr != nil {
				if gr.err != io.EOF {
					gr.err = fmt.Errorf("cannot read row #%d: %s", gr.row, gr.err)
				}
				return false
			}
			n, err := gr.reader.Read(gr.buff[:]) // first, read and get some bytes and store to buffer
			gr.rErr = err
			gr.tmpBuff = gr.buff[:n]
		}

		fmt.Println("bbb")
		n := bytes.IndexByte(gr.tmpBuff, '\n') // read from buffer
		if n >= 0 {
			tmp := gr.tmpBuff[:n]

			// next row found
			gr.tmpBuff = gr.tmpBuff[n+1:]
			gr.colBuff = tmp
			return true
		}
		fmt.Println("ccc")
	}
	// implement if cannot find \n
}

// Int returns next int column
func (gr *Reader) Int() int {
	if gr.err != nil {
		return 0
	}
	b, err := gr.nextColumn()
	if err != nil {
		gr.err = fmt.Errorf("cannot read `int`: %s", err)
		return 0
	}

	n, err := strconv.Atoi(bytesToString(b))
	if err != nil {
		gr.err = fmt.Errorf("cannot parse `int`: %s", err)
		return 0
	}
	return n
}

func (gr *Reader) nextColumn() ([]byte, error) {
	gr.col++
	if gr.tmpBuff == nil {
		return nil, fmt.Errorf("no more columns")
	}

	n := bytes.IndexByte(gr.colBuff, '\t') // look for tab
	if n < 0 {
		// tab is not found, the most right column
		tmp := gr.colBuff
		gr.colBuff = nil
		return tmp, nil
	}
	tmp := gr.colBuff[:n]
	gr.colBuff = gr.colBuff[n+1:]
	return tmp, nil
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) // faster than string(b)
}
