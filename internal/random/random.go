package random

import (
	"math/rand"
)

type Source struct {
	src rand.Source
}

func New(src rand.Source) *Source {
	return &Source{
		src: src,
	}
}

func (s *Source) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(s.src.Int63() & 0xff)
	}

	return len(p), nil
}
