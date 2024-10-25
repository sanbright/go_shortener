package generator

import "fmt"

func ExampleShortLinkGenerator_UniqGenerate() {
	var uniqString string

	// инициализируем генератор
	generator := NewShortLinkGenerator(10)

	// генерируем новую уникальную строку
	uniqString = generator.UniqGenerate()

	// 12kdaoije4
	fmt.Printf("%s", uniqString)
}

func ExampleCryptGenerator_EncodeValue() {
	var encoded string
	// инициализируем генератор
	crypt := NewCryptGenerator("r3pD#zq2*&nM979$3l$&!DcGi3piW3&%5wZ")

	// шифуем значение
	encoded, err := crypt.EncodeValue("This is test")

	if err != nil {
		panic(err)
	}

	//
	fmt.Printf("%s", encoded)
}
