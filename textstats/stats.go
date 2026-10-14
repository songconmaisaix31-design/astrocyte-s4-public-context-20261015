package textstats

import (
	"strings"
	"unicode/utf8"
)

type Stats struct {
	Words      int
	Characters int
	Lines      int
}

func Count(text string) Stats {
	lines := strings.Count(text, "\n")
	if text != "" && !strings.HasSuffix(text, "\n") {
		lines++
	}
	return Stats{Words: len(strings.Fields(text)), Characters: utf8.RuneCountInString(text), Lines: lines}
}

func CountAll(texts []string) Stats {
	var result Stats
	for _, text := range texts {
		stats := Count(text)
		result.Words += stats.Words
		result.Characters += stats.Characters
		result.Lines += stats.Lines
	}
	return result
}
