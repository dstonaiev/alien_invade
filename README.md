# alien_invade
Test task


# Run application
go run cmd/main.go [flags]

Flags used in application run:
-A - number of generated values, optional, default 2
-M - path to map file, default data/map.txt
-S - true - print city map with aliens after each step, false - at the end only
-L - path to log file


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
Assumed, yes, as an option and in case no ways to go
