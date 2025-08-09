package text_analysis

import (
	"context"

	"mmm-osint/internal/pkg/keyword"
	"mmm-osint/internal/pkg/pii"
)

type TextAnalysisService interface {
	AnalyzeText(ctx context.Context, text string) (*TextAnalysisResult, error)
	Close() error
}

type TextAnalysisResult struct {
	Text        string             `json:"text"`
	PIIResult   *pii.Result        `json:"pii_result"`
	Keywords    []keyword.Keyword  `json:"keywords"`
}

func (r *TextAnalysisResult) HasPII() bool {
	return r.PIIResult != nil && !r.PIIResult.IsEmpty()
}

func (r *TextAnalysisResult) HasKeywords() bool {
	return len(r.Keywords) > 0
}