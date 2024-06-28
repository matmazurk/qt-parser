package parser

import "strings"

type search struct {
	atomPath     string
	searchParams searchParams
}

// builder is helper structure to build parser with given search params with ease
type builder struct {
	searches []search
}

func NewBuilder() *builder {
	return &builder{}
}

// Find registers new search params, which will be included when building the parser
// atomPath should include all needed atoms formatted in unix path way, e.g.: 'moov/trak/tkhd'
// findingName is an alias for the single search and can be used later for reference
// if given atom is found, [offset; offset+bytesAmount] will be read from it
func (b *builder) Find(atomPath string, offset, bytesAmount uint64, findingName string) *builder {
	b.searches = append(b.searches, search{
		atomPath: atomPath,
		searchParams: searchParams{
			offset:      offset,
			bytesAmount: bytesAmount,
			findingName: findingName,
		},
	})
	return b
}

func (b *builder) Build() parser {
	ret := newParser()
	for _, s := range b.searches {
		currAtom := ret.root
		atomTypes := strings.Split(s.atomPath, "/")
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
		currAtom.searchParams = append(currAtom.searchParams, s.searchParams)
	}
	return ret
}
