package model

const stepsThreshold = 10000

type Alien struct {
	City         string
	stepsCounter int
}

func (a *Alien) Move(newCity string) {
	a.stepsCounter++
	a.City = newCity
}

func (a *Alien) ExceedThreshold() bool {
	return a.stepsCounter > stepsThreshold
}
