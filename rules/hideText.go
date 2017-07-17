package rules

import (
	"github.com/golang-commonmark/markdown"
)

type OpenHideText struct {
	Content string
	lvl     int
}

func (j *OpenHideText) Tag() string {
	return "hide"
}

func (j *OpenHideText) Opening() bool {
	return true
}

func (j *OpenHideText) Closing() bool {
	return true
}

func (j *OpenHideText) Block() bool {
	return false
}

func (j *OpenHideText) Level() int {
	return j.lvl
}

func (j *OpenHideText) SetLevel(lvl int) {
	j.lvl = lvl
}

func RuleHideText(s *markdown.StateInline, silent bool) (_ bool) {
	src := s.Src
	max := s.PosMax
	start := s.Pos
	marker := "\\\\*"
	marker2 := "*\\\\"
	if start+3 > max {
		return
	}

	if string(src[start:start+3]) != marker {
		return
	}

	if silent {
		return
	}

	open := OpenHideText{}

	for s.Pos < max {
		m := s.Pos + 3
		if m > max {
			m = max
		}
		open.Content += string(src[s.Pos])
		if src[s.Pos:m] == marker2 {
			open.Content += "\\\\"
			s.Pos += 3
			break
		}
		s.Pos++
	}
	s.PosMax = s.Pos

	s.PushOpeningToken(&open)

	return true
}
