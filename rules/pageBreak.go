package rules

import "github.com/golang-commonmark/markdown"

type PageBreak struct {
	lvl int
}

func (j *PageBreak) Tag() string {
	return "br"
}

func (j *PageBreak) Opening() bool {
	return true
}

func (j *PageBreak) Closing() bool {
	return true
}

func (j *PageBreak) Block() bool {
	return true
}

func (j *PageBreak) Level() int {
	return j.lvl
}

func (j *PageBreak) SetLevel(lvl int) {
	j.lvl = lvl
}

func RulePageBreak(s *markdown.StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	src := s.Src

	marker := src[pos : pos+5]

	if marker != "\\page" {
		return
	}

	if silent {
		return true
	}

	s.Line = startLine + 1
	s.PushToken(&PageBreak{})

	return true
}
