package revai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://api.rev.ai"
const defaultUserAgent = "revai-go-client"

type service struct {
	client *Client
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
	UserAgent  string

	APIKey string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Rev.ai API.
	Job *JobService
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

	c.common.client = c

	c.Job = (*JobService)(&c.common)

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

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	return req, nil
}

func (c *Client) newMultiPartRequest(mw *multipart.Writer, path string, body io.Reader) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		return err
	}
	defer resp.Body.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		return err
	}
	debug := b.String()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return &ErrBadStatusCode{
			OriginalBody: debug,
			Code:         resp.StatusCode,
		}
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(&b).Decode(v); err != nil {
		return fmt.Errorf("failed decoding response %w", err)
	}

	return nil
}

func makeReaderPart(mw *multipart.Writer, partName, filename string, partValue io.Reader) error {
	part, err := mw.CreateFormFile(partName, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, partValue); err != nil {
		return err
	}
	return nil
}

func makeStringPart(mw *multipart.Writer, partName string, partValue string) error {
	part, err := mw.CreateFormField(partName)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, strings.NewReader(partValue)); err != nil {
		return err
	}
	return nil
}
