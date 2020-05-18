package revai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVocabID = "cvY6kfhmV4srTd"

func TestCustomVocabularyService_Get(t *testing.T) {
	params := &GetCustomVocabularyParams{
		ID: testVocabID,
	}

	ctx := context.Background()

	vocabulary, err := testClient.CustomVocabulary.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, vocabulary.ID, "vocabulary id should not be nil")
}
