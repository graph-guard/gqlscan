<a href="https://github.com/graph-guard/gqlscan/actions?query=workflow%3ACI">
    <img src="https://github.com/graph-guard/gqlscan/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href="https://coveralls.io/github/graph-guard/gqlscan">
    <img src="https://coveralls.io/repos/github/graph-guard/gqlscan/badge.svg" alt="Coverage Status" />
</a>
<a href="https://goreportcard.com/report/github.com/graph-guard/gqlscan">
    <img src="https://goreportcard.com/badge/github.com/graph-guard/gqlscan" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/graph-guard/gqlscan">
    <img src="https://godoc.org/github.com/graph-guard/gqlscan?status.svg" alt="GoDoc">
</a>
<a href="https://go.dev/play/p/hWgkDaNqrPr">
    <img src="https://img.shields.io/badge/Demo-Playground-blueviolet.svg">
</a>

# gqlscan
[gqlscan](https://pkg.go.dev/github.com/graph-guard/gqlscan) provides functions for fast and allocation-free
lexical scanning and validation of GraphQL queries according
to the [GraphQL specification of October 2021](https://spec.graphql.org/October2021/).

The provided functions don't perform semantic analysis such as
making sure that declared variables are used or that
values match their declared types, etc. as this is outside the scope
of lexical analysis.
