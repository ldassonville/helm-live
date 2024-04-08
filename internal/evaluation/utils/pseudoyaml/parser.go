package pseudoyaml

import (
	"strings"
	"unicode"
)

type Section struct {
	StartAt int `json:"start_at"`
	EndAt   int `json:"end_at"`

	Indent  string `json:"indent"`
	Content string `json:"content"`

	ProcessingError bool `json:"processing_error"`
	linesBuffer     []string
}

func NewSection(startAt int, line string) *Section {
	return &Section{
		StartAt:     startAt,
		linesBuffer: []string{line},
	}
}

func (s *Section) AddLine(line string) {
	s.linesBuffer = append(s.linesBuffer, line)
}

func (s *Section) Close() {
	if len(s.linesBuffer) > 0 {
		s.EndAt = s.StartAt + len(s.linesBuffer) - 1
	} else {
		s.EndAt = s.StartAt
	}

	s.Content = s.rawContent()
}

func (s *Section) rawContent() string {
	if s.linesBuffer == nil {
		return ""
	}
	return strings.Join(s.linesBuffer[:], "\n")
}

func LookupFirstMatch(pathParts []string, value string) string {

	return extractByPath(pathParts, value, "")
}

func extractByPath(pathParts []string, value string, indent string) string {

	if len(pathParts) == 0 {
		return ""
	}
	sections := extractSection(value, pathParts[0], indent)

	if len(pathParts) == 1 {
		for _, section := range sections {
			for _, line := range section.linesBuffer {
				if strings.Contains(line, ":") {
					value := strings.Split(section.linesBuffer[0], ":")[1]
					return strings.TrimSpace(value)
				}
			}
		}
		return ""
	}

	for _, section := range sections {
		res := extractByPath(pathParts[1:], section.rawContent(), section.Indent)
		if len(res) > 0 {
			return res
		}
	}
	return ""
}

func extractSection(value string, prefix string, indent string) (sections []*Section) {

	lines := strings.Split(value, "\n")

	stillConcerned := false
	yamlIndentation := 0

	var section *Section

	for lineNumber, line := range lines {

		if !stillConcerned && !isComment(line) && strings.HasPrefix(line, indent+prefix) {
			stillConcerned = true
			yamlIndentation = countSpacePrefix(line)
			section = NewSection(lineNumber+1, line)
			sections = append(sections, section)
			continue
		}

		if stillConcerned {
			currentIndent := countSpacePrefix(line)

			if currentIndent <= yamlIndentation {
				stillConcerned = false
				section.Close()
				section = nil
				continue
			} else {
				section.Indent = strings.Repeat(" ", currentIndent)
			}
			section.AddLine(line)
		}
	}

	if section != nil {
		section.Close()
	}

	return
}

func countSpacePrefix(line string) int {
	total := 0
	start := 0
	for ; start < len(line); start++ {
		if !unicode.IsSpace(rune(line[start])) {
			return total
		}
		total++
	}
	return total
}

// isComment give it the current yaml line is a comment
func isComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}
