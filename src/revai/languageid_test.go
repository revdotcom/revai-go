package revai

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageIdService_SubmitFile(t *testing.T) {
	f, err := os.Open("./testdata/testaudio.mp3")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	params := &LanguageIdFileParams{
		Media:    f,
		Filename: f.Name(),
	}

	ctx := context.Background()

	newJob, err := testClient.LanguageId.SubmitFile(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestLanguageIdService_SubmitFileWithOption(t *testing.T) {
	f, err := os.Open("./testdata/testaudio.mp3")
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()

	params := &LanguageIdFileParams{
		Media:    f,
		Filename: f.Name(),
		Options: &LanguageIdFileOptions{
			Metadata: testMetadata,
		},
	}

	ctx := context.Background()

	newJob, err := testClient.LanguageId.SubmitFile(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, testMetadata, newJob.Metadata, "meta data should be set")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestLanguageIdService_SubmitURL(t *testing.T) {
	params := &LanguageIdParams{
		SourceConfig: &UrlConfig{
			Url: testMediaURL,
		},
	}

	ctx := context.Background()

	newJob, err := testClient.LanguageId.SubmitURL(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestLanguageIdService_Get(t *testing.T) {
	params := &GetLanguageIdParams{
		ID: testLanguageId.ID,
	}

	ctx := context.Background()

	newJob, err := testClient.LanguageId.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
}

func TestLanguageIdService_GetById(t *testing.T) {
	params := &GetLanguageIdParams{
		ID: testLanguageId.ID,
	}

	ctx := context.Background()

	newJob, err := testClient.LanguageId.GetJobById(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
}

func TestLanguageIdService_Delete(t *testing.T) {
	deletableJob := makeTestLanguageId()

	params := &DeleteParams{
		ID: deletableJob.ID,
	}

	ctx := context.Background()

	if job, err := testClient.LanguageId.Delete(ctx, params); err != nil {
		t.Error(err)
		return
	} else if job != nil {
		t.Error("Bad Status " + job.Status)
		return
	}
}

func TestLanguageIdService_List(t *testing.T) {
	params := &ListParams{}

	ctx := context.Background()

	jobs, err := testClient.LanguageId.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, jobs, "jobs should not be nil")
}

func TestLanguageIdService_ListWithLimit(t *testing.T) {
	params := &ListParams{
		Limit: 2,
	}

	ctx := context.Background()

	jobs, err := testClient.LanguageId.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 2, len(jobs), "it returns 2 jobs when limit is set to 2")
}
