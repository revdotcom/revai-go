package revai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMetadata = "test-metadata"

func TestJobService_SubmitFile(t *testing.T) {
	f, err := os.Open("./testdata/img.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	params := &NewJobParams{
		Media:    f,
		Filename: f.Name(),
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

func TestJobService_SubmitFileWithOption(t *testing.T) {
	f, err := os.Open("./testdata/img.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	params := &NewJobParams{
		Media:    f,
		Filename: f.Name(),
		JobOptions: &JobOptions{
			Metadata: testMetadata,
		},
	}

	ctx := context.Background()

	newJob, err := testClient.Job.SubmitFile(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, testMetadata, newJob.Metadata, "meta data should be set")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}
