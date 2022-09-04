//go:generate go run ../cmd/subcmd/examples.go ../cmd/subcmd/doc.go examples
package api

import "github.com/gofrs/uuid"

// Logging a user with username and password in.
var LoginRequestExample1 = LoginRequest{
	Username: "chad",
	Password: "secret",
}

// Successful login of a new user.
var LoginResponseExample1 = LoginResponse{
	SessionID: uuid.Nil,
}
