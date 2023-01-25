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

func InitApp(logger *log.Logger, mapFile string, rnd rand.Randomizer) (*AlienInvasionApp, error) {
	app := &AlienInvasionApp{
		logger:   logger,
		cityMap:  make(map[string]*model.City),
		cityList: make([]string, 0),
		alienMap: make(map[uint]int, 0),
		rnd:      rnd,
	}
	err := app.initMap(mapFile)
	return app, err
}

type AlienInvasionApp struct {
	logger   *log.Logger
	cityMap  map[string]*model.City
	cityList []string
	alienMap map[uint]int
	rnd      rand.Randomizer
}

func (app *AlienInvasionApp) SeedAliens(aliensNum uint) {
	cityListLen := len(app.cityList)
	for i := uint(1); i <= aliensNum; i++ {
		app.alienMap[i] = 0
		cityKey := app.cityList[app.rnd.GenInt(cityListLen)]
		city := app.cityMap[cityKey]
		//not needed to check precense as it first step and cityList fully matches to keys in cityMap at this stage (no map entries deleted yet)
		city.AlienCome(i)
	}
}

func (app *AlienInvasionApp) initMap(mapFile string) error {
	absPath, err := filepath.Abs(mapFile)
	app.logger.Printf("Absolute path to map file: %s\n", absPath)
	if err != nil {
		return fmt.Errorf("unable to evaluate file path %s. error %v", mapFile, err)
	}
	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("error opening file %s for city map. error %v", mapFile, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var combErr error
	for scanner.Scan() {
		city, err := app.addCity(scanner.Text())
		if err != nil {
			combErr = merr.Append(combErr, err)
		} else {
			app.cityMap[city.Name] = city
			app.cityList = append(app.cityList, city.Name)
		}
	}
	if combErr != nil {
		return combErr
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (app *AlienInvasionApp) addCity(cityLine string) (*model.City, error) {
	parts := strings.Split(cityLine, dataSeparator)
	cityName := parts[0]
	_, ok := app.cityMap[cityName]
	if ok {
		return nil, fmt.Errorf("city with name %s was added to map previosly. dublicate will be ignored", cityName)
	}
	city := &model.City{
		Name:        cityName,
		Destination: make(map[model.Direction]string),
	}
	if len(parts) > 1 {
		for i := 1; i < len(parts); i++ {
			sepIdx := strings.Index(parts[i], destSeparator)
			if sepIdx == 4 || sepIdx == 5 { //lenght of direction words either 4 or 5
				dir, ok := model.ParseString(parts[i][:sepIdx])
				if !ok {
					return nil, fmt.Errorf("invalid destination '%s' found in file", parts[i][:sepIdx])
				}
				city.Destination[dir] = parts[i][sepIdx+1:]
			} else {
				return nil, fmt.Errorf("failed to parse line '%s'. reason: invalid separator %s position", cityLine, destSeparator)
			}
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

			//singe actor zone
			nextCityKey := city.DrawDirection(app.rnd)

			alien := city.Aliens[0]
			if nextCityKey != cityKey {
				//alien desided to leave city
				app.alienMap[alien]++
				app.cityMap[nextCityKey].AlienCome(alien)
				city.Aliens = nil
			}
			//end of singe actor zone

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
		for dest, cityName := range city.Destination {
			if name == cityName {
				combErr = merr.Append(combErr, fmt.Errorf("invalid file parse: city's %s destination %v should not point to itself", name, dest))
			} else {
				otherCity, ok := app.cityMap[cityName]
				if !ok || otherCity.Destination[dest.GetOppos()] != name {
					combErr = merr.Append(combErr, fmt.Errorf("invalid file parse: city's %s destination %v should correspond to neighbor's city %s destination %v", name, dest, cityName, dest.GetOppos()))
				}
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
