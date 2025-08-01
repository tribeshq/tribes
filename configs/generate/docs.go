package main

import (
	"os"
	"text/template"
)

func generateDocsFile(path string, env []Env) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	funcMap := template.FuncMap{
		"backtick": func(s string) string {
			return "`" + s + "`"
		},
		"quote": func(s string) string {
			return `"` + s + `"`
		},
	}
	tmpl := template.Must(template.New("docs").Funcs(funcMap).Parse(docsTemplate))

	err = tmpl.Execute(file, env)
	if err != nil {
		panic(err)
	}
}

const docsTemplate string = `<!--
File generated by internal/config/generate.
DO NOT EDIT.
-->

<!-- markdownlint-disable line_length -->
# Tribes Rollup Configuration

This file documents the configuration options.

<!-- markdownlint-disable MD012 -->
{{- range .}}

## {{backtick .Name}}

{{.Description}}

* **Type:** {{backtick .GoType}}
{{- if .Default}}
* **Default:** {{.Default | quote | backtick}}
{{- end}}
{{- if .UsedBy}}
* **Used by:** {{range $i, $e := .UsedBy}}{{if $i}}, {{end}}{{$e}}{{end}}
{{- end}}
{{- end}}
`
