package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type citySuite struct {
	suite.Suite
}

func TestCitySuite(t *testing.T) {
	suite.Run(t, new(citySuite))
}

func (suite *citySuite) SetupTest() {
}

func (suite *citySuite) TestAlienCome() {

}

func (suite *citySuite) TestDrawDirection() {

}
