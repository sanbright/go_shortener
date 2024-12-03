package generator

import "testing"

func TestCryptGenerator_EncodeValue(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "UniqGenerate_5",
			text: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator := NewCryptGenerator("$$ecuRityKe453H@")

			code, err := generator.EncodeValue(tt.text)

			if err != nil {
				t.Error(err)
			}

			decode, err := generator.DecodeValue(code)

			if err != nil {
				t.Error(err)
			}

			if tt.text != decode {
				t.Errorf("%v: encode and decode failed. source = '%v', want = '%v'", tt.name, tt.text, decode)
			}
		})
	}
}
