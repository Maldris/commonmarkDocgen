package rules

import (
	"github.com/golang-commonmark/markdown"
)

type TableHeader struct {
	lvl   int
	Lines bool
}

func (j *TableHeader) Tag() string {
	return "thead"
}

func (j *TableHeader) Opening() bool {
	return true
}

func (j *TableHeader) Closing() bool {
	return true
}

func (j *TableHeader) Block() bool {
	return true
}

func (j *TableHeader) Level() int {
	return j.lvl
}

func (j *TableHeader) SetLevel(lvl int) {
	j.lvl = lvl
}

func RuleTableSettings(s *markdown.StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	src := s.Src

	if len(src) < pos+7 {
		return
	}

	marker := src[pos : pos+6]

	if marker != "\\thead" {
		return
	}

	if silent {
		return true
	}

	var lines bool = false
	if src[pos+6] == 't' || src[pos+6] == 'l' {
		lines = true
	}

	s.Line = startLine + 1
	s.PushToken(&TableHeader{Lines: lines})

	return true
}
