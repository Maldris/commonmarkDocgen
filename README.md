# commonmarkDocgen
generate PDF documents from a commonmark markdown template

Built using [jung kurt's gofpdf](https://github.com/jung-kurt/gofpdf) and [golang-commonmark's markdown parser](https://github.com/golang-commonmark/markdown)


## Remaining TODOs
 - add a full test suite that doesn't cause licensing issues
 - add some examples
 - write some full usage documentation to go with the examples
 - add support for:
     - images
     - strikethrough
     - inline code blocks


## Basic Usage

1. create a new document object, passing in the root template, its name, and either page size configuration, or nil
```
doc := NewDocument("name", template, nil)
```

2. if you need one, set a document header template (will be generated once, at the top of the first page of the document)
```
doc.SetDocumentHeader(`# I am a document header, great for titles, and document identifiers

---`)
```

3. if needed, set page headers and/or footer templates (will generate at the start and end of each page)
```
doc.SetPageHeader(`_header template_ at the top of each page`)
doc.SetPageFooter(`---
***footer template***, not quite the same, just a little different :::Page {{._page}}:::`)
```

4. execute the template, passing in any arguments used in the template
```
err := doc.Execute(arguments)
if err != nil {
 return err
}
```

5. render your pdf to either a local file, or a string
```
err = doc.RenderToFile(destination)
if err != nil {
 return err
}
```
OR
```
output, err := doc.RenderToString()
if err != nil {
 return err
}
```


## Advanced Concepts
### AddExtensionFunctions(funcs map[string]interface{})
If you have functions you wish to use in your templates, register them with **AddExtensionFunctions**
Note: functions defined in this way that are also on the included function set, will override the default funcationality

### RegisterSubTemplate(name, body string)
If your template needs to reference a subtemplate stored separately to your core template, add it to the system with **RegisterSubTemplate**


## Included functions

**Cell**: creates a special data type that stores a pointer to an interface{}, useful when you need to be able to modify a value in a parent scope inside the template (not normally allowed), very useful when manually numbering bullets or section headers when their order/quantity may change (additional documentation comming)

**ToUpper**: same as normal, converts arguement string to upper case

**Currency**: currency formats a float ($xxxxx.yy)

**Date**: formats a date in (2 January 2006) format

**IntDate**: takes an int64 unix timestamp and does the same date formatting as above

**num2word**: express a number in words

**add**: add two numbers

**subtract**: subtract two numbers

**multiply**: multiply two floats

**divide**: divide two floats

**dict**: create a dictionary (map[string]interface{}) from passed arguments (format of repeating string value pairs), useful to pass a subset of arguments to a subtemplates

**yn**: humanised bools, outputs yes or no instead of true or false

Note: a custom **eq** method is also provided to support Cell, keep this in mind if overwriting

### Using Cell
Using the **Cell** included function will create a new cell object which can be saved in the template from the returned data
```
{{$i := Cell 1}}
```
The cell itself is a pointer to a struct containing an interface field storing its value
When you want to display its content in the template, use **cell.Get()**, and to update the value, call **cell.Set()**
```
# {{$i.Get()}}. Introductions
{{$i.Set (add $i.Get 1)}}
# {{$i.Get()}}. Next Heading
```

## Extra Markup

This library comes with a few extra markup literals added, at both the block and inline level.

### Block

**\page** will insert a page break in the document

**\thead** is a header element that can be used to customize the behavior of following tables with a series of arguments that follow it in any order i.e. `\thead lt sfull c:1:2 ar`
> **lines** turn on and off drawing the tables lines
>> any argument beginning with l is interpreted as a lines argument
>>
>> `lt` and `ll` turn lines on, all others turn lines off (absent line argument is interpreted as lt)
>>
>> i.e. `\thead lf` will turn lines off

> **size mode** should the table wrap its content, or span to a certain size
>> **shalf** table will span full width of page
>>
>> **sfull** table will span half the page
>>
>> **swrap** (default) table will wrap its content
>>
>> i.e. `\thead shalf` will make any following tables half the page width

> **column ratios** allows setting a ratio of the width of each column in a table
>> any argument starting with **c** is a column argument
>>
>> columns are then denoted as width ratios separated by colons (:)
>>
>> each column will then be printed as a fraction of the sum of all argument components
>>
>> i.e. `\thead c:1:2:1` will applied to a 3 column table will have a central column twice the width of the outer two.

> **table alignment** allows the table, if smaller than page width, to be aligned to the left, center or right of the page
>> any argument starting with **a** is an alignment argument, with the following character denoting mode
>>
>> **al** aligns the table with the left margin
>>
>> **ac** (default) centers the table on the page
>>
>> **ar** aligns the table with the right margin
>>
>> i.e. `\thead al` will left align the table

### Inline

**hanging indent** any text surrounded by exclamation marks (!) will set a hanging indent of the start of the inline block, preventing any subsequent lines from reflowing to the left margin and instead to the distance the inline element starts from the left edge of the page

**justify** flanking text in colons (:) will set the justification type with:
> 1 being left justified i.e. :text:
>
> 2 being center justified i.e. ::text::
>
> 3 being right justified i.e. :::text:::

**hidden text** allows text to be written such that it is not visible when reading the document, but accessable when parsing the file (useful with tools like docusign that allow attaching behaviours to textural anchors).

Any text enclosed by the tags `\\*` and `*\\` will be printed with a width and font size of 0
