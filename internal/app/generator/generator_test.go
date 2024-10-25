package generator

import "testing"

func BenchmarkUniqGenerate(b *testing.B) {
	b.Run("generate 10 length", func(b *testing.B) {
		generator := NewShortLinkGenerator(10)

		for i := 0; i < b.N; i++ {
			generator.UniqGenerate()
		}
	})

	b.Run("generate 40 length", func(b *testing.B) {
		generator := NewShortLinkGenerator(40)

		for i := 0; i < b.N; i++ {
			generator.UniqGenerate()
		}
	})
}
