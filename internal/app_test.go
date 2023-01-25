package processor

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	mock "github.com/dstonaiev/alien_invade/test/mock"
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

	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())
	suite.Require().EqualError(err, "1 error occurred:\n\t* failed to parse line 'Bar south=Foo we'. reason: invalid separator = position\n\n")
}

func (suite *appSuite) TestSeedAliens() {
	const numOfAliens = 5
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)

	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())
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

	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())
	suite.Require().NoError(err)

	for i := 1; i <= 4; i++ {
		suite.app.alienMap[uint(i)] = 0
	}

	suite.app.cityMap["Tanzi"].AlienOut = uint(1)
	suite.app.cityMap["Alakam"].AlienOut = uint(2)
	suite.app.cityMap["Pandeva"].AlienOut = uint(3)
	suite.app.cityMap["Umper"].AlienOut = uint(4)

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().False(shouldStop)
	suite.Require().Equal(uint(0), suite.app.cityMap["Tanzi"].AlienOut)
	suite.Require().Equal(uint(0), suite.app.cityMap["Alakam"].AlienOut)
	suite.Require().Equal(uint(0), suite.app.cityMap["Pandeva"].AlienOut)
	suite.Require().Equal(uint(0), suite.app.cityMap["Umper"].AlienOut)

	//1st alien moved from Tanzi to Pandeva
	suite.Require().Len(suite.app.cityMap["Pandeva"].AliensIn, 1)
	suite.Require().Equal(1, suite.app.cityMap["Pandeva"].AliensIn[0])

	//2nd alien moved from Alakam to Tanzi
	suite.Require().Len(suite.app.cityMap["Tanzi"].AliensIn, 1)
	suite.Require().Equal(2, suite.app.cityMap["Tanzi"].AliensIn[0])

	//3rd alien moved from Pandeva to Umper
	suite.Require().Len(suite.app.cityMap["Umper"].AliensIn, 1)
	suite.Require().Equal(3, suite.app.cityMap["Umper"].AliensIn[0])

	//4th alien moved from Umper to K'en'hahh
	suite.Require().Len(suite.app.cityMap["K'en'hahh"].AliensIn, 1)
	suite.Require().Equal(4, suite.app.cityMap["K'en'hahh"].AliensIn[0])
}

func (suite *appSuite) TestWalkCitiesIsolatedOnly() {
	const numOfAliens = 20
	path, err := filepath.Abs("../test/data/mapIsolated.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())

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
	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())

	suite.Require().NoError(err)

	suite.app.SeedAliens(numOfAliens)

	//init map with big counters
	for i := uint(1); i <= numOfAliens; i++ {
		suite.app.alienMap[i] = 10001
	}

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().True(shouldStop)
}

func (suite *appSuite) TestValidateCityMapFileOK() {
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().NoError(err)
}

func (suite *appSuite) TestValidateCityMapFileIncomplete() {
	path, err := filepath.Abs("../test/data/mapIncomplete.txt")
	suite.Require().NoError(err)
	suite.app, err = InitApp(suite.logger, path, mock.InitMockRandomizer())

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().Error(err)
}

func (suite *appSuite) TestUpdateCities() {
	var buf bytes.Buffer
	path, err := filepath.Abs("../test/data/mapOK.txt")
	suite.Require().NoError(err)
	logger := log.New(&buf, "Client Log: ", log.LstdFlags)
	suite.app, err = InitApp(logger, path, mock.InitMockRandomizer())

	suite.Require().NoError(err)

	for i := 1; i <= 12; i++ {
		suite.app.alienMap[uint(i)] = 0
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
	_, ok := suite.app.alienMap[1]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[2]
	suite.Require().True(ok)
	_, ok = suite.app.alienMap[3]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[4]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[5]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[6]
	suite.Require().True(ok)
	_, ok = suite.app.alienMap[7]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[8]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[9]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[10]
	suite.Require().True(ok)
	_, ok = suite.app.alienMap[11]
	suite.Require().False(ok)
	_, ok = suite.app.alienMap[12]
	suite.Require().False(ok)

	suite.Require().Len(suite.app.cityMap, 5)

	_, ok = suite.app.cityMap["Kishore"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Alakam"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Tanzi"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Pandeva"]
	suite.Require().False(ok)

	c, ok := suite.app.cityMap["K'en'hahh"]
	suite.Require().True(ok)
	suite.Require().Equal(c.AlienOut, uint(2))

	c, ok = suite.app.cityMap["Basap"]
	suite.Require().True(ok)
	suite.Require().Equal(c.AlienOut, uint(6))

	c, ok = suite.app.cityMap["Umper"]
	suite.Require().True(ok)
	suite.Require().Equal(c.AlienOut, uint(10))
}
