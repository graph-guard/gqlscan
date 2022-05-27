// Package gqlscan provides functions for fast and allocation-free
// lexical scanning and validation of GraphQL queries according
// to the GraphQL specification of October 2021
// (https://spec.graphql.org/October2021/).
//
// The provided functions don't perform semantic analysis such as
// making sure that declared variables are used or that
// values match their declared types, etc. as this is outside the scope
// of lexical analysis.
package gqlscan

//go:generate go run cmd/gen/main.go
