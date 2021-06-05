package faces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type facesTestSuite struct {
	suite.Suite
}

type admin struct {
	Name string
	Role string `faces:"private"`
	user
}

type team struct {
	Users []user
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

	Reveal(actual, "private")

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
	Reveal(actual, "public")

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

	Reveal(actual)

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldPanicWithoutAnAddress() {
	var (
		u = user{
			Email:    "foo&bar.com",
			Password: "123",
		}
	)

	assert.Panics(suite.T(), func() {
		Reveal(u, "public")
	})
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

	Reveal(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterMultipleTags() {
	type identity struct {
		FirstName string `faces:"private,public"`
		LastName  string `faces:"confidential,private"`
		BirthDate string `faces:"confidential"`
	}

	var (
		actual = &identity{
			FirstName: "Foo",
			LastName:  "Bar",
			BirthDate: "01/01/2000",
		}
		expected = &identity{
			FirstName: "Foo",
			LastName:  "Bar",
			BirthDate: "",
		}
	)

	Reveal(actual, "public", "private")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterNestedStructure() {

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

	Reveal(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldFilterArrayOfStructures() {
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

	Reveal(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(facesTestSuite))
}

func BenchmarkFaces(b *testing.B) {
	u := &user{
		Email:    "foo&bar.com",
		Password: "123",
	}

	for n := 0; n < b.N; n++ {
		Reveal(u, "private")
	}
}

func BenchmarkFacesWithNestedStructure(b *testing.B) {
	a := &admin{
		Name: "Foo",
		Role: "Administrator",
		user: user{
			Email:    "foo&bar.com",
			Password: "123",
		},
	}

	for n := 0; n < b.N; n++ {
		Reveal(a, "public")
	}
}

func BenchmarkFacesWithArrayOfStructures(b *testing.B) {
	t := &team{
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

	for n := 0; n < b.N; n++ {
		Reveal(t, "public")
	}
}

func BenchmarkFacesWithoutTags(b *testing.B) {
	u := &user{
		Email:    "foo&bar.com",
		Password: "123",
	}

	for n := 0; n < b.N; n++ {
		Reveal(u)
	}
}

func BenchmarkFacesWithSomethingElseThanAStructure(b *testing.B) {
	m := &map[string]string{
		"a": "b",
	}

	for n := 0; n < b.N; n++ {
		// Will do nothing
		Reveal(m, "public")
	}
}
