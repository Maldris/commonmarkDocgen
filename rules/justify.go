package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang-commonmark/markdown"
)

type JustifyOpen struct {
	lvl  int
	Type int
}

type JustifyClose struct {
	lvl  int
	Type int
}

func (j *JustifyOpen) Tag() string {
	return "center"
}

func (j *JustifyOpen) Opening() bool {
	return true
}

func (j *JustifyOpen) Closing() bool {
	return false
}

func (j *JustifyOpen) Block() bool {
	return false
}

func (j *JustifyOpen) Level() int {
	return j.lvl
}

func (j *JustifyOpen) SetLevel(lvl int) {
	j.lvl = lvl
}

func (j *JustifyClose) Tag() string {
	return "center"
}

func (j *JustifyClose) Opening() bool {
	return false
}

func (j *JustifyClose) Closing() bool {
	return true
}

func (j *JustifyClose) Block() bool {
	return false
}

func (j *JustifyClose) Level() int {
	return j.lvl
}

func (j *JustifyClose) SetLevel(lvl int) {
	j.lvl = lvl
}

func RuleJustify(s *markdown.StateInline, silent bool) (_ bool) {
	src := s.Src
	max := s.PosMax
	start := s.Pos
	marker := src[start]

	if src[start] != ':' {
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

	s.PushOpeningToken(&JustifyOpen{Type: startCount - 1})

	s.Md.Inline.Tokenize(s)

	s.PushClosingToken(&JustifyClose{Type: startCount - 1})

	s.Pos = s.PosMax + startCount
	s.PosMax = max

	return true
}

func scanDelims(s *markdown.StateInline, start int) (canOpen bool, canClose bool, count int) {
	pos := start
	max := s.PosMax
	src := s.Src
	marker := src[start]

	lastChar, lastLen := utf8.DecodeLastRuneInString(src[:start])

	for pos < max && src[pos] == marker {
		pos++
	}
	count = pos - start

	nextChar, nextLen := utf8.DecodeRuneInString(src[pos:])

	isLastSpaceOrStart := lastLen == 0 || unicode.IsSpace(lastChar)
	isNextSpaceOrEnd := nextLen == 0 || unicode.IsSpace(nextChar)
	isLastPunct := !isLastSpaceOrStart && (isMarkdownPunct(lastChar) || unicode.IsPunct(lastChar))
	isNextPunct := !isNextSpaceOrEnd && (isMarkdownPunct(nextChar) || unicode.IsPunct(nextChar))

	leftFlanking := !isNextSpaceOrEnd && (!isNextPunct || isLastSpaceOrStart || isLastPunct)
	rightFlanking := !isLastSpaceOrStart && (!isLastPunct || isNextSpaceOrEnd || isNextPunct)

	if marker == '_' {
		canOpen = leftFlanking && (!rightFlanking || isLastPunct)
		canClose = rightFlanking && (!leftFlanking || isNextPunct)
	} else {
		canOpen = leftFlanking
		canClose = rightFlanking
	}

	return
}

func isMarkdownPunct(r rune) bool {
	if r > 0x7e {
		return false
	}
	return strings.IndexRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", r) > -1
}
