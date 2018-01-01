package docgen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Maldris/commonmarkDocgen/ext/num2words"
	"github.com/Maldris/commonmarkDocgen/rules"
	"github.com/golang-commonmark/markdown"
	"github.com/jung-kurt/gofpdf"
)

// Document is the root docgen object, and represents an instance of document generation
type Document struct {
	t            *template.Template
	template     string
	name         string
	subTemplates struct {
		docHeader  string
		pageHeader string
		pageFooter string
		main       []struct {
			name string
			body string
		}
	}
	params map[string]interface{}
	parser *markdown.Markdown
	fpdf   *gofpdf.Fpdf

	fontFamily string
	fontSize   int
	fontStyle  string

	lineHeight float64
	leftMargin float64
	writeMode  writeMode
	alignment  alignment

	table struct {
		lines bool
		rows  [][]cell
	}
	link struct {
		ref  string
		text string
	}
	buffer     string
	listStyle  listStyle
	listStart  int
	indents    []float64
	extensions template.FuncMap
}

type alignment uint

const (
	alignLeft alignment = iota
	alignCenter
	alignRight
)

type cell struct {
	text string
	head bool
}

type writeMode uint

const (
	normal writeMode = iota
	tableHead
	tableCell
	link
	strikethrough
)

type listStyle uint

const (
	none listStyle = iota
	dash
	dot
	carat
	numberedDash
	numberedDot
	numberedArc
)

/*PdfConfig A utility config object used to provide the base page size details when initialising the pdf
 * Fields:
 * 	Portrait: Flag identifying if the page layout should be portrait, if false, it will be landscape
 * 	Metric:   Flag identifying if the units of measurement used should be in metric (thus mm), if false, inches are used
 * 	Paper:    String representing the size of paper to use, options are: "A3", "A4", "A5", "Letter", or "Legal"
 */
type PdfConfig struct {
	Portrait bool
	Metric   bool
	Paper    string
}

/*NewDocument Creates a new document object that represents an instance of document generation
 * Params:
 * 	template (string): the template that is the body of the document
 * 	conf (*PdfConfig): config object to determin the pdf page size, if nil, the default of a portrait A4 page measured in mm will be used
 */
func NewDocument(name, templateStr string, conf *PdfConfig) *Document {
	if conf == nil {
		conf = &PdfConfig{
			Portrait: true,
			Metric:   true,
			Paper:    "A4",
		}
	}

	t := &template.Template{}

	doc := &Document{
		t:        t,
		name:     name,
		template: templateStr,
		parser:   markdown.New(),
		// fpdf:       pdf,
		fontSize:   nominalFontSize,
		fontStyle:  "",
		fontFamily: "Arial",
		lineHeight: 5,
		// leftMargin: leftMargin,
	}
	doc.pdfInit(conf)
	// pdf generator doesnt support typographic fancy quotes, overwriting the fancy ones to normal ones
	doc.table.lines = true
	doc.parser.Typographer = false
	doc.parser.Linkify = false
	doc.parser.Quotes = [4]rune{'"', '"', '\'', '\''}
	doc.parser.HTML = true
	markdown.RegisterBlockRule(1050, rules.RulePageBreak, nil)
	markdown.RegisterBlockRule(1055, rules.RuleTableSettings, nil)
	markdown.RegisterInlineRule(2000, rules.RuleHangIndent)
	markdown.RegisterInlineRule(2200, rules.RuleJustify)
	markdown.RegisterInlineRule(200, rules.RuleHideText)
	return doc
}

// AddExtensionFunctions adds functions that will be available to the template when executed, it can also be used to overwrite the default functions if necessary
func (d *Document) AddExtensionFunctions(funcs map[string]interface{}) {
	if d.extensions == nil {
		d.extensions = template.FuncMap{}
	}
	for key, val := range funcs {
		d.extensions[key] = val
	}
}

func (d *Document) pdfInit(conf *PdfConfig) {
	orientation := "L"
	if conf.Portrait {
		orientation = "P"
	}
	units := "mm"
	if !conf.Metric {
		units = "in"
	}
	if conf.Paper != "A3" && conf.Paper != "A4" && conf.Paper != "A5" && conf.Paper != "Letter" && conf.Paper != "Legal" {
		conf.Paper = "A4"
	}
	pdf := gofpdf.New(orientation, units, conf.Paper, "")
	// pdf := gofpdf.New("P", "mm", "A4", "")
	d.fpdf = pdf
	leftMargin, _, _, _ := pdf.GetMargins()
	d.leftMargin = leftMargin
	pdf.AddPage()
	pdf.SetFont(d.fontFamily, d.fontStyle, float64(d.fontSize))
}

// SetDocumentHeader is used to provide the template to use to generate the document header, a template only used once at the start of the document, but for which the page header still occurs
func (d *Document) SetDocumentHeader(template string) {
	d.subTemplates.docHeader = template
}

// SetPageHeader is used to provide a template that will be used to build a header section for each page in the pdf, except the first page
func (d *Document) SetPageHeader(template string) {
	d.subTemplates.pageHeader = template
	if template == "" {
		d.fpdf.SetHeaderFunc(func() {})
		return
	}
	d.fpdf.SetHeaderFunc(func() {
		d.params["_page"] = d.fpdf.PageNo()
		err := d.renderTemplate("_pageHeader", d.subTemplates.pageHeader)
		if err != nil {
			d.fpdf.SetError(err)
			return
		}
		d.fpdf.Write(d.lineHeight, "\n\n")
	})
}

// SetPageFooter is used to provide a template that will be used to build a footer section for each page in the pdf
func (d *Document) SetPageFooter(template string) {
	d.subTemplates.pageFooter = template
	if template == "" {
		d.fpdf.SetFooterFunc(func() {})
		return
	}
	d.fpdf.SetFooterFunc(func() {
		d.params["_page"] = d.fpdf.PageNo()
		finalMarkdown, err := templateSubstitution(d.t, "_footer", d.subTemplates.pageFooter, d.params, d.extensions)
		if err != nil {
			d.fpdf.SetError(err)
			return
		}
		strWd := d.fpdf.GetStringWidth(finalMarkdown)
		pgWd, pgHt := d.fpdf.GetPageSize()
		footHeight := strWd / pgWd
		_, _, _, bMarg := d.fpdf.GetMargins()
		if d.fpdf.GetY() <= pgHt-(bMarg+2*d.lineHeight+footHeight) {
			d.fpdf.SetY(pgHt - (bMarg + 2*d.lineHeight + footHeight))
			d.fpdf.Write(d.lineHeight, "\n\n")
		}
		tokens := d.parser.Parse([]byte(finalMarkdown))
		for _, tok := range tokens {
			d.render(tok)
		}
	})
}

func (d *Document) generateDocumentHeader() error {
	if d.subTemplates.docHeader == "" {
		return nil
	}
	err := d.renderTemplate("_header", d.subTemplates.docHeader)
	if err != nil {
		return err
	}
	d.fpdf.Write(d.lineHeight, "\n\n")
	return nil
}

// RegisterSubTemplate will register a new template with name and body, this new template can then be invoked inside the documents template
func (d *Document) RegisterSubTemplate(name, body string) {
	temp := struct {
		name string
		body string
	}{
		name: name,
		body: body,
	}
	d.subTemplates.main = append(d.subTemplates.main, temp)
}

// Execute takes in the parameters to use to generate the document, and does the template parse, and document generation, effectively executing all templates loaded into the document
func (d *Document) Execute(data map[string]interface{}) error {
	d.params = data
	err := d.parseSubtemplates()
	if err != nil {
		return err
	}
	err = d.generateDocumentHeader()
	if err != nil {
		return err
	}
	err = d.renderTemplate(d.name, d.template)
	return err
}

func (d *Document) parseSubtemplates() error {
	funcs := loadFuncs(d.extensions)
	var err error
	for _, temp := range d.subTemplates.main {
		_, err = d.t.New(temp.name).Funcs(funcs).Parse(temp.body)
		if err != nil {
			return err
		}
	}
	return err
}

// RenderToFile (called after Execute) this method renders the resultant pdf to a file at fname, as a fully qualified path with filename
func (d *Document) RenderToFile(fname string) error {
	if !d.fpdf.Ok() {
		err := d.fpdf.Error()
		fmt.Printf("[docgen] Error rendering PDF: %v", err)
		return err
	}
	d.fpdf.OutputFileAndClose(fname)
	return nil
}

// RenderToString (called after Execute) this method renders the output pdf as a byte string that can be stored in a database or similar
func (d *Document) RenderToString() (string, error) {
	if !d.fpdf.Ok() {
		err := d.fpdf.Error()
		fmt.Printf("[docgen] Error rendering PDF: %v", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	d.fpdf.Output(buf)
	return buf.String(), nil
}

func (d *Document) renderTemplate(name, temp string) error {
	finalMarkdown, err := templateSubstitution(d.t, name, temp, d.params, d.extensions)
	if err != nil {
		return err
	}
	tokens := d.parser.Parse([]byte(finalMarkdown))
	for _, tok := range tokens {
		d.render(tok)
	}
	return nil
}

func templateSubstitution(t *template.Template, name, tmp string, data interface{}, exts template.FuncMap) (string, error) {
	temp, _ := templateSplit(tmp)

	funcMap := loadFuncs(exts)
	var err error
	t, err = t.New(name).Funcs(funcMap).Parse(temp)
	if err != nil {
		fmt.Printf("[docgen] Template parse error: %v", err)
		return temp, err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		fmt.Printf("[docgen] Template render error: %v", err)
		return temp, err
	}

	return buf.String(), nil
}

func templateSplit(template string) (ret string, markUps []string) {
	for i := 0; i < len(template); i++ {
		if ((i + 3) < len(template)) && (template[i] == '{') && (template[i+1] == '#') && (template[i+2] == '%') { //opening tag of {#%
			var markupContent string
			i += 3 //move beyond the opening tag

			for i < len(template) {
				if ((i + 2) < len(template)) && (template[i] == '%') && (template[i+1] == '#') && (template[i+2] == '}') {
					i += 2

					spl := strings.Split(markupContent, "|||")
					markUps = append(markUps, spl[0])
					if len(spl) > 1 {
						ret += string(spl[1])
					}
					break
				} else {
					markupContent += string(template[i])
				}
				i++
			}
		} else {
			ret += string(template[i])
		}
	}
	return
}

func loadFuncs(exts template.FuncMap) template.FuncMap {
	funcMap := template.FuncMap{
		"eq":       eq,
		"Cell":     newCell,
		"ToUpper":  strings.ToUpper,
		"Currency": currencyFormat,
		"Date":     formatDate,
		"IntDate":  integerDateFormat,
		"num2word": num2words.Convert,
		"add":      add,
		"subtract": subtract,
		"multiply": multiply,
		"divide":   divide,
		"cb":       codeBlock,
		"dict":     dictionary,
		"yn":       boolString,
	}

	for key, val := range exts {
		funcMap[key] = val
	}
	return funcMap
}
