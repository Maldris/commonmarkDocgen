package docgen

import (
	"log"
	"os"
	"testing"
)

// TestMain is the root testing method
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGen tests document generation
func TestGen(t *testing.T) {
	err := OutputExamplePdf("test.pdf")
	if err != nil {
		t.Error(err)
	}
}

// OutputExamplePdf is a test method to generate an extremely simple pdf to test initial generation
func OutputExamplePdf(destination string) error {
	template := `{#% test marker %#}

I am a test document
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

\thead ln shalf

| a | b |
| --- | --- |
| test table | with cells |
| and more | cells |

\thead ll sfull

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

	doc.AddExtensionFunctions(map[string]interface{}{
		"logger": func() string {
			log.Print("test")
			return ""
		},
	})

	doc.RegisterSubTemplate("test-sub", `blah blah test subtemplate`)

	err := doc.Execute(inputMap)
	if err != nil {
		return err
	}

	err = doc.RenderToFile(destination)
	if err != nil {
		return err
	}
	_, err = doc.RenderToString()
	if err != nil {
		return err
	}
	return nil
}
