package revai

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	testFileName = "./testdata/testaudio.mp3"
)

var (
	testClient *Client

	testJob               *Job
	testVocab             *CustomVocabulary
	testLanguageId        *LanguageId
	testTopicExtraction   *TopicExtraction
	testSentimentAnalysis *SentimentAnalysis
)

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

func setup() {
	testClient = NewClient(os.Getenv("REV_AI_API_KEY"))
	testJob = makeTestJob()
	testVocab = makeTestVocab()
	testLanguageId = makeTestLanguageId()
	testTopicExtraction = makeTestTopicExtraction()
	testSentimentAnalysis = makeTestSentimentAnalysis()
	fmt.Println("sleeping for 1 minute to allow file to be processed")
	time.Sleep(1 * time.Minute)
}

func makeTestJob() *Job {
	f := getTestFile()

	params := &NewFileJobParams{
		Media:    f, // some io.Reader
		Filename: f.Name(),
	}

	ctx := context.Background()

	job, err := testClient.Job.SubmitFile(ctx, params)
	if err != nil {
		panic(err)
	}

	return job
}

func makeTestVocab() *CustomVocabulary {
	params := &CreateCustomVocabularyParams{
		CustomVocabularies: []Phrase{
			{
				Phrases: []string{"hello"},
			},
		},
	}

	ctx := context.Background()

	vocab, err := testClient.CustomVocabulary.Create(ctx, params)
	if err != nil {
		panic(err)
	}

	return vocab
}

func makeTestLanguageId() *LanguageId {
	f := getTestFile()

	params := &LanguageIdFileParams{
		Media:    f, // some io.Reader
		Filename: f.Name(),
	}

	ctx := context.Background()

	job, err := testClient.LanguageId.SubmitFile(ctx, params)
	if err != nil {
		panic(err)
	}

	return job
}

func makeTestTopicExtraction() *TopicExtraction {
	data := getTestJsonData()

	var transcript Transcript
	err := json.Unmarshal(data, &transcript)
	if err != nil {
		panic(err)
	}

	params := &TopicExtractionJsonParams{
		Transcript: transcript,
	}

	ctx := context.Background()

	job, err := testClient.TopicExtraction.SubmitTranscriptJson(ctx, params)
	if err != nil {
		panic(err)
	}

	return job
}

func makeTestSentimentAnalysis() *SentimentAnalysis {
	data := getTestJsonData()

	var transcript Transcript
	err := json.Unmarshal(data, &transcript)
	if err != nil {
		panic(err)
	}

	params := &SentimentAnalysisJsonParams{
		Transcript: transcript,
	}

	ctx := context.Background()

	job, err := testClient.SentimentAnalysis.SubmitTranscriptJson(ctx, params)
	if err != nil {
		panic(err)
	}

	return job
}

func getTestFile() *os.File {
	f, err := os.Open(testFileName)
	if err != nil {
		panic(err)
	}

	return f
}

func getTestJsonData() []byte {
	body, err := ioutil.ReadFile("./testdata/test.json")
	if err != nil {
		panic(err)
	}

	return body
}
