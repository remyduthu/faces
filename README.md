# 🌓 Faces

[![Go Reference](https://pkg.go.dev/badge/github.com/remydt/faces.svg)](https://pkg.go.dev/github.com/remydt/faces)

Faces filters a structure using field tags. This is useful when you want to
create several versions of the same structure without duplicating it. A use case
could be when creating an API and not wanting to provide private information to
end users.

## Usage

A simple example can be:

```go
// A simple user structure
type User struct {
  Email    string
  Password string `faces:"private"` // The password is part of the user's private face.
}

// Create a user
user := &User{
  Email:    "foo@bar.com",
  Password: "private-password",
}

// Filter the user
faces.Reveal(user, "public")

fmt.Printf("%#v", user)
// Output:
// &main.User{Email:"foo@bar.com", Password:""}
```

As you can see, the password is now hidden (or at least, reset). This package
can be combined with the JSON package as in the example below:

```go
// A simple user structure
type User struct {
  Email    string `json:"email"`
  Password string `faces:"private" json:"password,omitempty"` // We do not want to display the field when using json
}

// Create a user
user := &User{
  Email:    "foo@bar.com",
  Password: "private-password",
}

// Filter the user
faces.Reveal(user, "public")

// Convert the user to JSON
b, _ := json.Marshal(user)

fmt.Println(string(b))
// Output:
// {"email":"foo@bar.com"}
```

## Benchmarks

Here are the results of a benchmark executed on my laptop. The results do not represent an absolute truth but will allow me not to regress during future updates.

```
goos: darwin
goarch: amd64
pkg: github.com/remydt/faces
cpu: Intel(R) Core(TM) i7-8559U CPU @ 2.70GHz
BenchmarkFaces-8                                  	 5030667	       231.3 ns/op	      40 B/op	       4 allocs/op
BenchmarkFacesWithNestedStructure-8               	 2031132	       599.7 ns/op	     112 B/op	      10 allocs/op
BenchmarkFacesWithArrayOfStructures-8             	 1724448	       686.8 ns/op	     120 B/op	      11 allocs/op
BenchmarkFacesWithoutTags-8                       	534643616	         2.242 ns/op	       0 B/op	       0 allocs/op
BenchmarkFacesWithSomethingElseThanAStructure-8   	129993452	         9.237 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: 100.0% of statements
```
