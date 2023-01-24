package model

import (
	"strings"
)

var (
	directionsMap = map[string]Direction{
		"north": North,
		"south": South,
		"east":  East,
		"west":  West,
	}
)

func ParseString(str string) (Direction, bool) {
	c, ok := directionsMap[strings.ToLower(str)]
	return c, ok
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func (d *Direction) String() string {
	switch *d {
	case North:
		return "north"
	case South:
		return "south"
	case East:
		return "east"
	case West:
		return "west"
	}
	return "unknown"
}

func (d *Direction) GetOppos() Direction {
	switch *d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	}
	return *d
}
