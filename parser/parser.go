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
	root *atom
}

func newParser() parser {
	return parser{
		root: &atom{
			typ:    rootAtom,
			childs: map[atomType]*atom{},
		},
	}
}

// Parse tries to recursively look for child atoms of its root atom, reading given data from leaf atoms
// returned map consists of findingName as a key and the findings as value
// the findings for single findingName are always in the same order as particular atoms in parsed file
func (p parser) Parse(r io.Reader) (map[string][][]byte, error) {
	return handleAtom(r, p.root, 0, 0)
}

func (p parser) Root() *atom {
	return p.root
}

// handleAtom reads desired data for leaf atoms or recursively propagetes itself among container atom childs
// reader should be exacly after atom's header position
// size 0 denotes root atom
// alreadyRead should be set to atom's header size
func handleAtom(r io.Reader, a *atom, size, alreadyRead uint64) (map[string][][]byte, error) {
	if a.IsLeaf() {
		buf := make([]byte, size-alreadyRead)
		_, err := r.Read(buf)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read bytes from leaf atom '%s'", a.typ)
		}
		ret := map[string][][]byte{}
		for _, sp := range a.searchParams {
			ret[sp.findingName] = append(ret[sp.findingName], buf[sp.offset-alreadyRead:sp.offset-alreadyRead+sp.bytesAmount])
		}
		return ret, nil
	}

	findings := map[string][][]byte{}
	readBytes := alreadyRead
	for readBytes < size || size == 0 {
		var loopReadBytes uint64 = 0
		header := make([]byte, u64Bytes)
		n, err := r.Read(header)
		loopReadBytes += uint64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return findings, nil
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
			appendMap(findings, ret)
		} else {
			err := skipBytes(r, atomSize-loopReadBytes)
			readBytes += atomSize
			if err != nil {
				return nil, errors.Wrapf(err, "could not skip %d bytes of '%s' atom type", atomSize, currAtomType)
			}
		}
	}
	return findings, nil
}

// skipBytes mindlessly skips given count of bytes in a reader
func skipBytes(r io.Reader, count uint64) error {
	buf := make([]byte, count)
	_, err := r.Read(buf)
	return err
}

// appendMap appends src to dst by extending dst with new key/values
func appendMap(dst, src map[string][][]byte) {
	for k, v := range src {
		if _, ok := dst[k]; !ok {
			dst[k] = v
		} else {
			dst[k] = append(dst[k], v...)
		}
	}
}
