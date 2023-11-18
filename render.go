package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
)

func renderHTML() {

	tmpl, err := template.ParseFiles("template/template.html", "template/header.html")

	if err != nil {
		panic("Parse failed!!")
	}

	dat, _ := os.ReadFile("source/sample.md")
	var buf bytes.Buffer
	if err := goldmark.Convert(dat, &buf); err != nil {
		panic(err)
	}

	var htmlOutputBuffer bytes.Buffer

	tmpl.Execute(&htmlOutputBuffer, template.HTML(buf.String()))

	replacedHTML, _ := replaceCodeParts(htmlOutputBuffer.Bytes(), styles.Monokai)

	fmt.Println(replacedHTML)

	os.WriteFile("dist/index.html", []byte(replacedHTML), 0644)

}
