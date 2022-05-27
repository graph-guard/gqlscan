package main

import (
	"bytes"
	_ "embed"
	"flag"
	"go/format"
	"io/fs"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed gqlscan.gotmpl
var tmplGqlscan string

//go:embed scan_body.gotmpl
var tmplScanBody string

//go:embed callback.gotmpl
var tmplCallback string

//go:embed skip_irrelevant.gotmpl
var tmplSkipIrrelevant string

func main() {
	var fOutPath string
	flag.StringVar(
		&fOutPath,
		"out",
		"./gqlscan.go",
		"output file path",
	)
	flag.Parse()

	fl, err := os.OpenFile(
		fOutPath,
		os.O_CREATE|os.O_TRUNC|os.O_SYNC|os.O_WRONLY,
		fs.FileMode(0644),
	)
	if err != nil {
		log.Fatalf("opening output file: %v", err)
	}
	defer fl.Close()

	t, err := template.New("gqlscan").
		Funcs(sprig.TxtFuncMap()).
		Parse(tmplGqlscan)
	if err != nil {
		log.Fatalf("parsing main template: %v", err)
	}
	for _, r := range []struct {
		Name, Source string
	}{
		{"skip_irrelevant", tmplSkipIrrelevant},
		{"scan_body", tmplScanBody},
		{"callback", tmplCallback},
	} {
		if _, err := t.New(r.Name).Parse(r.Source); err != nil {
			log.Fatalf("parsing template %q: %v", r.Name, err)
		}
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "gqlscan", nil); err != nil {
		log.Fatalf("executing main template: %v", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("formatting: %v", err)
	}
	if _, err := fl.Write(p); err != nil {
		log.Fatalf("writing output: %v", err)
	}
}
