package revai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVocabID = "cvY6kfhmV4srTd"

func TestCustomVocabularyService_Create(t *testing.T) {
	params := &CreateCustomVocabularyParams{
		CustomVocabularies: []Phrase{
			{
				Phrases: []string{"paul lefalux"},
			},
		},
	}

	ctx := context.Background()

	vocabulary, err := testClient.CustomVocabulary.Create(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, vocabulary.ID, "vocabulary id should not be nil")
}

func TestCustomVocabularyService_Get(t *testing.T) {
	params := &GetCustomVocabularyParams{
		ID: testVocab.ID,
	}

	ctx := context.Background()

	vocabulary, err := testClient.CustomVocabulary.Get(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.NotNil(t, vocabulary.ID, "vocabulary id should not be nil")
}

func TestCustomVocabularyService_List(t *testing.T) {
	params := &ListParams{}

	ctx := context.Background()

	vocabularies, err := testClient.CustomVocabulary.List(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Greater(t, len(vocabularies), 0, "vocabularies should not be nil")
}

func TestCustomVocabularyService_Delete(t *testing.T) {
	deletableVocab := makeTestVocab()

	params := &DeleteParams{
		ID: deletableVocab.ID,
	}

	ctx := context.Background()

	err := testClient.CustomVocabulary.Delete(ctx, params)
	if err != nil {
		t.Error(err)
		return
	}
}
