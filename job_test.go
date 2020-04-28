package revai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobService_SubmitFile(t *testing.T) {
	f, err := os.Open("./testdata/testaudio.mp3")
	if err != nil {
		t.Error(err)
		return
	}

	params := &NewJobParams{
		Media: f,
	}

	ctx := context.Background()

	newJob, err := testClient.Job.SubmitFile(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}
