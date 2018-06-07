package docgen

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/Maldris/commonmarkDocgen/rules"
	"github.com/golang-commonmark/markdown"
)

const (
	nominalIndent    = 10
	bulletIndent     = 5
	nominalFontSize  = 12
	heading1FontSize = 24
	heading2FontSize = 22
	heading3FontSize = 20
	heading4FontSize = 18
	heading5FontSize = 16
	heading6FontSize = 14
	cellMargin       = 2 // marge top/bottom of cell
)

func (d *Document) render(tok markdown.Token) {
	if d.Debug {
		fmt.Printf("[docgen] [%v] %v :: %v # %v", tok.Tag(), reflect.TypeOf(tok), tok.Block(), tok)
	}
	switch tok.(type) {
	case *markdown.BlockquoteOpen:
		d.leftMargin += d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
	case *markdown.BlockquoteClose:
		d.leftMargin -= d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
		// d.fpdf.Write(d.lineHeight, "\n")
		d.fpdf.SetX(d.leftMargin)

	case *markdown.BulletListOpen:
		d.listStyle = dash
		d.leftMargin += d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
	case *markdown.BulletListClose:
		d.listStyle = none
		d.leftMargin -= d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.fpdf.SetX(d.leftMargin)

	case *markdown.OrderedListOpen:
		tk := tok.(*markdown.OrderedListOpen)
		d.listStart = tk.Map[0]
		d.listStyle = numberedDot
		d.leftMargin += d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
	case *markdown.OrderedListClose:
		d.listStyle = none
		d.leftMargin -= d.sizes.NominalIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.fpdf.SetX(d.leftMargin)

	case *markdown.ListItemOpen:
		switch d.listStyle {
		case dash:
			d.fpdf.Write(d.lineHeight, "-")
		case carat:
			d.fpdf.Write(d.lineHeight, ">")
		case numberedDot:
			tk := tok.(*markdown.ListItemOpen)
			d.fpdf.Writef(d.lineHeight, "%v.", (tk.Map[0]-d.listStart)+1)
		}
		d.leftMargin += d.sizes.BulletIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
	case *markdown.ListItemClose:
		d.leftMargin -= d.sizes.BulletIndent
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.fpdf.SetX(d.leftMargin)

	case *markdown.CodeBlock:
		tk := tok.(*markdown.CodeBlock)
		d.codeBlock(strings.Replace(tk.Content, "\n    ", "\n", -1))
	case *markdown.CodeInline: // TODO
		ci := tok.(*markdown.CodeInline)
		d.fpdf.Write(d.lineHeight, ci.Content)
	case *markdown.Fence:
		tk := tok.(*markdown.Fence)
		d.codeBlock(tk.Content)

	case *markdown.EmphasisOpen:
		d.applyStyle("I")
	case *markdown.EmphasisClose:
		d.removeStyle("I")
	case *markdown.StrongOpen:
		d.applyStyle("B")
	case *markdown.StrongClose:
		d.removeStyle("B")

	case *markdown.StrikethroughOpen: // TODO
		// d.writeMode = strikethrough
	case *markdown.StrikethroughClose: // TODO: make this work, the original attempt (commented below) didnt work
		// html := d.fpdf.HTMLBasicNew()
		// html.Write(d.lineHeight, "<del>"+d.buffer+"<del>")
		// d.writeMode = normal

	case *markdown.Softbreak:
		d.fpdf.Write(d.lineHeight, "\n")
	case *markdown.Hardbreak:
		d.fpdf.Write(d.lineHeight, "\n")

	case *markdown.HeadingOpen:
		hd := tok.(*markdown.HeadingOpen)
		switch hd.HLevel {
		case 1:
			d.fontSize = int(d.sizes.Heading1FontSize)
		case 2:
			d.fontSize = int(d.sizes.Heading2FontSize)
		case 3:
			d.fontSize = int(d.sizes.Heading3FontSize)
		case 4:
			d.fontSize = int(d.sizes.Heading4FontSize)
		case 5:
			d.fontSize = int(d.sizes.Heading5FontSize)
		case 6:
			d.fontSize = int(d.sizes.Heading6FontSize)
		}
		d.lineHeight *= float64(d.fontSize) / d.sizes.NominalFontSize
		d.flushTextStyling()
	case *markdown.HeadingClose:
		d.lineHeight /= float64(d.fontSize) / d.sizes.NominalFontSize
		d.fontSize = int(d.sizes.NominalFontSize)
		d.flushTextStyling()
		d.fpdf.Write(d.lineHeight, "\n\n")

	case *markdown.HTMLBlock:
		tk := tok.(*markdown.HTMLBlock)
		html := d.fpdf.HTMLBasicNew()
		html.Write(d.lineHeight, tk.Content)
		d.fpdf.SetX(d.leftMargin)
	case *markdown.HTMLInline:
		tk := tok.(*markdown.HTMLInline)
		html := d.fpdf.HTMLBasicNew()
		html.Write(d.lineHeight, tk.Content)

	case *markdown.Hr:
		wd, _ := d.fpdf.GetPageSize()
		d.fpdf.Line(d.leftMargin, d.fpdf.GetY(), wd-d.leftMargin, d.fpdf.GetY())
		d.fpdf.Write(d.lineHeight, "\n")

	case *markdown.Image: // TODO
	case *markdown.Inline:
		il := tok.(*markdown.Inline)
		for _, tok := range il.Children {
			d.render(tok)
		}

	case *markdown.LinkOpen:
		ln := tok.(*markdown.LinkOpen)
		d.writeMode = link
		d.link.ref = ln.Href
	case *markdown.LinkClose:
		d.applyStyle("U")
		if d.link.ref != "" {
			d.fpdf.WriteLinkString(d.lineHeight, d.link.text, d.link.ref)
		} else {
			d.fpdf.Write(d.lineHeight, d.link.text)
		}
		d.removeStyle("U")
		d.writeMode = normal

	case *markdown.ParagraphOpen:
	case *markdown.ParagraphClose:
		d.fpdf.Write(d.lineHeight, "\n\n")
		d.fpdf.SetX(d.leftMargin)

	case *markdown.TableOpen:
		d.table.rows = [][]cell{}
	case *markdown.TableClose:
		d.tableMulti()
		d.flushTextStyling()
		d.fpdf.SetLeftMargin(d.leftMargin)
		d.fpdf.SetX(d.leftMargin)
		d.fpdf.Write(d.lineHeight, "\n")

	case *markdown.TheadOpen:
	case *markdown.TheadClose:

	case *markdown.TrOpen:
		d.table.rows = append(d.table.rows, []cell{})
	case *markdown.TrClose:

	case *markdown.ThOpen:
		d.writeMode = tableHead
	case *markdown.ThClose:
		d.writeMode = normal

	case *markdown.TbodyOpen:
	case *markdown.TbodyClose:

	case *markdown.TdOpen:
		d.writeMode = tableCell
	case *markdown.TdClose:
		d.writeMode = normal

	case *markdown.Text:
		txt := tok.(*markdown.Text)
		content := strings.Replace(txt.Content, "~", "    ", -1) // REVIEW: temp till editor buttons for tab in ui
		content = strings.Replace(content, "\t", "    ", -1)
		switch d.writeMode {
		case normal:
			switch d.alignment {
			case alignLeft:
				d.fpdf.Write(d.lineHeight, content)
			case alignCenter:
				d.fpdf.WriteAligned(0, d.lineHeight, content, "C")
			case alignRight:
				d.fpdf.WriteAligned(0, d.lineHeight, content, "R")
			}
		case tableHead:
			// d.table.head = append(d.table.head, content)
			d.table.rows[len(d.table.rows)-1] = append(d.table.rows[len(d.table.rows)-1], cell{text: strings.Replace(content, "\\n", "\n", -1), head: true})
		case tableCell:
			d.table.rows[len(d.table.rows)-1] = append(d.table.rows[len(d.table.rows)-1], cell{text: strings.Replace(content, "\\n", "\n", -1), head: false})
		case link:
			d.link.text = content
		case strikethrough:
			d.buffer = content
		}

	case *rules.JustifyOpen:
		tk := tok.(*rules.JustifyOpen)
		switch tk.Type {
		case 0:
			d.alignment = alignLeft
		case 1:
			d.alignment = alignCenter
		case 2:
			d.alignment = alignRight
		}
	case *rules.JustifyClose:
		d.alignment = alignLeft

	case *rules.PageBreak:
		d.fpdf.AddPage()

	case *rules.OpenHangingIndent:
		if d.indents == nil {
			d.indents = []float64{}
		}
		d.indents = append(d.indents, d.leftMargin)
		d.leftMargin = d.fpdf.GetX()
		d.fpdf.SetLeftMargin(d.leftMargin)
	case *rules.ClosenHangingIndent:
		l := len(d.indents)
		if l > 0 {
			d.leftMargin = d.indents[l-1]
			d.fpdf.SetLeftMargin(d.leftMargin)
			d.indents = d.indents[:l-1]
		} else {
			d.fpdf.SetLeftMargin(10)
		}

	case *rules.TableHeader:
		tk := tok.(*rules.TableHeader)
		d.table.lines = tk.Lines
	case *rules.OpenHideText:
		tk := tok.(*rules.OpenHideText)
		d.fpdf.SetFontSize(0)
		d.fpdf.Write(d.lineHeight, tk.Content)
		d.fpdf.SetFontSize(float64(d.fontSize))

	}
}

func (d *Document) applyStyle(style string) {
	if !strings.Contains(d.fontStyle, style) {
		d.fontStyle += style
	}
	d.flushTextStyling()
}

func (d *Document) removeStyle(style string) {
	d.fontStyle = strings.Replace(d.fontStyle, style, "", -1)
	d.flushTextStyling()
}

func (d *Document) flushTextStyling() {
	d.fpdf.SetFont(d.fontFamily, d.fontStyle, float64(d.fontSize))
}

func (d *Document) codeBlock(content string) {
	wpage, _ := d.fpdf.GetPageSize()
	lmarge, _, rmarge, _ := d.fpdf.GetMargins()
	r, g, b := d.fpdf.GetFillColor()
	d.fpdf.SetFillColor(180, 180, 180)
	d.leftMargin += d.sizes.NominalIndent
	d.fpdf.SetLeftMargin(d.leftMargin)
	d.fpdf.MultiCell(wpage-(lmarge+rmarge)-20, d.lineHeight, content, "", "", true)
	d.leftMargin -= d.sizes.NominalIndent
	d.fpdf.SetLeftMargin(d.leftMargin)
	d.fpdf.SetFillColor(r, g, b)
	d.fpdf.Write(d.lineHeight, "\n")
}

func (d *Document) tableMulti() {
	line := d.fpdf.GetLineWidth()

	d.fpdf.SetLineWidth(line / 4)
	_, pageh := d.fpdf.GetPageSize()
	_, _, _, mbottom := d.fpdf.GetMargins()

	cols := d.calcTableColumnWidths()

	for _, row := range d.table.rows {
		curx, y := d.fpdf.GetXY()
		x := curx

		height := 0.0

		for i, cell := range row {
			if cell.head {
				d.applyStyle("B")
			}
			lines := d.fpdf.SplitLines([]byte(cell.text), cols[i])
			h := float64(len(lines))*d.lineHeight + float64(int(d.sizes.CellMargin)*len(lines))
			if h > height {
				height = h
			}
			if cell.head {
				d.removeStyle("B")
			}
		}
		// add a new page if the height of the row doesn't fit on the page
		if d.fpdf.GetY()+height > pageh-mbottom {
			d.fpdf.AddPage()
			y = d.fpdf.GetY()
		}
		for i, cell := range row {
			if cell.head {
				d.applyStyle("B")
			}
			width := cols[i]
			if d.table.lines {
				d.fpdf.Rect(x, y, width, height, "")
			}
			d.fpdf.MultiCell(width, d.lineHeight+d.sizes.CellMargin, cell.text, "", "", false)
			x += width
			d.fpdf.SetXY(x, y)
			if cell.head {
				d.removeStyle("B")
			}
		}
		d.fpdf.SetXY(curx, y+height)
	}

	d.fpdf.SetLineWidth(line)
}

func (d *Document) calcTableColumnWidths() []float64 {

	wdspace := math.Ceil(d.fpdf.GetStringWidth(" "))

	type column struct {
		max     float64
		total   float64
		count   float64
		average float64
	}
	var c []column
	for _, rowVal := range d.table.rows {
		for col, colVal := range rowVal {
			if len(c) < col+1 {
				c = append(c, column{})
			}

			if colVal.head {
				d.removeStyle("B")
				d.flushTextStyling()
			}

			wd := d.fpdf.GetStringWidth(colVal.text) + (wdspace * float64(strings.Count(colVal.text, " ")+1+4))
			if colVal.head {
				d.removeStyle("B")
				d.flushTextStyling()
			}
			if wd > c[col].max {
				c[col].max = wd
			}
			c[col].total += wd
			c[col].count++
			c[col].average = c[col].total / c[col].count
		}
	}

	total := float64(0)
	avgTotal := float64(0)
	for _, col := range c {
		total += col.max
		avgTotal += col.average
	}

	wpage, _ := d.fpdf.GetPageSize()
	lmarge, _, rmarge, _ := d.fpdf.GetMargins()

	var cols []float64
	if total <= (wpage - (lmarge + rmarge)) {
		for _, col := range c {
			cols = append(cols, col.max)
		}
	} else {
		for _, col := range c {
			cols = append(cols, (col.average/avgTotal)*(wpage-(lmarge+rmarge)))
		}
	}
	return cols

}
