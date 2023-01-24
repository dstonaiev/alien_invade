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

func (suite *appSuite) TestSeedAliens() {
	const numOfAliens = 5
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.app.SeedAliens(uint(numOfAliens))

	suite.Require().Len(suite.app.alienMap, numOfAliens)
	alienCnt := 0
	for _, city := range suite.app.cityMap {
		alienCnt += len(city.Aliens)
	}
	suite.Require().Equal(numOfAliens, alienCnt)
}

func (suite *appSuite) TestWalkCitiesOK() {
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().False(shouldStop)
}

func (suite *appSuite) TestWalkCitiesIsolatedOnly() {
	const numOfAliens = 20
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/mapIsolated.txt", mock.InitMockRandomizer())

	suite.app.SeedAliens(uint(numOfAliens))

	//city not found in map
	shouldStop := suite.app.WalkCities()

	suite.Require().True(shouldStop)
	suite.Require().Len(suite.app.alienMap, numOfAliens)
}

func (suite *appSuite) TestWalkCitiesAllCountersExceed() {
	const numOfAliens = 10
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/test4.txt", mock.InitMockRandomizer())
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
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	err := suite.app.ValidateCityMap()

	suite.Require().NoError(err)
}

func (suite *appSuite) TestValidateCityMapFileIncomplete() {
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	err := suite.app.ValidateCityMap()

	suite.Require().EqualError(err, "")
}

func (suite *appSuite) TestValidateCityMapFileCorrupted() {
	suite.app = InitApp(suite.logger, suite.workdir+"/../test/data/corrupted.txt", mock.InitMockRandomizer())

	err := suite.app.ValidateCityMap()

	suite.Require().EqualError(err, "")
}

func (suite *appSuite) TestUpdateCities() {
	var buf bytes.Buffer
	logger := log.New(&buf, "Client Log: ", log.LstdFlags)
	suite.app = InitApp(logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.app.UpdateCities()
	//delete some aliens
	//delete some cities

	//check updated map

	suite.Require().Equal("", buf.String())
}

func (suite *appSuite) TestPrintResult() {
	var buf bytes.Buffer
	logger := log.New(&buf, "Client Log: ", log.LstdFlags)
	suite.app = InitApp(logger, suite.workdir+"/../test/data/mapOK.txt", mock.InitMockRandomizer())

	suite.app.PrintResult()

	suite.Require().Equal("", buf.String())
}
