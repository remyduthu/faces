package faces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type facesTestSuite struct {
	suite.Suite
}

type user struct {
	Email           string
	Password        string `faces:"private"`
	unexportedField string
}

func (suite *facesTestSuite) TestShouldNotFilter() {
	var (
		actual = &user{
			Email:           "foo&bar.com",
			Password:        "123",
			unexportedField: "Foo",
		}
		expected = actual
	)

	FilterWithTags(actual, "private")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldNotFilterSomethingElseThanAStructure() {
	var (
		actual = &map[string]string{
			"a": "b",
		}
		expected = actual
	)

	// Will do nothing
	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldNotFilterWithoutTags() {
	var (
		actual = &user{
			Email:    "foo&bar.com",
			Password: "123",
		}
		expected = actual
	)

	FilterWithTags(actual)

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilter() {
	var (
		actual = &user{
			Email:    "foo&bar.com",
			Password: "123",
		}
		expected = &user{
			Email:    "foo&bar.com",
			Password: "",
		}
	)

	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterMultipleTags() {
	var (
		actual = &user{
			Email:    "foo&bar.com",
			Password: "123",
		}
		expected = actual
	)

	FilterWithTags(actual, "public", "private")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterNestedStructure() {
	type admin struct {
		Name string
		Role string `faces:"private"`
		user
	}

	var (
		actual = &admin{
			Name: "Foo",
			Role: "Administrator",
			user: user{
				Email:    "foo&bar.com",
				Password: "123",
			},
		}
		expected = &admin{
			Name: "Foo",
			Role: "",
			user: user{
				Email:    "foo&bar.com",
				Password: "",
			},
		}
	)

	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterArrayOfStructures() {
	type team struct {
		Users []user
	}

	var (
		actual = &team{
			Users: []user{
				{
					Email:    "foo&bar.com",
					Password: "123",
				},
				{
					Email:    "john&doe.com",
					Password: "456",
				},
			},
		}
		expected = &team{
			Users: []user{
				{
					Email:    "foo&bar.com",
					Password: "",
				},
				{
					Email:    "john&doe.com",
					Password: "",
				},
			},
		}
	)

	FilterWithTags(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(facesTestSuite))
}
