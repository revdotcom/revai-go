package revai

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentimentAnalysisService_SubmitPlainText(t *testing.T) {
	body, err := ioutil.ReadFile("./testdata/testtext.txt")
	if err != nil {
		t.Error(err)
		return
	}

	params := &SentimentAnalysisPlainTextParams{
		Text: string(body),
	}

	ctx := context.Background()

	newJob, err := testClient.SentimentAnalysis.SubmitPlainText(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestSentimentAnalysisService_SubmitJson(t *testing.T) {
	body, err := ioutil.ReadFile("./testdata/test.json")
	if err != nil {
		t.Error(err)
		return
	}

	var transcript Transcript
	err = json.Unmarshal(body, &transcript)
	if err != nil {
		t.Error(err)
		return
	}

	params := &SentimentAnalysisJsonParams{
		Metadata:   testMetadata,
		Transcript: transcript,
	}

	ctx := context.Background()

	newJob, err := testClient.SentimentAnalysis.SubmitTranscriptJson(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}
	if newJob != nil {
		fmt.Printf("%+v", *newJob)
	}
	assert.NotNil(t, newJob.ID, "new job id should not be nil")
	assert.Equal(t, testMetadata, newJob.Metadata, "meta data should be set")
	assert.Equal(t, "in_progress", newJob.Status, "response status should be in_progress")
}

func TestSentimentAnalysisService_Get(t *testing.T) {
	params := &GetSentimentAnalysisParams{
		ID: testSentimentAnalysis.ID,
	}

	ctx := context.Background()

	newJob, err := testClient.SentimentAnalysis.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob, "new job should not be nil")
}

func TestSentimentAnalysisService_GetJobById(t *testing.T) {
	params := &GetSentimentAnalysisParams{
		ID: testSentimentAnalysis.ID,
	}

	ctx := context.Background()

	newJob, err := testClient.SentimentAnalysis.GetJobById(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, newJob.ID, "new job id should not be nil")
}

func TestSentimentAnalysisService_Delete(t *testing.T) {
	deletableJob := makeTestSentimentAnalysis()

	params := &DeleteParams{
		ID: deletableJob.ID,
	}
	fmt.Printf("%+v", *deletableJob)

	ctx := context.Background()

	if job, err := testClient.SentimentAnalysis.Delete(ctx, params); err != nil {
		t.Error(err)
		return
	} else if job != nil {
		t.Error("Bad Status " + job.Status)
		return
	}
}

func TestSentimentAnalysisService_List(t *testing.T) {
	params := &ListParams{}

	ctx := context.Background()

	jobs, err := testClient.SentimentAnalysis.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, jobs, "jobs should not be nil")
}

func TestSentimentAnalysisService_ListWithLimit(t *testing.T) {
	params := &ListParams{
		Limit: 2,
	}

	ctx := context.Background()

	jobs, err := testClient.SentimentAnalysis.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, 2, len(jobs), "it returns 2 jobs when limit is set to 2")
}
