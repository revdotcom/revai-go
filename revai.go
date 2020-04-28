package revai

import (
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.rev.ai/speechtotext/v1"
const defaultUserAgent = "revai-go-client"

type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
	UserAgent  string

	APIKey string
}

type ClientOption func(*Client)

// NewClient creates a new client and sets defaults. It then updates the client with any options passed in.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		HTTPClient: &http.Client{},
		BaseURL:    baseURL,
		APIKey:     apiKey,
		UserAgent:  defaultUserAgent,
	}

	for _, option := range opts {
		option(c)
	}

	return c
}

// HTTPClient sets the http client for the rev.ai client
func HTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// UserAgent sets the user agent for the rev.ai client
func UserAgent(userAgent string) func(*Client) {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// BaseURL sets the base url for the rev.ai client
func BaseURL(u *url.URL) func(*Client) {
	return func(c *Client) {
		c.BaseURL = u
	}
}
