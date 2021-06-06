# 🌓 Faces

[![GoReference](https://pkg.go.dev/badge/github.com/remydt/faces.svg)](https://pkg.go.dev/github.com/remydt/faces)

Faces filters a structure using field tags. This is useful when you want to
create several versions of the same structure without duplicating it. A use case
could be when creating an API and not wanting to provide private information to
end users.

## Usage

You can find many examples in the [test file](./faces_test.go) but here is a
simple one:

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

Here are the results of a benchmark executed on my laptop. The results do not
represent an absolute truth but will allow me not to regress during future
updates.

```
goos: darwin
goarch: amd64
pkg: github.com/remydt/faces
cpu: Intel(R) Core(TM) i7-8559U CPU @ 2.70GHz
BenchmarkReveal-8                                  	 5044090	       243.2 ns/op	      40 B/op	       4 allocs/op
BenchmarkRevealWithSlice-8                         	 2272608	       515.0 ns/op	     104 B/op	       9 allocs/op
BenchmarkRevealWithMap-8                           	 1452736	       825.9 ns/op	     256 B/op	      12 allocs/op
BenchmarkRevealWithNestedStructure-8               	 2004470	       596.9 ns/op	     112 B/op	      10 allocs/op
BenchmarkRevealWithArrayOfStructures-8             	 1725184	       699.8 ns/op	     120 B/op	      11 allocs/op
BenchmarkRevealWithoutTags-8                       	537072181	         2.258 ns/op	       0 B/op	       0 allocs/op
BenchmarkRevealWithSomethingElseThanAStructure-8   	100000000	        10.50 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: 100.0% of statements
```
