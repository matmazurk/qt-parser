package parser

import "io"

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
	return nil, nil
}

func (p parser) ToFind() *atom {
	return p.toFind
}
