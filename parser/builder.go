package parser

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
	return parser{}
}
