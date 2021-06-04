# 🌓 Faces

[![Go Reference](https://pkg.go.dev/badge/github.com/remydt/faces.svg)](https://pkg.go.dev/github.com/remydt/faces)

Faces filters a structure using field tags. This is useful when you want to
create several versions of the same structure without duplicating it. A use case
could be when creating an API and not wanting to provide private information to
end users.

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
