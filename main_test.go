package revai

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	testFileName = "./testdata/testaudio.mp3"
)

var (
	testClient *Client

	testJob   *Job
	testVocab *CustomVocabulary
)

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

func setup() {
	testClient = NewClient(os.Getenv("REV_AI_API_KEY"))
	testJob = makeTestJob()
	testVocab = makeTestVocab()
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

func getTestFile() *os.File {
	f, err := os.Open(testFileName)
	if err != nil {
		panic(err)
	}

	return f
}
