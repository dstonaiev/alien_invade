package model

import "github.com/dstonaiev/alien_invade/internal/rand"

type City struct {
	// city name
	Name string
	//alien which want to leave city, if no aliens is 0, can't be more than one, otherwise city should be destroyed on the previous step
	AlienOut uint
	//alien's numbers, if more than 1 city should be destroyed with aliens on it
	AliensIn []uint
	//map value = city name
	Destination map[Direction]string
}

func (c *City) AlienCome(alien uint) {
	c.AliensIn = append(c.AliensIn, alien)
}

func (c *City) DrawDirection(rnd rand.Randomizer) string {
	keys := make([]Direction, 0, len(c.Destination))
	for d := range c.Destination {
		keys = append(keys, d)
	}
	destNum := rnd.GenInt(len(keys) + 1)
	if destNum == len(keys) {
		return c.Name
	}
	return c.Destination[keys[destNum]]
}
