//go:generate go run ../cmd/apigen.go
//go:generate go run ../cmd/subcmd/doc.go ../cmd/subcmd/examples.go
package api

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/gofrs/uuid"
)

const (
	LoginEndpoint      = "/login"
	CreatePostEndpoint = "/post"
	GetFeedEndpoint    = "/feed"
	FollowUserEndpoint = "/follow"
)

type API interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	CreatePost(ctx context.Context, req CreatePostRequest) (*CreatePostResponse, error)
	GetFeed(ctx context.Context, req GetFeedRequest) (*GetFeedResponse, error)
	FollowUser(ctx context.Context, req FollowUserRequest) (*FollowUserResponse, error)
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
	// SessionID identifies the session cookie as UUID.
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
	// Body is the post's content written in Markdown syntax.
	Body string `json:"body"`
	// URL associates the post with a link.
	URL string `json:"url"`
}

func (r *CreatePostRequest) Method() string {
	return http.MethodPost
}

type CreatePostResponse struct {
	// ID identifies the newly created post.
	ID int `json:"id"`
}

// CreatePost creates a new post and enters it to th feeds of the users following the author.
func (a *api) CreatePost(ctx context.Context, req CreatePostRequest) (*CreatePostResponse, error) {
	return &CreatePostResponse{
		ID: rand.Int(),
	}, nil
}

type GetFeedRequest struct {
	// PerPage specifies the number of posts returned per request.
	PerPage int `json:"pageSize"`
	// Page for pagination; the page to retrieve.
	Page int `json:"pageNumber"`
}

func (r *GetFeedRequest) Method() string {
	return http.MethodGet
}

type GetFeedResponse struct {
	// Posts is the lists of all posts.
	Posts []string `json:"posts"`
}

// GetFeed retrieves the feed for the given user.
func (a *api) GetFeed(ctx context.Context, req GetFeedRequest) (*GetFeedResponse, error) {
	return &GetFeedResponse{
		Posts: []string{},
	}, nil
}

type FollowUserRequest struct {
	UserID int `json:"userID"`
}

func (r *FollowUserRequest) Method() string {
	return http.MethodGet
}

type FollowUserResponse struct {
}

func (a *api) FollowUser(ctx context.Context, req FollowUserRequest) (*FollowUserResponse, error) {
	return nil, nil
}
