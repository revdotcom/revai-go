package revai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptionService_Get(t *testing.T) {
	params := &GetCaptionParams{
		JobID: testJob.ID,
	}

	ctx := context.Background()

	caption, err := testClient.Caption.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, caption.Value, "caption value should not be nil")
}
