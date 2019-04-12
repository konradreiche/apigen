//go:generate go run ../cmd/apigen.go
//go:generate go run ../cmd/subcmd/doc.go ../cmd/subcmd/examples.go
package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/gobuffalo/uuid"
)

const (
	LoginEndpoint      = "/login"
	CreatePostEndpoint = "/post"
	GetFeedEndpoint    = "/feed"
)

type API interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	CreatePost(ctx context.Context, req CreatePostRequest) (*CreatePostResponse, error)
	GetFeed(ctx context.Context, req GetFeedRequest) (*GetFeedResponse, error)
}

type api struct {
}

func NewAPI() API {
	return &api{}
}

type LoginRequest struct {
	// Username is used for logging the user in.
	Username string `json:"username"`
	// Password is used for logging the user in.
	Password string `json:"password"`
}

func (r *LoginRequest) Method() string {
	return http.MethodPost
}

type LoginResponse struct {
	SessionID uuid.UUID `json:"id"`
}

// Login authenticates a user and generates a session ID returned to the client.
func (a *api) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		SessionID: id,
	}, nil
}

type CreatePostRequest struct {
	Body string `json:"body"`
	URL  string `json:"url"`
}

func (r *CreatePostRequest) Method() string {
	return http.MethodPost
}

type CreatePostResponse struct {
	ID int `json:"id"`
}

func (a *api) CreatePost(ctx context.Context, req CreatePostRequest) (*CreatePostResponse, error) {
	return &CreatePostResponse{
		ID: rand.Int(),
	}, nil
}

type GetFeedRequest struct {
}

func (r *GetFeedRequest) Method() string {
	return http.MethodGet
}

type GetFeedResponse struct {
	Posts []string `json:"posts"`
}

func (a *api) GetFeed(ctx context.Context, req GetFeedRequest) (*GetFeedResponse, error) {
	return &GetFeedResponse{
		Posts: []string{},
	}, nil
}
