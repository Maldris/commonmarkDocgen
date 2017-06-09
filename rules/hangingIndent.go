package rules

import (
	"github.com/golang-commonmark/markdown"
)

type OpenHangingIndent struct {
	lvl int
}

type ClosenHangingIndent struct {
	lvl int
}

func (j *OpenHangingIndent) Tag() string {
	return "hang"
}

func (j *OpenHangingIndent) Opening() bool {
	return true
}

func (j *OpenHangingIndent) Closing() bool {
	return false
}

func (j *OpenHangingIndent) Block() bool {
	return false
}

func (j *OpenHangingIndent) Level() int {
	return j.lvl
}

func (j *OpenHangingIndent) SetLevel(lvl int) {
	j.lvl = lvl
}

func (j *ClosenHangingIndent) Tag() string {
	return "hang"
}

func (j *ClosenHangingIndent) Opening() bool {
	return false
}

func (j *ClosenHangingIndent) Closing() bool {
	return true
}

func (j *ClosenHangingIndent) Block() bool {
	return false
}

func (j *ClosenHangingIndent) Level() int {
	return j.lvl
}

func (j *ClosenHangingIndent) SetLevel(lvl int) {
	j.lvl = lvl
}

func RuleHangIndent(s *markdown.StateInline, silent bool) (_ bool) {
	src := s.Src
	max := s.PosMax
	start := s.Pos
	marker := src[start]

	if src[start] != '!' {
		return
	}

	if silent {
		return
	}

	canOpen, _, startCount := scanDelims(s, start)
	s.Pos += startCount
	if !canOpen {
		s.Pending.WriteString(src[start:s.Pos])
		return true
	}

	stack := []int{startCount}
	found := false

	for s.Pos < max {
		if src[s.Pos] == marker {
			canOpen, canClose, count := scanDelims(s, s.Pos)

			if canClose {
				oldCount := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				newCount := count

				for oldCount != newCount {
					if newCount < oldCount {
						stack = append(stack, oldCount-newCount)
						break
					}

					newCount -= oldCount

					if len(stack) == 0 {
						break
					}

					s.Pos += oldCount
					oldCount = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}

				if len(stack) == 0 {
					startCount = oldCount
					found = true
					break
				}

				s.Pos += count
				continue
			}

			if canOpen {
				stack = append(stack, count)
			}

			s.Pos += count
			continue
		}

		s.Md.Inline.SkipToken(s)
	}

	if !found {
		s.Pos = start
		return
	}

	s.PosMax = s.Pos
	s.Pos = start + startCount

	s.PushOpeningToken(&OpenHangingIndent{})

	s.Md.Inline.Tokenize(s)

	s.PushClosingToken(&ClosenHangingIndent{})

	s.Pos = s.PosMax + startCount
	s.PosMax = max

	return true
}
