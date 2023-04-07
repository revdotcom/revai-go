package revai

const (
	LanguageIdJobType        = "languageid"
	SentimentAnalysisJobType = "sentiment_analysis"
	TopicExtractionJobType   = "topic_extraction"
)

// Config for Source/Notification Configs
// Same usage as MediaURL and CallbackURL
// Allow exactly map{"Authorization": "YOUR_AUTH_HERE"} as of 4/7/23
type UrlConfig struct {
	Url         string            `json:"url,omitempty"`
	AuthHeaders map[string]string `json:"auth_headers,omitempty"`
}

// ListParams specifies the optional query parameters to most List methods.
type ListParams struct {
	Limit         int    `url:"limit,omitempty"`
	StartingAfter string `url:"starting_after,omitempty"`
}

// The following are generic json structures sent through callback URL
// Use for unmarshalling request body
// https://docs.rev.ai/api/JOB_TYPE/webhooks/
// JOB_TYPE options are: asynchronous, topic-extraction, sentiment-analysis, language-identification, custom-vocabulary
type GenericPostJson struct {
	Job                 *GenericJob `json:"job"`
	CustomVocabularyJob *GenericJob `json:"custom_vocabulary"`
}

type GenericJob struct {
	ID        string  `json:"id,omitempty"`
	Created   string  `json:"created_on,omitempty"`
	Completed string  `json:"completed_on,omitempty"`
	Metadata  string  `json:"metadata,omitempty"`
	Status    string  `json:"status,omitempty"`
	Duration  float64 `json:"duration_seconds,omitempty"`
	Type      string  `json:"type,omitempty"`
}
