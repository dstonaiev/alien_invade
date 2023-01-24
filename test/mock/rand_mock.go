package rand_test

import "github.com/dstonaiev/alien_invade/internal/rand"

func InitMockRandomizer() rand.Randomizer {
	return &MockRandomizer{
		counter: 0,
		valArr:  []int{213, 4354, 11, 435, 443, 11, 54, 45346, 65, 76, 657, 112, 78},
	}
}

type MockRandomizer struct {
	counter int
	valArr  []int
}

func (r *MockRandomizer) GenInt(rangeTo int) int {
	val := r.valArr[r.counter]
	r.counter++
	if r.counter == len(r.valArr) {
		r.counter = 0
	}
	return val % rangeTo
}
