package processor

import (
	"bytes"
	"log"
	"os"
	"testing"

	mock "github.com/dstonaiev/alien_invade/test/mock"
	"github.com/stretchr/testify/suite"
)

type appSuite struct {
	suite.Suite
	logger  *log.Logger
	app     *AlienInvasionApp
	workdir string
}

func TestAlienInvasionAppSuite(t *testing.T) {
	suite.Run(t, new(appSuite))
}

func (suite *appSuite) SetupTest() {
	suite.logger = log.New(os.Stdout, "testClient", 0)
	suite.workdir, _ = os.Getwd()
}

func (suite *appSuite) TestInitAppFileCorrupted() {
	var err error
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/corrupted.txt", mock.InitMockRandomizer())

	suite.Require().EqualError(err, "1 error occurred:\n\t* failed to parse line 'Bar south=Foo we'. reason: invalid separator = position\n\n")
}

func (suite *appSuite) TestSeedAliens() {
	var err error
	const numOfAliens = 5
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

	suite.app.SeedAliens(uint(numOfAliens))

	suite.Require().Len(suite.app.alienMap, numOfAliens)
	alienCnt := 0
	for _, city := range suite.app.cityMap {
		alienCnt += len(city.Aliens)
	}
	suite.Require().Equal(numOfAliens, alienCnt)
}

func (suite *appSuite) TestWalkCitiesOK() {
	var err error
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

	suite.app.SeedAliens(uint(10))

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().False(shouldStop)
}

func (suite *appSuite) TestWalkCitiesIsolatedOnly() {
	var err error
	const numOfAliens = 20
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapIsolated.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

	suite.app.SeedAliens(uint(numOfAliens))

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().True(shouldStop)
	suite.Require().Len(suite.app.alienMap, numOfAliens)
}

func (suite *appSuite) TestWalkCitiesAllCountersExceed() {
	var err error
	const numOfAliens = 10
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

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
	var err error
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().NoError(err)
}

func (suite *appSuite) TestValidateCityMapFileIncomplete() {
	var err error
	suite.app, err = InitApp(suite.logger, suite.workdir+"/../test/data/mapIncomplete.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

	err = suite.app.ValidateCityMap()

	suite.Require().Error(err)
}

func (suite *appSuite) TestUpdateCities() {
	var err error
	var buf bytes.Buffer
	logger := log.New(&buf, "Client Log: ", log.LstdFlags)
	suite.app, err = InitApp(logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.Require().NoError(err)

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

	suite.Require().Len(suite.app.cityMap, 5)

	_, ok := suite.app.cityMap["Kishore"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Alakam"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Tanzi"]
	suite.Require().False(ok)

	_, ok = suite.app.cityMap["Pandeva"]
	suite.Require().False(ok)

	c, ok := suite.app.cityMap["K'en'hahh"]
	suite.Require().True(ok)
	suite.Require().Len(c.Aliens, 1)
	suite.Require().Equal(c.Aliens[0], uint(2))

	c, ok = suite.app.cityMap["Basap"]
	suite.Require().True(ok)
	suite.Require().Len(c.Aliens, 1)
	suite.Require().Equal(c.Aliens[0], uint(6))

	c, ok = suite.app.cityMap["Umper"]
	suite.Require().True(ok)
	suite.Require().Len(c.Aliens, 1)
	suite.Require().Equal(c.Aliens[0], uint(10))
}
