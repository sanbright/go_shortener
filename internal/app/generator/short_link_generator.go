package generator

import "github.com/dchest/uniuri"

type ShortLinkGenerator struct {
}

func NewShortLinkGenerator() *ShortLinkGenerator {
	return &ShortLinkGenerator{}
}

func (generator *ShortLinkGenerator) UniqGenerate() string {
	return uniuri.NewLen(10)
}
