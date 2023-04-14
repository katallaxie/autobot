package ports

import "github.com/katallaxie/autobot/internal/analyzer"

// Analyzer ...
type Analyzer interface {
	Analyze(locale, content string) (sentence analyzer.Sentence)
}
