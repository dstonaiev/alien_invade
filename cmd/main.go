package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	processor "github.com/dstonaiev/alien_invade/internal"
)

var (
	aliensNum     uint
	mapFile       string
	logFile       string
	printEachStep bool
)

func init() {
	flag.UintVar(&aliensNum, "A", 2, "aliens number")
	flag.StringVar(&mapFile, "M", "data/map.txt", "file path use as world map")
	flag.BoolVar(&printEachStep, "S", false, "true - print city map with aliens after each step")
	flag.StringVar(&logFile, "L", "log", "log file")
	flag.Parse()
}

func main() {
	logger := initLog(logFile)
	app, err := processor.InitApp(logger, mapFile)
	if err != nil {
		log.Panicf("provided map file is corrupted, errors: %v", err)
	}

	if err = app.ValidateCityMap(); err != nil {
		log.Panicf("provided map file didn't pass validation, errors: %v", err)
	}

	if aliensNum == uint(0) {
		log.Panic("No aliens were generated")
	}

	// seed aliens
	app.SeedAliens(aliensNum)

	//phase result
	app.UpdateCities()

	stop := false

	//start invasion
	//check condition to continue process:
	//1. Some cities exist with non-empty ways to reach them
	//2. At least to aliens exists, if one only, he can't change map
	//3. At leat one alien made less than 10000 steps
	for !stop && !app.IsEmpty() {
		stop = true //should reset stop value on each iteration
		stop = stop && app.WalkCities()

		//phase result
		if !stop {
			app.UpdateCities()
		}
		if printEachStep {
			app.PrintResult()
		}
	}

	//Print result
	if !printEachStep {
		app.PrintResult()
	}
}

func initLog(logPath string) *log.Logger {
	absPath, err := filepath.Abs(logPath + "/log" + time.Now().Format(time.RFC3339) + ".log")
	log.Printf("Absolute path to map file: %s\n", absPath)
	if err != nil {
		log.Printf("unable to evaluate file path %s. error %v", logFile, err)
		return log.New(os.Stdout, "App Log: ", log.LstdFlags)
	}
	file, err := os.Create(absPath)
	if err != nil {
		log.Printf("error opening file %s for log dump. error %v", logFile, err)
		return log.New(os.Stdout, "App Log: ", log.LstdFlags)
	}
	return log.New(file, "App Log: ", log.LstdFlags)
}
