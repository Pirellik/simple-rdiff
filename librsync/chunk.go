package librsync

import (
	"encoding/binary"
	"fmt"
	"io"
)

type chunkType byte

const (
	chunkTypeReusable chunkType = iota
	chunkTypeModified
)

type chunk interface {
	chunkType() chunkType
	append(chunk)
	write(io.Writer) error
	patch(io.ReadSeeker, io.Writer) error
}

type reusable struct {
	startPosition uint64
	length        uint64
}

type modified struct {
	data []byte
}

func (r *reusable) chunkType() chunkType { return chunkTypeReusable }
func (m *modified) chunkType() chunkType { return chunkTypeModified }

func (r *reusable) append(c chunk) {
	casted := c.(*reusable)
	r.length += casted.length
}

func (m *modified) append(c chunk) {
	casted := c.(*modified)
	m.data = append(m.data, casted.data...)
}

func (r *reusable) write(out io.Writer) error {
	if err := binary.Write(out, binary.BigEndian, r.chunkType()); err != nil {
		return err
	}
	if err := binary.Write(out, binary.BigEndian, r.startPosition); err != nil {
		return err
	}
	if err := binary.Write(out, binary.BigEndian, r.length); err != nil {
		return err
	}
	return nil
}

func (m *modified) write(out io.Writer) error {
	if err := binary.Write(out, binary.BigEndian, m.chunkType()); err != nil {
		return err
	}
	if err := binary.Write(out, binary.BigEndian, uint64(len(m.data))); err != nil {
		return err
	}
	if err := binary.Write(out, binary.BigEndian, m.data); err != nil {
		return err
	}
	return nil
}

func (r *reusable) patch(base io.ReadSeeker, out io.Writer) error {
	if _, err := base.Seek(int64(r.startPosition), io.SeekStart); err != nil {
		return err
	}
	_, err := io.CopyN(out, base, int64(r.length))
	return err
}

func (m *modified) patch(base io.ReadSeeker, out io.Writer) error {
	_, err := out.Write(m.data)
	return err
}

func readChunk(in io.Reader) (chunk, error) {
	var cType chunkType
	if err := binary.Read(in, binary.BigEndian, &cType); err != nil {
		return nil, err
	}
	switch cType {
	case chunkTypeReusable:
		var startPosition uint64
		if err := binary.Read(in, binary.BigEndian, &startPosition); err != nil {
			return nil, err
		}
		var length uint64
		if err := binary.Read(in, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		return &reusable{
			startPosition: startPosition,
			length:        length,
		}, nil
	case chunkTypeModified:
		var length uint64
		if err := binary.Read(in, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		data := make([]byte, length)
		n, err := in.Read(data)
		if err != nil {
			return nil, err
		}
		if uint64(n) != length {
			return nil, fmt.Errorf("corrupted chunk - length mismatch, got = %d, want = %d", n, length)
		}
		return &modified{data: data}, nil
	default:
		return nil, fmt.Errorf("corrupted chunk - unknown type = %x", cType)
	}
}
