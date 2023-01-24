package model

type City struct {
	// city name
	Name string
	//alien's numbers, if more than 1 city should be destroyed with aliens on it
	Aliens []uint
	//map value = city name
	Destination map[Direction]string
}

func (c *City) AlienCome(alien uint) {
	c.Aliens = append(c.Aliens, alien)
}

func (c *City) IsBattle() bool {
	return len(c.Aliens) > 1
}
