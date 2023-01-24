package rand

import (
	"math/rand"
	"time"
)

type Randomizer interface {
	GenInt(int) int
}

func NewRandomizer() Randomizer {
	src := rand.NewSource(time.Now().UnixNano())
	return &randomizer{
		rnd: rand.New(src)}
}

type randomizer struct {
	rnd *rand.Rand
}

func (r *randomizer) GenInt(rangeTo int) int {
	return r.rnd.Intn(rangeTo)
}
