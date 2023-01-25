package processor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"github.com/dstonaiev/alien_invade/internal/model"
	"github.com/stretchr/testify/suite"
)

type appSuite struct {
	suite.Suite
	logger *log.Logger
	app    *AlienInvasionApp
}

func TestAlienInvasionAppSuite(t *testing.T) {
	suite.Run(t, new(appSuite))
}

func (suite *appSuite) SetupTest() {
	suite.logger = log.New(os.Stdout, "testClient", 0)
}

func (suite *appSuite) TestInitAppFileCorrupted() {
	path, err := filepath.Abs("../test/data/corrupted.txt")
	suite.Require().NoError(err)

	suite.app, err = InitApp(suite.logger, path)
	suite.Require().EqualError(err, "1 error occurred:\n\t* failed to parse line 'Bar south=Foo we'. reason: invalid separator = position\n\n")
}

func (suite *appSuite) TestSeedAliens() {
	const numOfAliens = 5
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)

	suite.app, err = InitApp(suite.logger, path)
	suite.Require().NoError(err)

	suite.app.SeedAliens(uint(numOfAliens))

	suite.Require().Len(suite.app.alienMap, numOfAliens)
	alienCnt := 0
	for _, city := range suite.app.cityMap {
		alienCnt += len(city.AliensIn)
	}
	suite.Require().Equal(numOfAliens, alienCnt)
}

func (suite *appSuite) TestWalkCitiesOK() {
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)

	suite.app, err = InitApp(suite.logger, path)
	suite.Require().NoError(err)

	suite.app.alienMap[uint(0)] = &model.Alien{City: "Tanzi"}
	suite.app.alienMap[uint(1)] = &model.Alien{City: "Alakam"}
	suite.app.alienMap[uint(2)] = &model.Alien{City: "Pandeva"}
	suite.app.alienMap[uint(3)] = &model.Alien{City: "Umper"}

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().False(shouldStop)

	suite.Require().True(checkContainValue(suite.app.cityMap[suite.app.alienMap[uint(0)].City], "Tanzi"))
	suite.Require().True(checkContainValue(suite.app.cityMap[suite.app.alienMap[uint(1)].City], "Alakam"))
	suite.Require().True(checkContainValue(suite.app.cityMap[suite.app.alienMap[uint(2)].City], "Pandeva"))
	suite.Require().True(checkContainValue(suite.app.cityMap[suite.app.alienMap[uint(3)].City], "Umper"))
}

func (suite *appSuite) TestWalkCitiesIsolatedOnly() {
	const numOfAliens = 20
	path, err := filepath.Abs("../test/data/mapIsolated.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path)

	suite.Require().NoError(err)

	suite.app.SeedAliens(uint(numOfAliens))

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().True(shouldStop)
	suite.Require().Len(suite.app.alienMap, numOfAliens)
}

func (suite *appSuite) TestWalkCitiesAllCountersExceed() {
	const numOfAliens = 10
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path)

	suite.Require().NoError(err)

	suite.app.SeedAliens(numOfAliens)

	//init map with big counters
	for i := uint(1); i <= numOfAliens; i++ {
		setFieldValue(suite.app.alienMap[i], "stepsCounter", 10001)
	}

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().True(shouldStop)
}

func (suite *appSuite) TestValidateCityMapFileOK() {
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path)

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().NoError(err)
}

func (suite *appSuite) TestValidateCityMapFileIncomplete() {
	path, err := filepath.Abs("../test/data/mapIncomplete.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path)

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().Error(err)
}

func (suite *appSuite) TestUpdateCities() {
	var buf bytes.Buffer
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)
	logger := log.New(&buf, "Client Log: ", log.LstdFlags)
	suite.app, err = InitApp(logger, path)

	suite.Require().NoError(err)

	for i := 1; i <= 12; i++ {
		suite.app.alienMap[uint(i)] = &model.Alien{}
	}
	suite.app.cityMap["Tanzi"].AlienCome(uint(1))
	suite.app.cityMap["Tanzi"].AlienCome(uint(5))
	suite.app.cityMap["Tanzi"].AlienCome(uint(7))

	suite.app.cityMap["Kishore"].AlienCome(uint(3))
	suite.app.cityMap["Kishore"].AlienCome(uint(8))

	suite.app.cityMap["K'en'hahh"].AlienCome(uint(2))

	suite.app.cityMap["Alakam"].AlienCome(uint(9))
	suite.app.cityMap["Alakam"].AlienCome(uint(11))

	suite.app.cityMap["Pandeva"].AlienCome(uint(4))
	suite.app.cityMap["Pandeva"].AlienCome(uint(12))

	suite.app.cityMap["Basap"].AlienCome(uint(6))

	suite.app.cityMap["Umper"].AlienCome(uint(10))

	suite.app.UpdateCities()

	suite.Require().Len(suite.app.alienMap, 3)
	_, ok := suite.app.alienMap[2]
	suite.Require().True(ok)
	_, ok = suite.app.alienMap[6]
	suite.Require().True(ok)
	_, ok = suite.app.alienMap[10]
	suite.Require().True(ok)

	suite.Require().Len(suite.app.cityMap, 5)

	_, ok = suite.app.cityMap["Kishore"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Alakam"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Tanzi"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Pandeva"]
	suite.Require().False(ok)
}

func checkContainValue(city *model.City, dest string) bool {
	if city.Name == dest {
		return true
	}
	for _, cityName := range city.Destination {
		if dest == cityName {
			return true
		}
	}
	return false
}

func setFieldValue(target any, fieldName string, value any) {
	rv := reflect.ValueOf(target)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	if !rv.CanAddr() {
		panic("target must be addressable")
	}
	if rv.Kind() != reflect.Struct {
		panic(fmt.Sprintf(
			"unable to set the '%s' field value of the type %T, target must be a struct",
			fieldName,
			target,
		))
	}
	rf := rv.FieldByName(fieldName)

	reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}
