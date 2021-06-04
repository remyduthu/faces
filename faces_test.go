package faces

import (
	"testing"

	"github.com/go-playground/assert"
	"github.com/stretchr/testify/suite"
)

type facesTestSuite struct {
	suite.Suite
}

func (suite *facesTestSuite) TestShouldNotFilter() {
	var (
		actual = &struct {
			A string `faces:"public"`
			B string
		}{
			"a",
			"b",
		}
		expected = actual
	)

	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldNotFilterWithoutTags() {
	var (
		actual = &struct {
			A string
			B string
		}{
			"a",
			"b",
		}
		expected = actual
	)

	FilterWithTags(actual)

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilter() {
	var (
		actual = &struct {
			A string `faces:"private"`
			B string
		}{
			"a",
			"b",
		}
		expected = &struct {
			A string `faces:"private"`
			B string
		}{
			"",
			"b",
		}
	)

	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(facesTestSuite))
}
