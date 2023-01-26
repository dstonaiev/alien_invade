# Summary
Alien invasion simulator

# Description
Mad aliens are about to invade the earth and you are tasked with simulating the invasion.
Given a map containing the names of cities in the non-existent world of X. The map is in the file, with one city per line.
The city name is first, followed by 1-4 directions (north, south, east, or west). Each one represents a road to another city that lies in that direction.

The aliens (generated according to number given in command line) start out at random places on the map, and wander around randomly, following links. Each iteration, the aliens can travel in any of the directions leading out of a city.

When few aliens end up in the same place, they fight, and int the process kill each other and destroy the city. When a city is destroyed, it is removed from the map, and so are any roads that lead into or out of it.

# Run application
go run cmd/main.go [flags]

Flags used in application run:
-A - number of generated aliens, optional, default 2
-M - path to map file, default data/map.txt
-S - true - print city map with aliens after each step, false - at the end only
-L - path to log folder

# Q & A
1) May I assume city path bidirectional, i.e. if city A contains path to B, so city B definitely contains path to A?
Assumed bidirectional, for this purpose data validation applied

2) What happens with city and alien if only one alien in city and no ways to go (neighbors cities was destroyed)?
Alien remains in the city till the end of game. City will be ignored for steps

3) Will all aliens kill each other If by chance more than 2 aliens appeared in the city in the same time? 
Assumed, everyone will be killed in this case

4) How big data size may be? Does multithreading approach expected or linear can work as well?
Linear approach selected as no much heavy data operations expected in single processing, no much time gain in multithreading, however possible conflicts, "lost" aliens, which leave city and can't come to the next one because it destroyed

5) May alien prefer do not move to another city?
Assumed, yes, as an option (in this case step counter increased) and in case no ways to go (no step counter increasing)
