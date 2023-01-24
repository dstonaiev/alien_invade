package processor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dstonaiev/alien_invade/internal/model"
	"github.com/dstonaiev/alien_invade/internal/rand"
	merr "github.com/hashicorp/go-multierror"
)

const (
	stepsThreshold = 10000
	dataSeparator  = " "
	destSeparator  = "="
)

func InitApp(logger *log.Logger, mapFile string, rnd rand.Randomizer) *AlienInvasionApp {
	app := &AlienInvasionApp{
		logger:   logger,
		cityMap:  make(map[string]*model.City),
		cityList: make([]string, 0),
		alienMap: make(map[uint]int, 0),
		rnd:      rnd,
	}
	app.initMap(mapFile)
	return app
}

type AlienInvasionApp struct {
	logger   *log.Logger
	cityMap  map[string]*model.City
	cityList []string
	alienMap map[uint]int
	rnd      rand.Randomizer
}

func (app *AlienInvasionApp) SeedAliens(aliensNum uint) {
	for i := uint(1); i <= aliensNum; i++ {
		app.alienMap[i] = 0
		cityKey := app.cityList[app.rnd.GenInt(len(app.cityList))]
		city := app.cityMap[cityKey]
		city.AlienCome(i)
	}
}

func (app *AlienInvasionApp) initMap(mapFile string) {
	absPath, err := filepath.Abs(mapFile)
	app.logger.Printf("Absolute path to map file: %s\n", absPath)
	if err != nil {
		app.logger.Fatalf("unable to evaluate file path %s. error %v", mapFile, err)
	}
	file, err := os.Open(absPath)
	if err != nil {
		app.logger.Fatalf("error opening file %s for city map. error %v", mapFile, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		city, err := app.addCity(scanner.Text())
		if err != nil {
			app.logger.Printf("Unable to parse city line, error %v", err)
		} else {
			app.cityMap[city.Name] = city
			app.cityList = append(app.cityList, city.Name)
		}
	}

	if err := scanner.Err(); err != nil {
		app.logger.Fatal(err)
	}
}

func (app *AlienInvasionApp) addCity(cityLine string) (*model.City, error) {
	parts := strings.Split(cityLine, dataSeparator)
	cityName := parts[0]
	_, ok := app.cityMap[cityName]
	if ok {
		return nil, fmt.Errorf("city with name %s was added to map previosly. dublicate will be ignored\n", cityName)
	}
	city := &model.City{
		Name:        cityName,
		Destination: make(map[model.Direction]string),
	}
	if len(parts) > 1 {
		for i := 1; i < len(parts); i++ {
			sepIdx := strings.Index(parts[i], destSeparator)
			dir, ok := model.ParseString(parts[i][:sepIdx])
			if !ok {
				return nil, fmt.Errorf("Invalid destination '%s' found in file", parts[i][:sepIdx])
			}
			city.Destination[dir] = parts[i][sepIdx+1:]
		}
	}
	return city, nil
}

func (app *AlienInvasionApp) WalkCities() bool {
	maxStepsForEveryone := true
	//check if any alien move may happen during walk stage. used to avoid "lost city only" trap
	noMove := true

	for _, cityKey := range app.cityList {
		city, ok := app.cityMap[cityKey]
		if ok && len(city.Aliens) > 0 && len(city.Destination) > 0 {
			noMove = false
			keys := getDirectionKeys(city)

			alien := city.Aliens[0]
			destNum := app.rnd.GenInt(len(keys) + 1)
			if destNum < len(keys) {
				//alien desided to leave city
				nextCity := city.Destination[keys[destNum]]
				//alien move
				app.alienMap[alien]++
				app.cityMap[nextCity].AlienCome(city.Aliens[0])
				city.Aliens = nil
			}
			maxStepsForEveryone = maxStepsForEveryone && (app.alienMap[alien] > stepsThreshold)
		}
	}
	return maxStepsForEveryone || noMove
}

func (app *AlienInvasionApp) IsEmpty() bool {
	return len(app.alienMap) <= 1 && len(app.cityMap) == 0
}

func (app *AlienInvasionApp) ValidateCityMap() error {
	var combErr error
	for name, city := range app.cityMap {
		for dest, val := range city.Destination {
			otherCity := app.cityMap[val]
			if otherCity == nil || otherCity.Destination[dest.GetOppos()] != name {
				combErr = merr.Append(fmt.Errorf("Invalid file parse: city's %s destionation %s should correspond to neighbor's city %s destination %v\n", name, dest.String(), val, dest.GetOppos()))
			}
		}
	}
	return combErr
}

func (app *AlienInvasionApp) UpdateCities() {
	for _, cityKey := range app.cityList {
		city, ok := app.cityMap[cityKey]
		if ok {
			if city.IsBattle() {
				app.logger.Printf("%s has been destroyed by ", cityKey)
				ln := len(city.Aliens)
				for i := 0; i < ln; i++ {
					app.logger.Printf("alien %d", city.Aliens[i])
					delete(app.alienMap, city.Aliens[i])
					if i == ln-2 {
						app.logger.Print(" and ")
					} else if i < ln-1 {
						app.logger.Print(", ")
					}
				}
				app.logger.Println()

				for dest, neighCity := range city.Destination {
					delete(app.cityMap[neighCity].Destination, dest.GetOppos())
				}
				delete(app.cityMap, cityKey)
			}
		}
	}
}

func (app *AlienInvasionApp) PrintResult() {
	if len(app.cityMap) == 0 {
		app.logger.Println("World X was fully destroyed")
	} else {
		app.logger.Println("Remaining cities:")
		for name, city := range app.cityMap {
			app.logger.Print(name)
			if len(city.Aliens) > 0 {
				app.logger.Printf(" aliens: %+q", city.Aliens)
			}
			if len(city.Destination) == 0 {
				app.logger.Print(" LOST")
			}
			app.logger.Println()
		}
	}
}

func getDirectionKeys(city *model.City) []model.Direction {
	keys := make([]model.Direction, 0, len(city.Destination))
	for d := range city.Destination {
		keys = append(keys, d)
	}
	return keys
}
