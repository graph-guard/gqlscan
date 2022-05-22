package gqlscan_test

import (
	_ "embed"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/graph-guard/gqlscan"

	"github.com/stretchr/testify/require"
)

type Expect struct {
	Decl  string
	Type  gqlscan.Token
	Value string
}

type TestInput struct {
	decl   string
	input  string
	expect []Expect
}

var testdata = []TestInput{
	{Decl(), `{foo}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "foo"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `query {foo}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "foo"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: {foo: false})}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "foo"},
		{Decl(), gqlscan.TokenFalse, "false"},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: false)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenFalse, "false"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: true)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenTrue, "true"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: null)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: [])}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: [[]])}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: 0)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: 0.0)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenFloat, "0.0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: 42)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenInt, "42"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: -42)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenInt, "-42"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: -42.5678)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenFloat, "-42.5678"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(f: -42.5678e2)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenFloat, "-42.5678e2"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{ f (f: {x: 2}) }`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "f"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "x"},
		{Decl(), gqlscan.TokenInt, "2"},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text done }`, []Expect{
		// Fragment f1
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f1"},
		{Decl(), gqlscan.TokenFragTypeCond, "Query"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "todos"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f2"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Query Todos
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenOprName, "Todos"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f1"},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f2"},
		{Decl(), gqlscan.TokenFragTypeCond, "Todo"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "id"},
		{Decl(), gqlscan.TokenField, "text"},
		{Decl(), gqlscan.TokenField, "done"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `query Q($variable: Foo, $v: [ [ Bar ] ]) {
		foo_alias: foo(x: null) {
			foobar_alias: foo_bar
		}
		bar
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	} mutation M($variable: Foo, $v: [ [ Bar ] ]) {
		foo(x: null) {
			foo_bar
		}
		bar
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	}
	fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text(
		foo: 2,
		bar: "ok",
		baz: null,
	) done }
	subscription S($v:Input!){
		sub(i: $v) {f}
	}`, []Expect{
		// Query Q
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenOprName, "Q"},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "variable"},
		{Decl(), gqlscan.TokenVarTypeName, "Foo"},
		{Decl(), gqlscan.TokenVarName, "v"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "Bar"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFieldAlias, "foo_alias"},
		{Decl(), gqlscan.TokenField, "foo"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "x"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFieldAlias, "foobar_alias"},
		{Decl(), gqlscan.TokenField, "foo_bar"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenField, "bar"},
		{Decl(), gqlscan.TokenField, "baz"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragInline, "A"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_A"},
		{Decl(), gqlscan.TokenFragRef, "namedFragment1"},
		{Decl(), gqlscan.TokenFragRef, "namedFragment2"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenFragInline, "B"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_B"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz1"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "bool"},
		{Decl(), gqlscan.TokenFalse, "false"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz2"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "bool"},
		{Decl(), gqlscan.TokenTrue, "true"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz3"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "string"},
		{Decl(), gqlscan.TokenStr, "okay"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz4"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "array"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz5"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "variable"},
		{Decl(), gqlscan.TokenVarRef, "variable"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz6"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "variable"},
		{Decl(), gqlscan.TokenVarRef, "v"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz7"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "object"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "number0"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenObjField, "number1"},
		{Decl(), gqlscan.TokenInt, "2"},
		{Decl(), gqlscan.TokenObjField, "number2"},
		{Decl(), gqlscan.TokenFloat, "123456789.1234e2"},

		{Decl(), gqlscan.TokenObjField, "arr0"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "x"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},

		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Mutation M
		{Decl(), gqlscan.TokenDefMut, ""},
		{Decl(), gqlscan.TokenOprName, "M"},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "variable"},
		{Decl(), gqlscan.TokenVarTypeName, "Foo"},
		{Decl(), gqlscan.TokenVarName, "v"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "Bar"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "foo"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "x"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "foo_bar"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenField, "bar"},
		{Decl(), gqlscan.TokenField, "baz"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragInline, "A"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_A"},
		{Decl(), gqlscan.TokenFragRef, "namedFragment1"},
		{Decl(), gqlscan.TokenFragRef, "namedFragment2"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenFragInline, "B"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_B"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz1"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "bool"},
		{Decl(), gqlscan.TokenFalse, "false"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz2"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "bool"},
		{Decl(), gqlscan.TokenTrue, "true"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz3"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "string"},
		{Decl(), gqlscan.TokenStr, "okay"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz4"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "array"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz5"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "variable"},
		{Decl(), gqlscan.TokenVarRef, "variable"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz6"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "variable"},
		{Decl(), gqlscan.TokenVarRef, "v"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz7"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "object"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "number0"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenObjField, "number1"},
		{Decl(), gqlscan.TokenInt, "2"},
		{Decl(), gqlscan.TokenObjField, "number2"},
		{Decl(), gqlscan.TokenFloat, "123456789.1234e2"},

		{Decl(), gqlscan.TokenObjField, "arr0"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "x"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},

		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Fragment f1
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f1"},
		{Decl(), gqlscan.TokenFragTypeCond, "Query"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "todos"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f2"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Query Todos
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenOprName, "Todos"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f1"},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f2"},
		{Decl(), gqlscan.TokenFragTypeCond, "Todo"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "id"},
		{Decl(), gqlscan.TokenField, "text"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "foo"},
		{Decl(), gqlscan.TokenInt, "2"},
		{Decl(), gqlscan.TokenArg, "bar"},
		{Decl(), gqlscan.TokenStr, "ok"},
		{Decl(), gqlscan.TokenArg, "baz"},
		{Decl(), gqlscan.TokenNull, "null"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenField, "done"},
		{Decl(), gqlscan.TokenSelEnd, ""},

		// Subscription S
		{Decl(), gqlscan.TokenDefSub, ""},
		{Decl(), gqlscan.TokenOprName, "S"},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "v"},
		{Decl(), gqlscan.TokenVarTypeName, "Input"},
		{Decl(), gqlscan.TokenVarTypeNotNull, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "sub"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "i"},
		{Decl(), gqlscan.TokenVarRef, "v"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},

	// Comments
	{Decl(), "  #comment1\n  #comment2\n  {x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{  #comment1\n  #comment2\n  x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x  #comment1\n  #comment2\n  }", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x}  #comment1\n  #comment2\n", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x(  #comment1\n  #comment2\n  y:0)}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x(y  #comment1\n  #comment2\n  :0)}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x(y:  #comment1\n  #comment2\n  0)}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x(y:0  #comment1\n  #comment2\n  )}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{x(y:0)  #comment1\n  #comment2\n  }", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query  #comment1\n  #comment2\n  {x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "mutation  #comment1\n  #comment2\n  {x}", []Expect{
		{Decl(), gqlscan.TokenDefMut, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "fragment  #comment1\n  #comment2\n  f on X{x}", []Expect{
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f"},
		{Decl(), gqlscan.TokenFragTypeCond, "X"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "fragment f  #comment1\n  #comment2\n  on X{x}", []Expect{
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f"},
		{Decl(), gqlscan.TokenFragTypeCond, "X"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "fragment f on  #comment1\n  #comment2\n  X{x}", []Expect{
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f"},
		{Decl(), gqlscan.TokenFragTypeCond, "X"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "fragment f on X  #comment1\n  #comment2\n  {x}", []Expect{
		{Decl(), gqlscan.TokenDefFrag, ""},
		{Decl(), gqlscan.TokenFragName, "f"},
		{Decl(), gqlscan.TokenFragTypeCond, "X"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{  ...  #comment1\n  #comment2\n  f  }", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{  ...  f  #comment1\n  #comment2\n  }", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFragRef, "f"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query(  #comment1\n  #comment2\n  $x: T){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x  #comment1\n  #comment2\n  : T){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x:  #comment1\n  #comment2\n  T){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x:[  #comment1\n  #comment2\n  T]){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x:[T  #comment1\n  #comment2\n  ]){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x:[T]  #comment1\n  #comment2\n  ){x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "query($x:[T])  #comment1\n  #comment2\n  {x}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{f#comment\n{f2}}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f2"},
		{Decl(), gqlscan.TokenSelEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},

	// String escape
	{Decl(), `{x(s:"\"")}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "s"},
		{Decl(), gqlscan.TokenStr, `\"`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{x(s:"\\")}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "s"},
		{Decl(), gqlscan.TokenStr, `\\`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{x(s:"\\\"")}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "s"},
		{Decl(), gqlscan.TokenStr, `\\\"`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},

	{Decl(), `{x(y:1e8)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenFloat, `1e8`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{x(y:0e8)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenFloat, `0e8`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{x(y:0e+8)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenFloat, `0e+8`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{x(y:0e-8)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "y"},
		{Decl(), gqlscan.TokenFloat, `0e-8`},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `mutation{x}`, []Expect{
		{Decl(), gqlscan.TokenDefMut, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `mutation($x:T){x}`, []Expect{
		{Decl(), gqlscan.TokenDefMut, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "x"},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `mutation M{x}`, []Expect{
		{Decl(), gqlscan.TokenDefMut, ""},
		{Decl(), gqlscan.TokenOprName, "M"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(o:{o2:{x:[]}})}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "o"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "o2"},
		{Decl(), gqlscan.TokenObj, ""},
		{Decl(), gqlscan.TokenObjField, "x"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenObjEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(a:[0])}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenArr, ""},
		{Decl(), gqlscan.TokenInt, "0"},
		{Decl(), gqlscan.TokenArrEnd, ""},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `query($v:T ! ){x(a:$v)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "v"},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeNotNull, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenVarRef, "v"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `query ($v: [ [ T ! ] ! ] ! ) {x(a:$v)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenVarList, ""},
		{Decl(), gqlscan.TokenVarName, "v"},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeArr, ""},
		{Decl(), gqlscan.TokenVarTypeName, "T"},
		{Decl(), gqlscan.TokenVarTypeNotNull, ""},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarTypeNotNull, ""},
		{Decl(), gqlscan.TokenVarTypeArrEnd, ""},
		{Decl(), gqlscan.TokenVarTypeNotNull, ""},
		{Decl(), gqlscan.TokenVarListEnd, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "x"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenVarRef, "v"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{ bob : alice }`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenFieldAlias, "bob"},
		{Decl(), gqlscan.TokenField, "alice"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	Input(`query # This is a test with many comments
	# sample comment text line
	{ # sample comment text line
		# sample comment text line
		a # sample comment text line
		# sample comment text line
		{ # sample comment text line
			# sample comment text line
			b # sample comment text line
			# sample comment text line
			( # sample comment text line
				# sample comment text line
				x # sample comment text line
				# sample comment text line
				: # sample comment text line
				# sample comment text line
				1 # sample comment text line
			# sample comment text line
			) # sample comment text line
			# sample comment text line
			{ # sample comment text line
				# sample comment text line
				c # sample comment text line
				# sample comment text line
				d # sample comment text line
			# sample comment text line
			} # sample comment text line
		# sample comment text line
		} # sample comment text line
	# sample comment text line
	} # sample comment text line
	# sample comment text line
	query # sample comment text line
	# sample comment text line
	( # sample comment text line
		# sample comment text line
		$v # sample comment text line
		# sample comment text line
		: # sample comment text line
		# sample comment text line
		T # sample comment text line
		# sample comment text line
		! # sample comment text line
		# sample comment text line
		@d1 # sample comment text line
		# sample comment text line
		@d2 # sample comment text line
		# sample comment text line
		( # sample comment text line
			# sample comment text line
			a # sample comment text line
			# sample comment text line
			: # sample comment text line
			# sample comment text line
			0 # sample comment text line
		# sample comment text line
		) # sample comment text line
		# sample comment text line
		@d3 # sample comment text line
	# sample comment text line
	) # sample comment text line
	# sample comment text line
	@d1 # sample comment text line
	# sample comment text line
	@d2 # sample comment text line
	# sample comment text line
	( # sample comment text line
		# sample comment text line
		a # sample comment text line
		# sample comment text line
		: # sample comment text line
		# sample comment text line
		0 # sample comment text line
	# sample comment text line
	) # sample comment text line
	# sample comment text line
	@d3 # sample comment text line
	# sample comment text line
	{ # sample comment text line
		# sample comment text line
		... # sample comment text line
		# sample comment text line
		on # sample comment text line
		# sample comment text line
		T # sample comment text line
		# sample comment text line
		{ # sample comment text line
		# sample comment text line
			# sample comment text line
			... # sample comment text line
			# sample comment text line
			f # sample comment text line
		# sample comment text line
		} # sample comment text line
	} # sample comment text line
	`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "a"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "b"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "x"),
		Token(gqlscan.TokenInt, "1"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "c"),
		Token(gqlscan.TokenField, "d"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragInline, "T"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
	),
	{Decl(), `{f}
		#0
		#01
		#012
		#0123
		#01234
		#012345
		#0123456
		#01234567
		#012345678
		#0123456789
		#01234567810
		#01234567810a
		#01234567810ab
		#01234567810abc
		#01234567810abcd`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(a:
		"\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000"
		b:123456789
	)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenStr,
			`\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000`},
		{Decl(), gqlscan.TokenArg, "b"},
		{Decl(), gqlscan.TokenInt, "123456789"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), "{f(a:" + string_2695b + ")}", []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenStr, string_2695b[1 : len(string_2695b)-1]},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `{f(
		a:""""""
		b:"""abc"""
		c:"""\n\t" """
		d:"""
			foo
				bar
		"""
	)}`, []Expect{
		{Decl(), gqlscan.TokenDefQry, ""},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenArgList, ""},
		{Decl(), gqlscan.TokenArg, "a"},
		{Decl(), gqlscan.TokenStrBlock, ""},
		{Decl(), gqlscan.TokenArg, "b"},
		{Decl(), gqlscan.TokenStrBlock, "abc"},
		{Decl(), gqlscan.TokenArg, "c"},
		{Decl(), gqlscan.TokenStrBlock, `\n\t" `},
		{Decl(), gqlscan.TokenArg, "d"},
		{Decl(), gqlscan.TokenStrBlock,
			"\n\t\t\tfoo\n\t\t\t\tbar\n\t\t"},
		{Decl(), gqlscan.TokenArgListEnd, ""},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	{Decl(), `subscription S{f}`, []Expect{
		{Decl(), gqlscan.TokenDefSub, ""},
		{Decl(), gqlscan.TokenOprName, "S"},
		{Decl(), gqlscan.TokenSel, ""},
		{Decl(), gqlscan.TokenField, "f"},
		{Decl(), gqlscan.TokenSelEnd, ""},
	}},
	Input(`mutation @d1 @d2 (a:0) @d3 {f}`,
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`subscription @d1 @d2 (a:0) @d3 {f}`,
		Token(gqlscan.TokenDefSub),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query @d1 @d2 (a:0) @d3 {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query Q @d1 @d2 (a:0) @d3 {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Q"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query ($v: String) @d1 @d2 (a:$v) @d3 {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query @d1 @d2 (a:$v) {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query (
		$v: String @d1 @d2 (a:$v) @d3
	) {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query (
		$v1: String @d1 @d2 (a:0)
		$v2: String! @d1 @d2 (a:0)
		$v3: [String] @d1
	) {f}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v1"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenVarName, "v2"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenVarName, "v3"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{
		a (a: 0) @d1 @d2 (a:$v) @d3 {
			aa (a: 0) @d1 @d2 (a:$v) @d3
		}
		b @d1 @d2 (a:$v) @d3 {
			ba @d2 (a:$v)
			bb @d3 { bba }
		}
	}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "a"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "aa"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenField, "b"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "ba"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "bb"),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "bba"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{
		...f @d1 @d2 (a:$v) @d3
		...f2 @d1 @d2 (a:$v)
		x
		... on X @d1 @d2 (a:$v) @d3 {
			x
		}
		... on Y @d1 @d2 (a:$v) {
			x
		}
	}
	fragment f on X @d1 @d2 (a:$v) @d3 { x }
	fragment f2 on Y @d1 @d2 (a:$v) { x }`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenFragRef, "f2"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "x"),

		Token(gqlscan.TokenFragInline, "X"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenFragInline, "Y"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f"),
		Token(gqlscan.TokenFragTypeCond, "X"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenDirName, "d3"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f2"),
		Token(gqlscan.TokenFragTypeCond, "Y"),
		Token(gqlscan.TokenDirName, "d1"),
		Token(gqlscan.TokenDirName, "d2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
}

//go:embed t_s_2695b.txt
var string_2695b string

//go:embed t_blks_2747b.graphql
var blockstring_2747b string

//go:embed t_blks_2747b_expect.txt
var blockstring_2747b_expect string

func TestScan(t *testing.T) {
	for _, td := range testdata {
		t.Run(td.decl, func(t *testing.T) {
			t.Run("Scan", func(t *testing.T) {
				j := 0
				prevHead := 0
				err := gqlscan.Scan(
					[]byte(td.input),
					func(i *gqlscan.Iterator) (err bool) {
						require.True(
							t, j < len(td.expect),
							"exceeding expectation set at: %d {T: %s; V: %s}",
							j, i.Token().String(), i.Value(),
						)
						require.Equal(
							t, td.expect[j].Type.String(), i.Token().String(),
							"unexpected type at index %d (%s)",
							j, td.expect[j].Decl,
						)
						require.Equal(
							t, td.expect[j].Value, string(i.Value()),
							"unexpected type at index %d (%s)",
							j, td.expect[j].Decl,
						)
						require.GreaterOrEqual(t, i.IndexHead(), prevHead)
						require.GreaterOrEqual(t, i.IndexHead(), i.IndexTail())
						i.Value()
						j++
						return false
					},
				)
				require.Zero(t, err.Error())
				require.False(t, err.IsErr())
				for _, e := range td.expect[j:] {
					t.Errorf(
						"missing {T: %s; V: %s}",
						e.Type, e.Value,
					)
				}
			})

			t.Run("ScanAll", func(t *testing.T) {
				j := 0
				prevHead := 0
				err := gqlscan.ScanAll(
					[]byte(td.input),
					func(i *gqlscan.Iterator) {
						require.True(
							t, j < len(td.expect),
							"exceeding expectation set at: %d (%s) {T: %s; V: %s}",
							j, td.expect[j].Decl, i.Token().String(), i.Value(),
						)
						require.Equal(
							t, td.expect[j].Type.String(), i.Token().String(),
							"unexpected type at index %d (%s)",
							j, td.expect[j].Decl,
						)
						require.Equal(
							t, td.expect[j].Value, string(i.Value()),
							"unexpected value at index %d (%s)",
							j, td.expect[j].Decl,
						)
						require.GreaterOrEqual(t, i.IndexHead(), prevHead)
						require.GreaterOrEqual(t, i.IndexHead(), i.IndexTail())
						i.Value()
						j++
					},
				)
				require.Zero(t, err.Error())
				require.False(t, err.IsErr())
				for _, e := range td.expect[j:] {
					t.Errorf(
						"missing {T: %s; V: %s}",
						e.Type, e.Value,
					)
				}
			})
		})
	}
}

type TestInputErr struct {
	decl      string
	input     string
	expectErr string
}

var testdataErr = []TestInputErr{
	{Decl(), // Unexpected token as query.
		"q",
		"error at index 0 ('q'): unexpected token; expected definition",
	},
	{Decl(), // Missing square bracket in type.
		"query($a: [A){f}",
		"error at index 11 ('A'): invalid type; " +
			"expected variable type",
	},
	{Decl(), // Missing square bracket in type.
		"query($a: [[A]){f}",
		"error at index 13 (']'): invalid type; " +
			"expected variable type",
	},
	{Decl(), // Unexpected square bracket in variable type.
		"query($a: A]){f}",
		"error at index 11 (']'): unexpected token; " +
			"expected variable name",
	},
	{Decl(), // Unexpected square bracket in variable type.
		"query($a: [[A]]]){f}",
		"error at index 15 (']'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Missing query closing curly bracket.
		"{",
		"error at index 1: unexpected end of file; expected selection",
	},
	{Decl(), // Invalid field name.
		"{1abc}",
		"error at index 1 ('1'): unexpected token; " +
			"expected field name or alias",
	},
	{Decl(), // Trailing closing curly bracket.
		"{f}}",
		"error at index 3 ('}'): unexpected token; expected definition",
	},
	{Decl(), // Query missing closing curly bracket.
		"{}",
		"error at index 1 ('}'): unexpected token; " +
			"expected field name or alias",
	},
	{Decl(), // Variable missing name.
		"query($ ",
		"error at index 7 (' '): unexpected token; " +
			"expected variable name",
	},
	{Decl(), // Empty args.
		"{f()}",
		"error at index 3 (')'): unexpected token; expected argument name",
	},
	{Decl(), // Argument missing column.
		"{f(x null))}",
		"error at index 5 ('n'): " +
			"unexpected token; expected column after argument name",
	},
	{Decl(), // Argument with trailing closing parenthesis.
		"{f(x:null))}",
		"error at index 10 (')'): unexpected token; " +
			"expected field name or alias",
	},
	{Decl(), // Argument missing closing parenthesis.
		"{f(}",
		"error at index 3 ('}'): unexpected token; expected argument name",
	},
	{Decl(), // Invalid argument.
		"{f(x:abc))}",
		"error at index 5 ('a'): invalid value; expected value",
	},
	{Decl(), // String argument missing closing quotes.
		`{f(x:"))}`,
		"error at index 9: unexpected end of file; expected end of string",
	},
	{Decl(), // Invalid negative number.
		`{f(x:-))}`,
		"error at index 6 (')'): invalid number value; expected value",
	},
	{Decl(), // Number missing fraction.
		`{f(x:1.))}`,
		"error at index 7 (')'): invalid number value; expected value",
	},
	{Decl(), // Number missing exponent.
		`{f(x:1.2e))}`,
		"error at index 9 (')'): invalid number value; expected value",
	},
	{Decl(), // Number with leading zero.
		`{f(x:0123))}`,
		"error at index 6 ('1'): invalid number value; expected value",
	},

	// --- Unexpected EOF ---
	{Decl(), // Unexpected EOF.
		"",
		"error at index 0: unexpected end of file; expected definition",
	},
	{Decl(), // Unexpected EOF.
		"query",
		"error at index 5: unexpected end of file; " +
			"expected variable list or selection set",
	},
	{Decl(), // Unexpected EOF.
		"query Name",
		"error at index 10: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"query Name ",
		"error at index 11: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"mutation Name",
		"error at index 13: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"mutation Name ",
		"error at index 14: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"subscription Name",
		"error at index 17: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"subscription Name ",
		"error at index 18: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"query(",
		"error at index 6: unexpected end of file; " +
			"expected variable name",
	},
	{Decl(), // Unexpected EOF.
		"query( ",
		"error at index 7: unexpected end of file; " +
			"expected variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($",
		"error at index 6: unexpected end of file; " +
			"expected variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($v",
		"error at index 8: unexpected end of file; " +
			"expected column after variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($v ",
		"error at index 9: unexpected end of file; " +
			"expected column after variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($v:",
		"error at index 9: unexpected end of file; " +
			"expected variable type",
	},
	{Decl(), // Unexpected EOF.
		"query($v: ",
		"error at index 10: unexpected end of file; " +
			"expected variable type",
	},
	{Decl(), // Unexpected EOF.
		"query($v: T",
		"error at index 11: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($v: T ",
		"error at index 12: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected EOF.
		"query($v: T)",
		"error at index 12: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"query($v: T) ",
		"error at index 13: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"{",
		"error at index 1: unexpected end of file; " +
			"expected selection",
	},
	{Decl(), // Unexpected EOF.
		"{ ",
		"error at index 2: unexpected end of file; " +
			"expected selection",
	},
	{Decl(), // Unexpected EOF.
		"{foo",
		"error at index 4: unexpected end of file; " +
			"expected field name or alias",
	},
	{Decl(), // Unexpected EOF.
		"{foo ",
		"error at index 5: unexpected end of file; " +
			"expected field name or alias",
	},
	{Decl(), // Unexpected EOF.
		"{foo(",
		"error at index 5: unexpected end of file; " +
			"expected argument name",
	},
	{Decl(), // Unexpected EOF.
		"{foo( ",
		"error at index 6: unexpected end of file; " +
			"expected argument name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name",
		"error at index 9: unexpected end of file; " +
			"expected column after argument name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name ",
		"error at index 10: unexpected end of file; " +
			"expected column after argument name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name:",
		"error at index 10: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: ",
		"error at index 11: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: {",
		"error at index 12: unexpected end of file; " +
			"expected object field name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: { ",
		"error at index 13: unexpected end of file; " +
			"expected object field name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: {field",
		"error at index 17: unexpected end of file; " +
			"expected column after object field name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: {field ",
		"error at index 18: unexpected end of file; " +
			"expected column after object field name",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: {field:",
		"error at index 18: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		"{foo(name: {field: ",
		"error at index 19: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: "`,
		"error at index 12: unexpected end of file; " +
			"expected end of string",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: ""`,
		"error at index 13: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: f`,
		"error at index 12: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: t`,
		"error at index 12: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: n`,
		"error at index 12: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: 0`,
		"error at index 12: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: 0 `,
		"error at index 13: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: -`,
		"error at index 12: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: 0.`,
		"error at index 13: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{foo(name: 0.1e`,
		"error at index 15: unexpected end of file; expected value",
	},
	{Decl(), // Unexpected EOF.
		`{.`,
		"error at index 2: unexpected end of file; expected fragment",
	},
	{Decl(), // Unexpected EOF.
		`{..`,
		"error at index 3: unexpected end of file; expected fragment",
	},
	{Decl(), // Unexpected EOF.
		`{...`,
		"error at index 4: unexpected end of file; expected fragment",
	},
	{Decl(), // Unexpected EOF.
		`{... `,
		"error at index 5: unexpected end of file; expected fragment",
	},
	{Decl(), // Unexpected EOF.
		`{... on`,
		"error at index 7: unexpected end of file; expected fragment",
	},
	{Decl(), // Unexpected EOF.
		`{... on `,
		"error at index 8: unexpected end of file; " +
			"expected inlined fragment",
	},
	{Decl(), // Unexpected EOF.
		"fragment",
		"error at index 8: unexpected end of file; " +
			"expected fragment name",
	},
	{Decl(), // Unexpected EOF.
		"{x",
		"error at index 2: unexpected end of file; " +
			"expected field name or alias",
	},

	{Decl(), // Invalid value.
		"{x(p:falsa",
		"error at index 5 ('f'): invalid value; " +
			"expected value",
	},
	{Decl(), // Invalid value.
		"{x(p:truu",
		"error at index 5 ('t'): invalid value; " +
			"expected value",
	},
	{Decl(), // Invalid value.
		"{x(p:nuli",
		"error at index 5 ('n'): invalid value; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		"{x(p:[",
		"error at index 6: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"query($x:T)x",
		"error at index 11 ('x'): unexpected token; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"mutation M",
		"error at index 10: unexpected end of file; " +
			"expected selection set",
	},
	{Decl(), // Unexpected token.
		"query\x00",
		"error at index 5 (0x0): unexpected token; " +
			"expected operation name",
	},
	{Decl(), // Unexpected token.
		"{x(y:12e)}",
		"error at index 8 (')'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:12.)}",
		"error at index 8 (')'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:12x)}",
		"error at index 7 ('x'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:12.12x)}",
		"error at index 10 ('x'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected EOF.
		"{x(y:12.12",
		"error at index 10: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{Decl(), // Unexpected EOF.
		"{x(y:12.",
		"error at index 8: unexpected end of file; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:12e111x",
		"error at index 11 ('x'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:12ex",
		"error at index 8 ('x'): invalid number value; " +
			"expected value",
	},
	{Decl(), // Unexpected token.
		"{x(y:{f})}",
		"error at index 7 ('}'): unexpected token; " +
			"expected column after object field name",
	},
	{Decl(), // Unexpected token.
		"{x(\x00:1)}",
		"error at index 3 (0x0): unexpected token; " +
			"expected argument name",
	},
	{Decl(), // Unexpected EOF.
		"{x(y\x00:1)}",
		"error at index 4 (0x0): unexpected token; " +
			"expected argument name",
	},
	{Decl(), // Unexpected token.
		"query M [",
		"error at index 8 ('['): unexpected token; " +
			"expected selection set",
	},
	{Decl(), // Unexpected token.
		"mutation M|",
		"error at index 10 ('|'): unexpected token; " +
			"expected selection set",
	},
	{Decl(), // Unexpected EOF.
		"fragment f on",
		"error at index 13: unexpected end of file; " +
			"expected fragment type condition",
	},
	{Decl(), // Unexpected token.
		"mutation\x00",
		"error at index 8 (0x0): unexpected token; " +
			"expected operation name",
	},
	{Decl(), // Unexpected token.
		"subscription\x00",
		"error at index 12 (0x0): unexpected token; " +
			"expected operation name",
	},
	{Decl(), // Unexpected token.
		"fragment\x00",
		"error at index 8 (0x0): unexpected token; " +
			"expected fragment name",
	},
	{Decl(), // Unexpected EOF.
		"{x(y:$",
		"error at index 6: unexpected end of file; " +
			"expected referenced variable name",
	},
	{Decl(), // Unexpected EOF.
		"mutation",
		"error at index 8: unexpected end of file; " +
			"expected variable list or selection set",
	},
	{Decl(), // Unexpected EOF.
		"{x(y:null)",
		"error at index 10: unexpected end of file; " +
			"expected selection set or selection",
	},
	{Decl(), // Unexpected token.
		"query($v |",
		"error at index 9 ('|'): unexpected token; " +
			"expected column after variable name",
	},
	{Decl(), // Unexpected token.
		"query($v:[T] |)",
		"error at index 13 ('|'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected token.
		"fragment X at",
		"error at index 11 ('a'): unexpected token; " +
			"expected keyword 'on'",
	},
	{Decl(), // Unexpected EOF.
		"query($a:[A]",
		"error at index 12: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected EOF.
		"fragment f ",
		"error at index 11: unexpected end of file; " +
			"expected keyword 'on'",
	},
	{Decl(), // Unexpected EOF.
		"{f{x} ",
		"error at index 6: unexpected end of file; " +
			"expected selection or end of selection set",
	},
	{Decl(), // Unexpected token.
		"{f(x:\"abc\n\")}",
		"error at index 9 (0xa): unexpected token; " +
			"expected end of string",
	},
	{Decl(), // Unexpected token.
		"{.f}",
		"error at index 2 ('f'): unexpected token; " +
			"expected fragment",
	},
	{Decl(), // Unexpected token.
		"{..f}",
		"error at index 3 ('f'): unexpected token; " +
			"expected fragment",
	},
	{Decl(), // Unexpected token.
		"query($v:T ! !){x(a:$v)}",
		"error at index 13 ('!'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected token.
		"query($v: [ T ! ] ! ! ){x(a:$v)}",
		"error at index 20 ('!'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{Decl(), // Unexpected token.
		"{alias : alias2 : x}",
		"error at index 16 (':'): unexpected token; " +
			"expected field name or alias",
	},
	{Decl(), // Unexpected EOF.
		"{f:",
		"error at index 3: unexpected end of file; " +
			"expected field name",
	},
	{Decl(), // Unexpected EOF.
		"{f: ",
		"error at index 4: unexpected end of file; " +
			"expected field name",
	},
	{Decl(), // Invalid escape sequence.
		`{f(a:"\a")}`,
		"error at index 7 ('a'): unexpected token; " +
			"expected escaped sequence",
	},
	{Decl(), // Invalid escape sequence.
		`{f(a:"\u")}`,
		"error at index 8 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Invalid escape sequence.
		`{f(a:"\u1")}`,
		"error at index 9 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Invalid escape sequence.
		`{f(a:"\u12")}`,
		"error at index 10 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:"\u`,
		"error at index 8: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:"\u1`,
		"error at index 9: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:"\u12`,
		"error at index 10: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:"\u123`,
		"error at index 11: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Invalid escape sequence.
		`{f(a:"\u123")}`,
		"error at index 11 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:"""`,
		`error at index 8: unexpected end of file; ` +
			"expected end of block string",
	},
	{Decl(), // Unexpected EOF.
		`{f(a:""" `,
		"error at index 9: unexpected end of file; " +
			"expected end of block string",
	},
	{Decl(), // Control character in string.
		`{f(a:"0123456` + string(rune(0x00)) + `")}`,
		"error at index 13 (0x0): unexpected token; " +
			"expected end of string",
	},
	{Decl(), // Unexpected EOF.
		`{f #c`,
		"error at index 5: unexpected end of file; " +
			"expected selection, selection set or end of selection set",
	},
	InputErr( // Unexpected EOF.
		"query @",
		"error at index 7: unexpected end of file; "+
			"expected directive name",
	),
	InputErr( // Unexpected EOF.
		"query @ ",
		"error at index 8: unexpected end of file; "+
			"expected directive name",
	),
	InputErr( // Unexpected EOF.
		"query @directive",
		"error at index 16: unexpected end of file; "+
			"expected variable list or selection set",
	),
	InputErr( // Unexpected EOF.
		"query @directive ",
		"error at index 17: unexpected end of file; "+
			"expected variable list or selection set",
	),
	InputErr( // Unexpected EOF.
		"query @directive(",
		"error at index 17: unexpected end of file; "+
			"expected argument name",
	),
	InputErr( // Unexpected EOF.
		"query @directive( ",
		"error at index 18: unexpected end of file; "+
			"expected argument name",
	),
	InputErr( // Unexpected EOF.
		"{f @",
		"error at index 4: unexpected end of file; "+
			"expected directive name",
	),
	InputErr( // Unexpected EOF.
		"{f @ ",
		"error at index 5: unexpected end of file; "+
			"expected directive name",
	),
	InputErr( // Unexpected EOF.
		"{f @d",
		"error at index 5: unexpected end of file; "+
			"expected selection, selection set or end of selection set",
	),
	InputErr( // Unexpected EOF.
		"{f @d",
		"error at index 5: unexpected end of file; "+
			"expected selection, selection set or end of selection set",
	),
	InputErr( // Unexpected token; Variables after directives
		"query @d (a:0) (a:0) {f}",
		"error at index 15 ('('): unexpected token; "+
			"expected selection set",
	),
	InputErr( // Unexpected token; Arguments after directives
		"{f @d (a:0) (a:0)}",
		"error at index 12 ('('): unexpected token; "+
			"expected field name or alias",
	),
}

func TestScanErr(t *testing.T) {
	for _, td := range testdataErr {
		t.Run(td.decl, func(t *testing.T) {
			t.Run("Scan", func(t *testing.T) {
				err := gqlscan.Scan(
					[]byte(td.input),
					func(*gqlscan.Iterator) (err bool) {
						return false
					},
				)
				require.Equal(t, td.expectErr, err.Error())
				require.True(t, err.IsErr())
			})

			t.Run("ScanAll", func(t *testing.T) {
				err := gqlscan.ScanAll(
					[]byte(td.input),
					func(*gqlscan.Iterator) {},
				)
				require.Equal(t, td.expectErr, err.Error())
				require.True(t, err.IsErr())
			})
		})
	}
}

func TestScanFuncErr(t *testing.T) {
	const input = `
		{x}
		query($v: [T!]!) {x}
		mutation($v: [T!]!) {x}
		query Q($variable: Foo, $v: [ [ Bar ] ]) {
		foo(x: null) {
			foo_bar
		}
		bar
		tar(x: """block string""")
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	} mutation M($variable: Foo, $v: [ [ Bar ] ]) {
		foo_alias: foo(x: null) {
			foobar_alias: foo_bar
		}
		bar
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	}
	fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text(
		foo: 2,
		bar: "ok",
		baz: null,
	) done }
	subscription S($v:Input!){
		sub(i: $v) {f}
	}`

	numOfTokensInInput := 0
	err := gqlscan.Scan([]byte(input), func(*gqlscan.Iterator) (err bool) {
		numOfTokensInInput++
		return false
	})
	require.False(t, err.IsErr())

	for i := 0; i < numOfTokensInInput; i++ {
		c := 0
		err := gqlscan.Scan(
			[]byte(input),
			func(j *gqlscan.Iterator) (err bool) {
				if c == i {
					return true
				}
				c++
				return false
			},
		)
		require.True(t, err.IsErr())
		require.Equal(
			t, gqlscan.ErrCallbackFn, err.Code,
			"at index %d", err.AtIndex,
		)
		require.Regexp(t,
			"error at index [0-9]*(\\s+\\(.*\\))?: "+
				"callback function returned error; expected .*$",
			err.Error(),
		)
	}
}

func TestLevel(t *testing.T) {
	const input = `query Q($variable: Foo, $v: [ [ Bar ] ]) {
		foo_alias: foo(x: null) {
			foobar_alias: foo_bar
		}
		bar
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	} mutation M($variable: Foo, $v: [ [ Bar ] ]) {
		foo(x: null) {
			foo_bar
		}
		bar
		baz {
			baz_fuzz {
				... on A {
					baz_fuzz_taz_A
					...namedFragment1
					... namedFragment2
				}
				... on B {
					baz_fuzz_taz_B
				}
				baz_fuzz_taz1(bool: false)
				baz_fuzz_taz2(bool: true)
				baz_fuzz_taz3(string: "okay")
				baz_fuzz_taz4(array: [])
				baz_fuzz_taz5(variable: $variable)
				baz_fuzz_taz6(variable: $v)
				baz_fuzz_taz7(object: {
					number0: 0
					number1: 2
					number2: 123456789.1234e2
					arr0: [[] [{x:null}]]
				})
			}
		}
	}
	fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text(
		foo: 2,
		bar: "ok",
		baz: null,
	) done }`

	expect := []struct {
		Decl  string
		Type  gqlscan.Token
		Value string
		Level int
	}{
		// Query Q
		{Decl(), gqlscan.TokenDefQry, "", 0},
		{Decl(), gqlscan.TokenOprName, "Q", 0},
		{Decl(), gqlscan.TokenVarList, "", 0},
		{Decl(), gqlscan.TokenVarName, "variable", 0},
		{Decl(), gqlscan.TokenVarTypeName, "Foo", 0},
		{Decl(), gqlscan.TokenVarName, "v", 0},
		{Decl(), gqlscan.TokenVarTypeArr, "", 0},
		{Decl(), gqlscan.TokenVarTypeArr, "", 0},
		{Decl(), gqlscan.TokenVarTypeName, "Bar", 0},
		{Decl(), gqlscan.TokenVarTypeArrEnd, "", 0},
		{Decl(), gqlscan.TokenVarTypeArrEnd, "", 0},
		{Decl(), gqlscan.TokenVarListEnd, "", 0},
		{Decl(), gqlscan.TokenSel, "", 0},
		{Decl(), gqlscan.TokenFieldAlias, "foo_alias", 1},
		{Decl(), gqlscan.TokenField, "foo", 1},
		{Decl(), gqlscan.TokenArgList, "", 1},
		{Decl(), gqlscan.TokenArg, "x", 1},
		{Decl(), gqlscan.TokenNull, "null", 1},
		{Decl(), gqlscan.TokenArgListEnd, "", 1},

		{Decl(), gqlscan.TokenSel, "", 1},
		{Decl(), gqlscan.TokenFieldAlias, "foobar_alias", 2},
		{Decl(), gqlscan.TokenField, "foo_bar", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 2},
		{Decl(), gqlscan.TokenField, "bar", 1},
		{Decl(), gqlscan.TokenField, "baz", 1},
		{Decl(), gqlscan.TokenSel, "", 1},
		{Decl(), gqlscan.TokenField, "baz_fuzz", 2},
		{Decl(), gqlscan.TokenSel, "", 2},
		{Decl(), gqlscan.TokenFragInline, "A", 3},
		{Decl(), gqlscan.TokenSel, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{Decl(), gqlscan.TokenFragRef, "namedFragment1", 4},
		{Decl(), gqlscan.TokenFragRef, "namedFragment2", 4},
		{Decl(), gqlscan.TokenSelEnd, "", 4},
		{Decl(), gqlscan.TokenFragInline, "B", 3},
		{Decl(), gqlscan.TokenSel, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{Decl(), gqlscan.TokenSelEnd, "", 4},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "bool", 3},
		{Decl(), gqlscan.TokenFalse, "false", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "bool", 3},
		{Decl(), gqlscan.TokenTrue, "true", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "string", 3},
		{Decl(), gqlscan.TokenStr, "okay", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "array", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "variable", 3},
		{Decl(), gqlscan.TokenVarRef, "variable", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "variable", 3},
		{Decl(), gqlscan.TokenVarRef, "v", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "object", 3},
		{Decl(), gqlscan.TokenObj, "", 3},
		{Decl(), gqlscan.TokenObjField, "number0", 3},
		{Decl(), gqlscan.TokenInt, "0", 3},
		{Decl(), gqlscan.TokenObjField, "number1", 3},
		{Decl(), gqlscan.TokenInt, "2", 3},
		{Decl(), gqlscan.TokenObjField, "number2", 3},
		{Decl(), gqlscan.TokenFloat, "123456789.1234e2", 3},

		{Decl(), gqlscan.TokenObjField, "arr0", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenObj, "", 3},
		{Decl(), gqlscan.TokenObjField, "x", 3},
		{Decl(), gqlscan.TokenNull, "null", 3},
		{Decl(), gqlscan.TokenObjEnd, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},

		{Decl(), gqlscan.TokenObjEnd, "", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenSelEnd, "", 3},
		{Decl(), gqlscan.TokenSelEnd, "", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 1},

		// Mutation M
		{Decl(), gqlscan.TokenDefMut, "", 0},
		{Decl(), gqlscan.TokenOprName, "M", 0},
		{Decl(), gqlscan.TokenVarList, "", 0},
		{Decl(), gqlscan.TokenVarName, "variable", 0},
		{Decl(), gqlscan.TokenVarTypeName, "Foo", 0},
		{Decl(), gqlscan.TokenVarName, "v", 0},
		{Decl(), gqlscan.TokenVarTypeArr, "", 0},
		{Decl(), gqlscan.TokenVarTypeArr, "", 0},
		{Decl(), gqlscan.TokenVarTypeName, "Bar", 0},
		{Decl(), gqlscan.TokenVarTypeArrEnd, "", 0},
		{Decl(), gqlscan.TokenVarTypeArrEnd, "", 0},
		{Decl(), gqlscan.TokenVarListEnd, "", 0},
		{Decl(), gqlscan.TokenSel, "", 0},
		{Decl(), gqlscan.TokenField, "foo", 1},
		{Decl(), gqlscan.TokenArgList, "", 1},
		{Decl(), gqlscan.TokenArg, "x", 1},
		{Decl(), gqlscan.TokenNull, "null", 1},
		{Decl(), gqlscan.TokenArgListEnd, "", 1},
		{Decl(), gqlscan.TokenSel, "", 1},
		{Decl(), gqlscan.TokenField, "foo_bar", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 2},
		{Decl(), gqlscan.TokenField, "bar", 1},
		{Decl(), gqlscan.TokenField, "baz", 1},
		{Decl(), gqlscan.TokenSel, "", 1},
		{Decl(), gqlscan.TokenField, "baz_fuzz", 2},
		{Decl(), gqlscan.TokenSel, "", 2},
		{Decl(), gqlscan.TokenFragInline, "A", 3},
		{Decl(), gqlscan.TokenSel, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{Decl(), gqlscan.TokenFragRef, "namedFragment1", 4},
		{Decl(), gqlscan.TokenFragRef, "namedFragment2", 4},
		{Decl(), gqlscan.TokenSelEnd, "", 4},
		{Decl(), gqlscan.TokenFragInline, "B", 3},
		{Decl(), gqlscan.TokenSel, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{Decl(), gqlscan.TokenSelEnd, "", 4},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "bool", 3},
		{Decl(), gqlscan.TokenFalse, "false", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "bool", 3},
		{Decl(), gqlscan.TokenTrue, "true", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "string", 3},
		{Decl(), gqlscan.TokenStr, "okay", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "array", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "variable", 3},
		{Decl(), gqlscan.TokenVarRef, "variable", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "variable", 3},
		{Decl(), gqlscan.TokenVarRef, "v", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{Decl(), gqlscan.TokenArgList, "", 3},
		{Decl(), gqlscan.TokenArg, "object", 3},
		{Decl(), gqlscan.TokenObj, "", 3},
		{Decl(), gqlscan.TokenObjField, "number0", 3},
		{Decl(), gqlscan.TokenInt, "0", 3},
		{Decl(), gqlscan.TokenObjField, "number1", 3},
		{Decl(), gqlscan.TokenInt, "2", 3},
		{Decl(), gqlscan.TokenObjField, "number2", 3},
		{Decl(), gqlscan.TokenFloat, "123456789.1234e2", 3},

		{Decl(), gqlscan.TokenObjField, "arr0", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArr, "", 3},
		{Decl(), gqlscan.TokenObj, "", 3},
		{Decl(), gqlscan.TokenObjField, "x", 3},
		{Decl(), gqlscan.TokenNull, "null", 3},
		{Decl(), gqlscan.TokenObjEnd, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},
		{Decl(), gqlscan.TokenArrEnd, "", 3},

		{Decl(), gqlscan.TokenObjEnd, "", 3},
		{Decl(), gqlscan.TokenArgListEnd, "", 3},
		{Decl(), gqlscan.TokenSelEnd, "", 3},
		{Decl(), gqlscan.TokenSelEnd, "", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 1},

		// Fragment f1
		{Decl(), gqlscan.TokenDefFrag, "", 0},
		{Decl(), gqlscan.TokenFragName, "f1", 0},
		{Decl(), gqlscan.TokenFragTypeCond, "Query", 0},
		{Decl(), gqlscan.TokenSel, "", 0},
		{Decl(), gqlscan.TokenField, "todos", 1},
		{Decl(), gqlscan.TokenSel, "", 1},
		{Decl(), gqlscan.TokenFragRef, "f2", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 2},
		{Decl(), gqlscan.TokenSelEnd, "", 1},

		// Query Todos
		{Decl(), gqlscan.TokenDefQry, "", 0},
		{Decl(), gqlscan.TokenOprName, "Todos", 0},
		{Decl(), gqlscan.TokenSel, "", 0},
		{Decl(), gqlscan.TokenFragRef, "f1", 1},
		{Decl(), gqlscan.TokenSelEnd, "", 1},

		// Fragment f2
		{Decl(), gqlscan.TokenDefFrag, "", 0},
		{Decl(), gqlscan.TokenFragName, "f2", 0},
		{Decl(), gqlscan.TokenFragTypeCond, "Todo", 0},
		{Decl(), gqlscan.TokenSel, "", 0},
		{Decl(), gqlscan.TokenField, "id", 1},
		{Decl(), gqlscan.TokenField, "text", 1},
		{Decl(), gqlscan.TokenArgList, "", 1},
		{Decl(), gqlscan.TokenArg, "foo", 1},
		{Decl(), gqlscan.TokenInt, "2", 1},
		{Decl(), gqlscan.TokenArg, "bar", 1},
		{Decl(), gqlscan.TokenStr, "ok", 1},
		{Decl(), gqlscan.TokenArg, "baz", 1},
		{Decl(), gqlscan.TokenNull, "null", 1},
		{Decl(), gqlscan.TokenArgListEnd, "", 1},
		{Decl(), gqlscan.TokenField, "done", 1},
		{Decl(), gqlscan.TokenSelEnd, "", 1},
	}

	t.Run("Scan", func(t *testing.T) {
		j := 0
		err := gqlscan.Scan(
			[]byte(input),
			func(i *gqlscan.Iterator) (err bool) {
				require.True(
					t, j < len(expect),
					"exceeding expectation set at: %d (%s) {T: %s; V: %s}",
					j, expect[j].Decl, i.Token().String(), i.Value(),
				)
				require.Equal(
					t, expect[j].Type.String(), i.Token().String(),
					"unexpected type at index %d (%s)",
					j, expect[j].Decl,
				)
				require.Equal(
					t, expect[j].Value, string(i.Value()),
					"unexpected value at index %d (%s)",
					j, expect[j].Decl,
				)
				require.Equal(
					t, expect[j].Level, i.LevelSelect(),
					"unexpected selection level at index %d (%s)",
					j, expect[j].Decl,
				)
				i.Value()
				j++
				return false
			},
		)
		require.Zero(t, err.Error())
		require.False(t, err.IsErr())
		for _, e := range expect[j:] {
			t.Errorf(
				"missing {T: %s; V: %s}",
				e.Type, e.Value,
			)
		}
	})

	t.Run("ScanAll", func(t *testing.T) {
		j := 0
		err := gqlscan.ScanAll(
			[]byte(input),
			func(i *gqlscan.Iterator) {
				require.True(
					t, j < len(expect),
					"exceeding expectation set at: %d (%s) {T: %s; V: %s}",
					j, expect[j].Decl, i.Token().String(), i.Value(),
				)
				require.Equal(
					t, expect[j].Type.String(), i.Token().String(),
					"unexpected type at index %d (%s)",
					j, expect[j].Decl,
				)
				require.Equal(
					t, expect[j].Value, string(i.Value()),
					"unexpected value at index %d (%s)",
					j, expect[j].Decl,
				)
				require.Equal(
					t, expect[j].Level, i.LevelSelect(),
					"unexpected selection level at index %d (%s)",
					j, expect[j].Decl,
				)
				i.Value()
				j++
			},
		)
		require.Zero(t, err.Error())
		require.False(t, err.IsErr())
		for _, e := range expect[j:] {
			t.Errorf(
				"missing {T: %s; V: %s}",
				e.Type, e.Value,
			)
		}
	})
}

func TestZeroValueToString(t *testing.T) {
	var expect gqlscan.Expect
	require.Zero(t, expect.String())

	var token gqlscan.Token
	require.Zero(t, token.String())
}

var testdataBlockStrings = []struct {
	decl         string
	input        string
	tokenIndex   int
	buffer       []byte
	expectWrites [][]byte
}{
	{decl: Decl(),
		input:        `{f(a:"0")}`,
		tokenIndex:   5,
		buffer:       nil,
		expectWrites: [][]byte{},
	},
	{decl: Decl(),
		input:      `{f(a:"0")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("0"),
		},
	},
	{decl: Decl(),
		input:      `{f(a:"01234567")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("01234567"),
		},
	},
	{decl: Decl(),
		input:      `{f(a:"0123456789ab")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("01234567"),
			[]byte("89ab"),
		},
	},
	{decl: Decl(),
		input:        `{f(a:"""""")}`,
		tokenIndex:   5,
		buffer:       make([]byte, 8),
		expectWrites: [][]byte{},
	},
	{decl: Decl(),
		input:      `{f(a:"""abc""")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("abc"),
		},
	},
	{decl: Decl(),
		input:      `{f(a:"""\n\t" """)}`,
		tokenIndex: 5,
		buffer:     make([]byte, 10),
		expectWrites: [][]byte{
			[]byte(`\n\t" `),
		},
	},
	{decl: Decl(),
		input: `{f(a:"""
					


					1234567
					12345678
					


				""")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("1234567\n"),
			[]byte("12345678"),
		},
	},
	{decl: Decl(),
		input: `{f(a:"""
					first line
					 second\t\tline
				 """)}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("first li"),
			[]byte("ne\n seco"),
			[]byte(`nd\t\tli`),
			[]byte(`ne`),
		},
	},
	{decl: Decl(),
		input:        `{f(a:"""\"""""")}`,
		tokenIndex:   5,
		buffer:       make([]byte, 3),
		expectWrites: [][]byte{[]byte(`"""`)},
	},
	{decl: Decl(),
		input: `{f(a:"""
					a
					 b
					"
					\
					\"""
				""")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 1),
		expectWrites: [][]byte{
			[]byte("a"), []byte("\n"),
			[]byte(" "), []byte("b"), []byte("\n"),
			[]byte(`"`), []byte("\n"),
			[]byte(`\`), []byte("\n"),
			[]byte(`"`), []byte(`"`), []byte(`"`),
		},
	},
	{decl: Decl(),
		// Three non-empty and two empty lines.
		// The second non-empty line consists of two spaces.
		input: `{f(a:"""` +
			"\n   a\n   \n     \n   \n   b\n" +
			`""")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("a\n\n  \n\nb"),
		},
	},
	{decl: Decl(),
		input:      blockstring_2747b,
		tokenIndex: 5,
		buffer:     make([]byte, 4096), // 4 KiB buffer
		expectWrites: [][]byte{
			[]byte(blockstring_2747b_expect),
		},
	},
}

func TestScanInterpreted(t *testing.T) {
	for _, td := range testdataBlockStrings {
		t.Run(td.decl, func(t *testing.T) {
			require := require.New(t)
			writes, c := [][]byte{}, 0
			err := gqlscan.Scan(
				[]byte(td.input),
				func(i *gqlscan.Iterator) (err bool) {
					if c != td.tokenIndex {
						c++
						return false
					}
					i.ScanInterpreted(td.buffer, func(b []byte) (stop bool) {
						w := make([]byte, len(b))
						copy(w, b)
						writes = append(writes, w)
						return false
					})
					return true
				},
			)
			require.Equal(
				gqlscan.ErrCallbackFn, err.Code,
				"unexpected error: %s", err.Error(),
			)
			require.Len(writes, len(td.expectWrites))
			for i, e := range td.expectWrites {
				require.Equal(
					string(e), string(writes[i]),
					"unexpected write at index %d (%s)", i, td.decl,
				)
			}
		})
	}
}

func TestScanInterpretedStop(t *testing.T) {
	const s = `
		first line\"""
		second\tline
	`
	const e = "first line\"\"\"\nsecond\\tline"
	const q = `{f(a:"""` + s + `""")}`

	itrFn := func(t *testing.T) func(i *gqlscan.Iterator) {
		c := -1
		require := require.New(t)
		return func(i *gqlscan.Iterator) {
			c++
			if c != 5 {
				return
			}
			for stopAt := 0; stopAt < len(e); stopAt++ {
				buf, callCount := make([]byte, 1), 0
				var r strings.Builder
				r.Grow(len(e))
				i.ScanInterpreted(buf, func(buffer []byte) (stop bool) {
					r.Write(buffer)
					callCount++
					return callCount > stopAt
				})
				require.Equal(e[:stopAt+1], r.String())
			}
		}
	}

	t.Run("Scan", func(t *testing.T) {
		require := require.New(t)
		fn, in := itrFn(t), []byte(q)
		err := gqlscan.Scan(in, func(i *gqlscan.Iterator) (err bool) {
			fn(i)
			return false
		})
		require.False(err.IsErr())
		require.Equal(q, string(in), "making sure input isn't mutated")
	})

	t.Run("ScanAll", func(t *testing.T) {
		require := require.New(t)
		fn, in := itrFn(t), []byte(q)
		err := gqlscan.ScanAll(in, func(i *gqlscan.Iterator) { fn(i) })
		require.False(err.IsErr())
		require.Equal(q, string(in), "making sure input isn't mutated")
	})
}

func Decl() string { return decl(2) }

func decl(skipFrames int) string {
	_, _, line, _ := runtime.Caller(skipFrames)
	return fmt.Sprintf("%d", line)
}

func Token(t gqlscan.Token, v ...string) Expect {
	var val string
	if len(v) > 1 {
		panic("expected single value")
	} else if len(v) > 0 {
		val = v[0]
	}
	return Expect{
		Decl:  decl(2),
		Type:  t,
		Value: val,
	}
}

func Input(input string, e ...Expect) TestInput {
	if len(e) < 1 {
		panic("requires at least one expectation")
	}
	return TestInput{
		decl:   decl(2),
		input:  input,
		expect: e,
	}
}

func InputErr(input string, e string) TestInputErr {
	if len(e) < 1 {
		panic("requires at least one expectation")
	}
	return TestInputErr{
		decl:      decl(2),
		input:     input,
		expectErr: e,
	}
}
