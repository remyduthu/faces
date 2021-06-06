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

func (suite *facesTestSuite) TestShouldNotReveal() {
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

func (suite *facesTestSuite) TestShouldNotRevealSomethingElseThanAStructure() {
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

func (suite *facesTestSuite) TestShouldNotRevealWithoutTags() {
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

func (suite *facesTestSuite) TestShouldReveal() {
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

func (suite *facesTestSuite) TestShouldRevealSlice() {
	var (
		actual = []*user{
			{
				Email:    "foo&bar.com",
				Password: "123",
			},
			{
				Email:    "john&doe.com",
				Password: "456",
			},
		}
		expected = []*user{
			{
				Email:    "foo&bar.com",
				Password: "",
			},
			{
				Email:    "john&doe.com",
				Password: "",
			},
		}
	)

	Reveal(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldRevealMapValues() {
	var (
		actual = map[string]*user{
			"foo": {
				Email:    "foo&bar.com",
				Password: "123",
			},
			"bar": {
				Email:    "john&doe.com",
				Password: "456",
			},
		}
		expected = map[string]*user{
			"foo": {
				Email:    "foo&bar.com",
				Password: "",
			},
			"bar": {
				Email:    "john&doe.com",
				Password: "",
			},
		}
	)

	Reveal(actual, "public")

	assert.Equal(suite.T(), expected, actual)
}

func (suite *facesTestSuite) TestShouldRevealMultipleTags() {
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

func (suite *facesTestSuite) TestShouldRevealNestedStructure() {

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

func (suite *facesTestSuite) TestShouldRevealArrayOfStructures() {
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

func BenchmarkReveal(b *testing.B) {
	u := &user{
		Email:    "foo&bar.com",
		Password: "123",
	}

	for n := 0; n < b.N; n++ {
		Reveal(u, "private")
	}
}

func BenchmarkRevealWithSlice(b *testing.B) {
	users := []user{
		{
			Email:    "foo&bar.com",
			Password: "123",
		},
		{
			Email:    "john&doe.com",
			Password: "456",
		},
	}

	for n := 0; n < b.N; n++ {
		Reveal(users, "private")
	}
}

func BenchmarkRevealWithMap(b *testing.B) {
	users := map[string]*user{
		"foo": {
			Email:    "foo&bar.com",
			Password: "123",
		},
		"bar": {
			Email:    "john&doe.com",
			Password: "456",
		},
	}

	for n := 0; n < b.N; n++ {
		Reveal(users, "private")
	}
}

func BenchmarkRevealWithNestedStructure(b *testing.B) {
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

func BenchmarkRevealWithArrayOfStructures(b *testing.B) {
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

func BenchmarkRevealWithoutTags(b *testing.B) {
	u := &user{
		Email:    "foo&bar.com",
		Password: "123",
	}

	for n := 0; n < b.N; n++ {
		Reveal(u)
	}
}

func BenchmarkRevealWithSomethingElseThanAStructure(b *testing.B) {
	m := &map[string]string{
		"a": "b",
	}

	for n := 0; n < b.N; n++ {
		// Will do nothing
		Reveal(m, "public")
	}
}
