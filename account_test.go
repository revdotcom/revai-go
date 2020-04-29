package revai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountService_Get(t *testing.T) {
	ctx := context.Background()

	newJob, err := testClient.Account.Get(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.Email, "account email should not be nil")
}
