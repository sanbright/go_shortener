package service

import "github.com/dchest/uniuri"

func UniqGenerate() string {
	return uniuri.NewLen(10)
}
