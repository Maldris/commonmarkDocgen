package docgen

// OutputExamplePdf is a test method to generate an extremely simple pdf to test initial generation
func OutputExamplePdf(destination string) {
	template := `I am a test document
  hear me roar
  I can have included templated content like '{{.include}}'
  **and formatted content through markdown**
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah
  blah`

	doc := NewDocument(template, nil)

	doc.SetDocumentHeader(`# I am a document header, great for titles, and document identifiers

---`)

	doc.SetPageHeader(`_header template_ all fancy and shit`)

	doc.SetPageFooter(`---
***footer template***, not quite the same, just a little different :::Page {{._page}}:::`)

	inputMap := map[string]interface{}{
		"include": "Im template included content",
	}
	doc.Execute(inputMap)

	doc.RenderToFile(destination)
}
