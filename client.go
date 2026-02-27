package pocketsmith

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

// An HttpClient is an interface over http.Client.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option configures a pocketsmith client.
type Option func(*Client) error

// Client defines a client for this package.
type Client struct {
	// config.
	baseURL    *url.URL // The endpoint to query against.
	token      string
	httpClient HttpClient // The http client used when sending / receiving data from the endpoint.

	// metadata.
	//authedUser *User // the authed user attached to the token.
}

func (c *Client) makeURL(p string) *url.URL {
	u := *c.baseURL
	u.Path = path.Join(u.Path, p)
	return &u
}

func (c *Client) decodeJSON(resp *http.Response, v interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Developer-Key", c.token)
	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

// New creates and returns a new Client, initialized with the provided token.
// The client itself is set up with tracing, logging, and HTTP configuration.
// Additional options can be provided to modify its behavior, via the options
// slice. The client is used for making requests and interacting with the
// Pockestmith API.
func New(token string, options ...Option) (*Client, error) {
	// check args.
	if token == "" {
		return nil, errors.New("pocketsmith: empty token")
	}

	u, err := url.Parse("https://api.pocketsmith.com/v2")
	if err != nil {
		return nil, err
	}

	// default client.
	c := &Client{
		token:      token,
		httpClient: http.DefaultClient,
		baseURL:    u,
	}

	// overwrite client with any given options.
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}
