package generator

import "github.com/dchest/uniuri"

// IShortLinkGenerator интерфейс генератора строк, для коротких ссылок.
type IShortLinkGenerator interface {
	UniqGenerate() string
}

// ShortLinkGenerator генератор уникальных строк, фиксированной длинны.
type ShortLinkGenerator struct {
	// length - длинна генерируемых строк
	length int
}

// NewShortLinkGenerator length зададет длинну генерируемой уникальной строки.
func NewShortLinkGenerator(length int) *ShortLinkGenerator {
	return &ShortLinkGenerator{length: length}
}

// UniqGenerate генерация уникальной стоки с фиксированной длинной.
func (generator *ShortLinkGenerator) UniqGenerate() string {
	return uniuri.NewLen(generator.length)
}
