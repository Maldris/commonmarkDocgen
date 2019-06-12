package rules

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

type TableHeader struct {
	lvl       int
	Lines     bool
	Size      SizeMode
	Cols      []float64
	Colsum    float64
	Alignment AlignMode
}

type SizeMode int

const (
	SizeUndefined SizeMode = iota // 0
	SizeWrap                      // 1
	SizeHalf                      // 2
	SizeFull                      // 3
)

type AlignMode int

const (
	AlignCenter AlignMode = iota // 0
	AlignLeft                    // 1
	AlignRight                   // 2
)

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

	tok := TableHeader{
		Lines:     true,
		Size:      SizeWrap,
		Cols:      []float64{},
		Alignment: AlignCenter,
	}

	exprs := strings.Split(src[pos+6:s.BMarks[startLine+1]], " ")

	for _, exp := range exprs {
		exp = strings.Trim(exp, " \r\n\t")
		switch {
		case strings.HasPrefix(exp, "a"):
			switch exp {
			case "al":
				tok.Alignment = AlignLeft
			case "ar":
				tok.Alignment = AlignRight
			default:
				tok.Alignment = AlignCenter
			}
		case strings.HasPrefix(exp, "l"):
			tok.Lines = exp == "lt" || exp == "ll"
		case strings.HasPrefix(exp, "s"):
			switch exp {
			case "shalf":
				tok.Size = SizeHalf
			case "sfull":
				tok.Size = SizeFull
			case "swrap":
				fallthrough
			default:
				tok.Size = SizeWrap
			}
		case strings.HasPrefix(exp, "c"):
			exp = exp[1:]
			ts := strings.Split(exp, ":")
			sum := float64(0)
			for _, t := range ts {
				if t == "" {
					continue
				}
				f, err := strconv.ParseFloat(t, 64)
				if err != nil {
					fmt.Printf("[docgen] Error parsing column size: %v", err)
					return
				}
				tok.Cols = append(tok.Cols, f)
				sum += f
			}
			tok.Colsum = sum
		}
	}

	s.Line = startLine + 1
	s.PushToken(&tok)

	return true
}
