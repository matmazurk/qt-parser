package parser

import (
	"encoding/binary"
	"io"
	"maps"

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

func (p parser) Parse(r io.Reader) (map[string]string, error) {
	ret, _, err := handleAtom(r, p.toFind, 0, 0)
	return ret, err
}

func handleAtom(r io.Reader, a *atom, size, alreadyRead uint64) (map[string]string, uint64, error) {
	if len(a.params) > 0 {
		buf := make([]byte, size-alreadyRead)
		n, err := r.Read(buf)
		if err != nil {
			return nil, uint64(n), errors.Wrapf(err, "could not read bytes from leaf atom '%s'", a.typ)
		}
		ret := map[string]string{}
		for _, sp := range a.params {
			val := buf[sp.offset-alreadyRead : sp.offset-alreadyRead+sp.bytesAmount]
			ret[sp.findingName] = string(val)
		}
		return ret, size, nil
	}

	rett := map[string]string{}
	readBytes := alreadyRead
	for readBytes < size || size == 0 {
		var loopReadBytes uint64 = 0
		header := make([]byte, u64Bytes)
		n, err := r.Read(header)
		loopReadBytes += uint64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return rett, uint64(n), nil
			}
			return nil, uint64(n), errors.Wrap(err, "could not read atom header")
		}

		atomSize := uint64(binary.BigEndian.Uint32(header[:u32Bytes]))
		currAtomType := atomType(header[u32Bytes:])
		if atomSize == extendedSize {
			bigSize := make([]byte, u64Bytes)
			n, err := r.Read(bigSize)
			loopReadBytes += uint64(n)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return rett, readBytes, nil
				}
			}
			atomSize = binary.BigEndian.Uint64(bigSize)
		}

		if currAtom, ok := a.childs[currAtomType]; ok {
			ret, _, err := handleAtom(r, currAtom, atomSize, loopReadBytes)
			if err != nil {
				return nil, readBytes, errors.Wrapf(err, "could not handle atom '%s'", currAtomType)
			}
			readBytes += atomSize
			maps.Copy(rett, ret)
		} else {
			_, err := skipBytes(r, atomSize-loopReadBytes)
			readBytes += atomSize
			if err != nil {
				if errors.Is(err, io.EOF) {
					return rett, readBytes, nil
				}
				return nil, readBytes, errors.Wrapf(err, "could not skip %d bytes of '%s' atom type", atomSize, currAtomType)
			}
		}
	}
	return rett, readBytes, nil
}

func skipBytes(r io.Reader, count uint64) (int, error) {
	buf := make([]byte, count)
	return r.Read(buf)
}

func (p parser) ToFind() *atom {
	return p.toFind
}
