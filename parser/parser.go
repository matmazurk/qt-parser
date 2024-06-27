package parser

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	u32Bytes = 4
	u64Bytes = 8

	extendedSize = 1
)

type parser struct {
	toFind *atom
}

func newParser() parser {
	return parser{
		toFind: &atom{
			typ:    rootAtom,
			childs: map[atomType]*atom{},
		},
	}
}

func (p parser) Parse(r io.Reader) (map[string][][]byte, error) {
	return handleAtom(r, p.toFind, 0, 0)
}

func handleAtom(r io.Reader, a *atom, size, alreadyRead uint64) (map[string][][]byte, error) {
	if len(a.params) > 0 {
		buf := make([]byte, size-alreadyRead)
		_, err := r.Read(buf)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read bytes from leaf atom '%s'", a.typ)
		}
		ret := map[string][][]byte{}
		for _, sp := range a.params {
			ret[sp.findingName] = append(ret[sp.findingName], buf[sp.offset-alreadyRead:sp.offset-alreadyRead+sp.bytesAmount])
		}
		return ret, nil
	}

	rett := map[string][][]byte{}
	readBytes := alreadyRead
	for readBytes < size || size == 0 {
		var loopReadBytes uint64 = 0
		header := make([]byte, u64Bytes)
		n, err := r.Read(header)
		loopReadBytes += uint64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return rett, nil
			}
			return nil, errors.Wrap(err, "could not read atom header")
		}
		if n < u64Bytes {
			readBytes += uint64(n)
			continue
		}

		atomSize := uint64(binary.BigEndian.Uint32(header[:u32Bytes]))
		currAtomType := atomType(header[u32Bytes:])
		if atomSize == extendedSize {
			bigSize := make([]byte, u64Bytes)
			n, err := r.Read(bigSize)
			loopReadBytes += uint64(n)
			if err != nil {
				return nil, err
			}
			atomSize = binary.BigEndian.Uint64(bigSize)
		}

		if currAtom, ok := a.childs[currAtomType]; ok {
			ret, err := handleAtom(r, currAtom, atomSize, loopReadBytes)
			if err != nil {
				return nil, errors.Wrapf(err, "could not handle atom '%s'", currAtomType)
			}
			readBytes += atomSize
			appendMap(rett, ret)
		} else {
			err := skipBytes(r, atomSize-loopReadBytes)
			readBytes += atomSize
			if err != nil {
				return nil, errors.Wrapf(err, "could not skip %d bytes of '%s' atom type", atomSize, currAtomType)
			}
		}
	}
	return rett, nil
}

func skipBytes(r io.Reader, count uint64) error {
	buf := make([]byte, count)
	_, err := r.Read(buf)
	return err
}

func (p parser) ToFind() *atom {
	return p.toFind
}

func appendMap(dst, src map[string][][]byte) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		} else {
			dst[k] = append(dst[k], v...)
		}
	}
}
