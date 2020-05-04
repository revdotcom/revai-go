package revai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMetadata = "test-metadata"
const testMediaURL = "https://support.rev.com/hc/en-us/article_attachments/200043975/FTC_Sample_1_-_Single.mp3"
const testJobID = "VTsmOAwsM46v"

func TestJobService_SubmitFile(t *testing.T) {
	f, err := os.Open("./testdata/img.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	params := &NewFileJobParams{
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

	params := &NewFileJobParams{
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

func TestJobService_Submit(t *testing.T) {
	params := &NewJobParams{
		MediaURL: testMediaURL,
	}

	ctx := context.Background()

	newJob, err := testClient.Job.Submit(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestJobService_SubmitWithOption(t *testing.T) {
	params := &NewJobParams{
		MediaURL: testMediaURL,
		Metadata: testMetadata,
	}

	ctx := context.Background()

	newJob, err := testClient.Job.Submit(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, testMetadata, newJob.Metadata, "meta data should be set")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestJobService_Get(t *testing.T) {
	params := &GetJobParams{
		ID: testJobID,
	}

	ctx := context.Background()

	newJob, err := testClient.Job.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
}

func TestJobService_Delete(t *testing.T) {
	params := &DeleteJobParams{
		ID: testJobID,
	}

	ctx := context.Background()

	if err := testClient.Job.Delete(ctx, params); err != nil {
		t.Error(err)
		return
	}
}

func TestJobService_List(t *testing.T) {
	params := &ListJobParams{}

	ctx := context.Background()

	jobs, err := testClient.Job.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, jobs, "jobs should not be nil")
}

func TestJobService_ListWithLimit(t *testing.T) {
	params := &ListJobParams{
		Limit: 2,
	}

	ctx := context.Background()

	jobs, err := testClient.Job.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 2, len(jobs), "it returns 2 jobs when limit is set to 2")
}
