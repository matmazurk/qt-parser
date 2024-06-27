package parser

type atomType string

const rootAtom atomType = "root"

type searchParams struct {
	offset      uint64
	bytesAmount uint64
	findingName string
}

type atom struct {
	typ    atomType
	childs map[atomType]*atom

	params []searchParams
}
