package revai

import (
	"os"
	"testing"
)

var (
	testClient *Client
)

func TestMain(m *testing.M) {
	testClient = NewClient(os.Getenv("REV_AI_API_KEY"))

	os.Exit(m.Run())
}
