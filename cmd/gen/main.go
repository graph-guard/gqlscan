package main

import (
	"bytes"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates
var tmpls embed.FS

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

	t := template.New("").Funcs(sprig.TxtFuncMap())
	if err := fs.WalkDir(
		tmpls,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			fileName := filepath.Base(path)
			if !strings.HasSuffix(fileName, ".gotmpl") {
				return nil
			}
			name := fileName[:len(fileName)-len(".gotmpl")]

			c, err := fs.ReadFile(tmpls, path)
			if err != nil {
				return fmt.Errorf("reading template (%s): %w", path, err)
			}

			if name != "gqlscan" {
				c = append([]byte(fmt.Sprintf("\n/*<%s>*/\n", name)), c...)
				if c[len(c)-1] != '\n' {
					c = append(c, '\n')
				}
				c = append(c, []byte(fmt.Sprintf("/*</%s>*/\n", name))...)
			}

			if _, err := t.New(name).
				Funcs(sprig.TxtFuncMap()).
				Parse(string(c)); err != nil {
				return fmt.Errorf("parsing template (%s): %w", path, err)
			}
			return nil
		},
	); err != nil {
		log.Fatalf("walking templates: %v", err)
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "gqlscan", nil); err != nil {
		log.Fatalf("executing main template: %v", err)
	}

	p, err := format.Source(buf.Bytes())
	if err != nil {
		if _, err := fl.Write(buf.Bytes()); err != nil {
			log.Fatalf("writing unformatted output: %v", err)
		}
		log.Fatalf("formatting: %v", err)
	}
	if _, err := fl.Write(p); err != nil {
		log.Fatalf("writing formatted output: %v", err)
	}
}
