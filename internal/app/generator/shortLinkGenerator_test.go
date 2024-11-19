package generator

import (
	"testing"
)

func TestShortLinkGenerator_UniqGenerate(t *testing.T) {

	type want struct {
		lenght int
	}

	tests := []struct {
		name   string
		lenght int
		want   want
	}{
		{
			name:   "UniqGenerate_5",
			lenght: 5,
			want: want{
				lenght: 5,
			},
		},
		{
			name:   "UniqGenerate_10",
			lenght: 10,
			want: want{
				lenght: 10,
			},
		},
		{
			name:   "UniqGenerate_30",
			lenght: 30,
			want: want{
				lenght: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewShortLinkGenerator(tt.lenght)

			code := generator.UniqGenerate()

			if len(code) != tt.want.lenght {
				t.Errorf("%v: lenght_generated_code = '%v', want = '%v'", tt.name, len(code), tt.want.lenght)
			}
		})
	}
}
