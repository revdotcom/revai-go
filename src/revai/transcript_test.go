package revai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranscriptService_Get(t *testing.T) {
	params := &GetTranscriptParams{
		ID: testJob.ID,
	}

	ctx := context.Background()

	transcript, err := testClient.Transcript.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, transcript.Monologues, "transcript monologue should not be nil")
}

func TestTranscriptService_GetText(t *testing.T) {
	params := &GetTranscriptParams{
		ID: testJob.ID,
	}

	ctx := context.Background()

	transcript, err := testClient.Transcript.GetText(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, transcript.Value, "transcript value should not be nil")
}
