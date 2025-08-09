package text_analysis

import (
	"context"

	"mmm-osint/internal/pkg/keyword"
	"mmm-osint/internal/pkg/pii"
)

type textAnalysisService struct {
	piiExtractor     pii.Extractor
	keywordExtractor keyword.Extractor
}

func NewTextAnalysisService(piiExtractor pii.Extractor, keywordExtractor keyword.Extractor) TextAnalysisService {
	return &textAnalysisService{
		piiExtractor:     piiExtractor,
		keywordExtractor: keywordExtractor,
	}
}

func (s *textAnalysisService) AnalyzeText(ctx context.Context, text string) (*TextAnalysisResult, error) {
	result := &TextAnalysisResult{
		Text: text,
	}

	piiResult, err := s.piiExtractor.ExtractPII(ctx, text)
	if err != nil {
		return nil, err
	}
	result.PIIResult = piiResult

	keywords, err := s.keywordExtractor.ExtractKeywordsWithScores(ctx, text, keyword.DefaultOptions())
	if err != nil {
		return nil, err
	}
	result.Keywords = keywords

	return result, nil
}

func (s *textAnalysisService) Close() error {
	if err := s.piiExtractor.Close(); err != nil {
		return err
	}
	if err := s.keywordExtractor.Close(); err != nil {
		return err
	}
	return nil
}