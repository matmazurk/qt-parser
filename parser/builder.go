package parser

import "strings"

type search struct {
	atomsPath   string
	offset      uint
	bytesAmount uint
	findingName string
}

type builder struct {
	searches []search
}

func NewBuilder() *builder {
	return &builder{}
}

func (b *builder) Find(atomsPath string, offset, bytesAmount uint, findingName string) *builder {
	b.searches = append(b.searches, search{
		atomsPath:   atomsPath,
		offset:      offset,
		bytesAmount: bytesAmount,
		findingName: findingName,
	})
	return b
}

func (b *builder) Build() parser {
	ret := newParser()
	for _, s := range b.searches {
		currAtom := ret.toFind
		atomTypes := strings.Split(s.atomsPath, "/")
		if len(atomTypes) == 0 {
			continue
		}
		for _, at := range atomTypes {
			cat := atomType(at)
			if a, ok := currAtom.childs[cat]; ok {
				currAtom = a
			} else {
				newAtom := &atom{
					typ:    cat,
					childs: map[atomType]*atom{},
				}
				currAtom.childs[cat] = newAtom
				currAtom = newAtom
			}
		}
		currAtom.params = append(currAtom.params, searchParams{
			offset:      s.offset,
			bytesAmount: s.bytesAmount,
			findingName: s.findingName,
		})
	}
	return ret
}
