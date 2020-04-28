package revai

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testUserAgent = "test-user-agent"
	testBaseURL   = "http://test.com"
)

func TestNewClient(t *testing.T) {
	c := NewClient("API-KEY")

	assert.NotNil(t, c)
}

func TestNewClientWithOptions(t *testing.T) {
	testHttpClient := &http.Client{
		Timeout: 1 * time.Second,
	}

	u, err := url.Parse(testBaseURL)
	if err != nil {
		t.Error(err)
	}

	c := NewClient(
		"API-KEY",
		HTTPClient(testHttpClient),
		UserAgent("test-user-agent"),
		BaseURL(u),
	)

	assert.NotNil(t, c)
	assert.Equal(t, c.UserAgent, testUserAgent, "user agent should match passed in user agent")
	assert.Equal(t, c.BaseURL, u, "base url should match passed in base url")
	assert.Equal(t, c.HTTPClient, testHttpClient, "http client should match passed in http client")
}

func TestUserAgent(t *testing.T) {
	c := NewClient(
		"api-key",
		UserAgent(testUserAgent),
	)

	assert.Equal(t, c.UserAgent, testUserAgent, "user agent should match passed in user agent")
}

func TestBaseURL(t *testing.T) {
	u, err := url.Parse(testBaseURL)
	if err != nil {
		t.Error(err)
	}

	c := NewClient(
		"api-key",
		BaseURL(u),
	)

	assert.Equal(t, c.BaseURL, u, "base url should match passed in base url")
}

func TestHttpClient(t *testing.T) {
	testHttpClient := &http.Client{
		Timeout: 1 * time.Minute,
	}

	c := NewClient(
		"api-key",
		HTTPClient(testHttpClient),
	)

	assert.Equal(t, c.HTTPClient, testHttpClient, "http client should match passed in http client")
}
