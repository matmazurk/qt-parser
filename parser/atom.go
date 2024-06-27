package parser

type atomType string

const rootAtom atomType = "root"

type searchParams struct {
	offset      uint
	bytesAmount uint
	findingName string
}

type atom struct {
	typ    atomType
	childs map[atomType][]atom

	params []searchParams
}
