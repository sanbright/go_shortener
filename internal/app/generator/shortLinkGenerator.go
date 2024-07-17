package generator

import "github.com/dchest/uniuri"

type ShortLinkGeneratorInterface interface {
	UniqGenerate() string
}

type ShortLinkGenerator struct {
	length int
}

func NewShortLinkGenerator(length int) *ShortLinkGenerator {
	return &ShortLinkGenerator{length: length}
}

func (generator *ShortLinkGenerator) UniqGenerate() string {
	return uniuri.NewLen(generator.length)
}
