package parser

type atomType string

const rootAtom atomType = "root"

type searchParams struct {
	// offset within whole atom, including header
	offset uint64
	// amount of bytes to read [offset; offset+bytesAmount]
	bytesAmount uint64
	// alias for given search
	findingName string
}

type atom struct {
	typ    atomType
	childs map[atomType]*atom

	// data to read for given atom
	// should be filled only for leaf atoms
	searchParams []searchParams
}

func (a *atom) IsLeaf() bool {
	return len(a.searchParams) > 0
}
