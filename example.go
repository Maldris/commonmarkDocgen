package docgen

// OutputExamplePdf is a test method to generate an extremely simple pdf to test initial generation
func OutputExamplePdf(destination string) error {
	template := `I am a test document
hear me roar
I can have included templated content like '{{.include}}'
**and formatted content through markdown**
blah
blah
blah ~ ; test indent with a really reeeeeeeeeeeeeeee eeeeeeeeeeeeeeeeeeeaaaaaaaaaaa aaaaaaaaaaaaaaaaaaaa aaaaaaaaaaaaaaaaaa aaallllllllllllllllll lllllllllllllllllllll lllllyyyyyyyyyyyyyyyyyyy yyyyyyyyy long line
blah
blah
blah
blah

\hn

| a | b |
| --- | --- |
| test table | with cells |
| and more | cells |

\hl

| a | b |
| --- | --- |
| test table | with cells |
| and more | cells |

blah
blah
blah
blah
blah
blah ~ !test indent with a really reeeeeeeeeeeeeeee eeeeeeeeeeeeeeeeeeeaaaaaaaaaaa aaaaaaaaaaaaaaaaaaaa aaaaaaaaaaaaaaaaaa aaallllllllllllllllll lllllllllllllllllllll lllllyyyyyyyyyyyyyyyyyyy yyyyyyyyy long line
blah!
blah
blah ~ test! indent with a really reeeeeeeeeeeeeeee eeeeeeeeeeeeeeeeeeeaaaaaaaaaaa aaaaaaaaaaaaaaaaaaaa aaaaaaaaaaaaaaaaaa aaallllllllllllllllll lllllllllllllllllllll lllllyyyyyyyyyyyyyyyyyyy yyyyyyyyy long line

blah!`

	doc := NewDocument("test", template, nil)

	doc.SetDocumentHeader(`# I am a document header, great for titles, and document identifiers

---`)

	doc.SetPageHeader(`_header template_ all fancy and shit`)

	doc.SetPageFooter(`---
***footer template***, not quite the same, just a little different :::Page {{._page}}:::`)

	inputMap := map[string]interface{}{
		"include": "Im template included content",
	}
	err := doc.Execute(inputMap)
	if err != nil {
		return err
	}

	doc.RenderToFile(destination)
	return nil
}
