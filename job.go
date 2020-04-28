package revai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

// client.Job.SubmitLocalFile()
// client.Job.SubmitURL()
// client.Job.Details()

type JobService service

type NewJob struct {
	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
}

type NewJobParams struct {
	Media io.Reader
}

func (s *JobService) SubmitFile(ctx context.Context, params *NewJobParams) (*NewJob, error) {
	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)

	if err := makeReaderPart(mw, "media", params.Media); err != nil {
		return nil, err
	}

	if err := mw.Close(); err != nil {
		return nil, fmt.Errorf("failed closing multipart-form writer %w", err)
	}

	req, err := s.client.newMultiPartRequest(mw, "/speechtotext/v1/jobs", body)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	var newJob NewJob
	if err := s.client.do(ctx, req, &newJob); err != nil {
		return nil, err
	}

	return &newJob, nil
}
