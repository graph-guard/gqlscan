package gqlscan_test

import (
	_ "embed"
	"fmt"
	"path/filepath"
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
	Input(`{foo}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "foo"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query {foo}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "foo"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: {foo: false})}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "foo"),
		Token(gqlscan.TokenFalse, "false"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: false)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenFalse, "false"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: true)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenTrue, "true"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: null)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: [])}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: [[]])}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: 0)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: 0.0)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenFloat, "0.0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: 42)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenInt, "42"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: -42)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenInt, "-42"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: -42.5678)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenFloat, "-42.5678"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(f: -42.5678e2)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenFloat, "-42.5678e2"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{ f (f: {x: 2}) }`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "f"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "x"),
		Token(gqlscan.TokenInt, "2"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text done }`,
		// Fragment f1
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f1"),
		Token(gqlscan.TokenFragTypeCond, "Query"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "todos"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f2"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		// Query Todos
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Todos"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f1"),
		Token(gqlscan.TokenSelEnd),

		// Fragment f2
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f2"),
		Token(gqlscan.TokenFragTypeCond, "Todo"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "id"),
		Token(gqlscan.TokenField, "text"),
		Token(gqlscan.TokenField, "done"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query Q(
		$variable: Foo,
		$ v: [ [ Bar ] ] = [[{f:0}] null [null]]
	) {
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
	} mutation M($variable: Foo={f:2}, $v: [ [ Bar ] ]) {
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
	}`,
		// Query Q
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Q"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "variable"),
		Token(gqlscan.TokenVarTypeName, "Foo"),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "Bar"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "f"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFieldAlias, "foo_alias"),
		Token(gqlscan.TokenField, "foo"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "x"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFieldAlias, "foobar_alias"),
		Token(gqlscan.TokenField, "foo_bar"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenField, "bar"),
		Token(gqlscan.TokenField, "baz"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragInline, "A"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz_taz_A"),
		Token(gqlscan.TokenFragRef, "namedFragment1"),
		Token(gqlscan.TokenFragRef, "namedFragment2"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenFragInline, "B"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz_taz_B"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz1"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "bool"),
		Token(gqlscan.TokenFalse, "false"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "bool"),
		Token(gqlscan.TokenTrue, "true"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz3"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "string"),
		Token(gqlscan.TokenStr, "okay"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz4"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "array"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz5"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "variable"),
		Token(gqlscan.TokenVarRef, "variable"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz6"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "variable"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz7"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "object"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "number0"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenObjField, "number1"),
		Token(gqlscan.TokenInt, "2"),
		Token(gqlscan.TokenObjField, "number2"),
		Token(gqlscan.TokenFloat, "123456789.1234e2"),

		Token(gqlscan.TokenObjField, "arr0"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "x"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArrEnd),

		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		// Mutation M
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenOprName, "M"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "variable"),
		Token(gqlscan.TokenVarTypeName, "Foo"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "f"),
		Token(gqlscan.TokenInt, "2"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "Bar"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "foo"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "x"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "foo_bar"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenField, "bar"),
		Token(gqlscan.TokenField, "baz"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragInline, "A"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz_taz_A"),
		Token(gqlscan.TokenFragRef, "namedFragment1"),
		Token(gqlscan.TokenFragRef, "namedFragment2"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenFragInline, "B"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "baz_fuzz_taz_B"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz1"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "bool"),
		Token(gqlscan.TokenFalse, "false"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz2"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "bool"),
		Token(gqlscan.TokenTrue, "true"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz3"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "string"),
		Token(gqlscan.TokenStr, "okay"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz4"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "array"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz5"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "variable"),
		Token(gqlscan.TokenVarRef, "variable"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz6"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "variable"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "baz_fuzz_taz7"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "object"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "number0"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenObjField, "number1"),
		Token(gqlscan.TokenInt, "2"),
		Token(gqlscan.TokenObjField, "number2"),
		Token(gqlscan.TokenFloat, "123456789.1234e2"),

		Token(gqlscan.TokenObjField, "arr0"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "x"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArrEnd),

		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		// Fragment f1
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f1"),
		Token(gqlscan.TokenFragTypeCond, "Query"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "todos"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f2"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),

		// Query Todos
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Todos"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f1"),
		Token(gqlscan.TokenSelEnd),

		// Fragment f2
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f2"),
		Token(gqlscan.TokenFragTypeCond, "Todo"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "id"),
		Token(gqlscan.TokenField, "text"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "foo"),
		Token(gqlscan.TokenInt, "2"),
		Token(gqlscan.TokenArg, "bar"),
		Token(gqlscan.TokenStr, "ok"),
		Token(gqlscan.TokenArg, "baz"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenField, "done"),
		Token(gqlscan.TokenSelEnd),

		// Subscription S
		Token(gqlscan.TokenDefSub),
		Token(gqlscan.TokenOprName, "S"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "Input"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "sub"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "i"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
	),

	// Comments
	Input("  #comment1\n  #comment2\n  {x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{  #comment1\n  #comment2\n  x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x  #comment1\n  #comment2\n  }",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x}  #comment1\n  #comment2\n",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x(  #comment1\n  #comment2\n  y:0)}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x(y  #comment1\n  #comment2\n  :0)}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x(y:  #comment1\n  #comment2\n  0)}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x(y:0  #comment1\n  #comment2\n  )}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{x(y:0)  #comment1\n  #comment2\n  }",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query  #comment1\n  #comment2\n  {x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("mutation  #comment1\n  #comment2\n  {x}",
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("fragment  #comment1\n  #comment2\n  f on X{x}",
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f"),
		Token(gqlscan.TokenFragTypeCond, "X"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("fragment f  #comment1\n  #comment2\n  on X{x}",
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f"),
		Token(gqlscan.TokenFragTypeCond, "X"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("fragment f on  #comment1\n  #comment2\n  X{x}",
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f"),
		Token(gqlscan.TokenFragTypeCond, "X"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("fragment f on X  #comment1\n  #comment2\n  {x}",
		Token(gqlscan.TokenDefFrag),
		Token(gqlscan.TokenFragName, "f"),
		Token(gqlscan.TokenFragTypeCond, "X"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{  ...  #comment1\n  #comment2\n  f  }",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{  ...  f  #comment1\n  #comment2\n  }",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFragRef, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query(  #comment1\n  #comment2\n  $x: T){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x  #comment1\n  #comment2\n  : T){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x:  #comment1\n  #comment2\n  T){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x:[  #comment1\n  #comment2\n  T]){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x:[T  #comment1\n  #comment2\n  ]){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x:[T]  #comment1\n  #comment2\n  ){x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("query($x:[T])  #comment1\n  #comment2\n  {x}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{f#comment\n{f2}}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f2"),
		Token(gqlscan.TokenSelEnd),
		Token(gqlscan.TokenSelEnd),
	),

	// String escape
	Input(`{x(s:"\"")}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "s"),
		Token(gqlscan.TokenStr, `\"`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{x(s:"\\")}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "s"),
		Token(gqlscan.TokenStr, `\\`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{x(s:"\\\"")}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "s"),
		Token(gqlscan.TokenStr, `\\\"`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),

	Input(`{x(y:1e8)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenFloat, `1e8`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{x(y:0e8)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenFloat, `0e8`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{x(y:0e+8)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenFloat, `0e+8`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{x(y:0e-8)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "y"),
		Token(gqlscan.TokenFloat, `0e-8`),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`mutation{x}`,
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`mutation($x:T){x}`,
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "x"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`mutation M{x}`,
		Token(gqlscan.TokenDefMut),
		Token(gqlscan.TokenOprName, "M"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(o:{o2:{x:[]}})}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "o"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "o2"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "x"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenObjEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(a:[0])}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenInt, "0"),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query($v:T ! ){x(a:$v)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query ($v: [ [ T ! ] ! ] ! ) {x(a:$v)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "T"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{ bob : alice }`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenFieldAlias, "bob"),
		Token(gqlscan.TokenField, "alice"),
		Token(gqlscan.TokenSelEnd),
	),
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
		$ # sample comment text line
		# sample comment text line
		v # sample comment text line
		# sample comment text line
		: # sample comment text line
		# sample comment text line
		Int # sample comment text line
		# sample comment text line
		! # sample comment text line
		# sample comment text line
		= # sample comment text line
		# sample comment text line
		42 # sample comment text line
		# sample comment text line
		@ # sample comment text line
		# sample comment text line
		d1 # sample comment text line
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
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenInt, "42"),
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
	Input(`{f}
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
		#01234567810abcd`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(a:
		"\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000"
		b:123456789
	)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenStr,
			`\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000`),
		Token(gqlscan.TokenArg, "b"),
		Token(gqlscan.TokenInt, "123456789"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input("{f(a:"+string_2695b+")}",
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenStr, string_2695b[1:len(string_2695b)-1]),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`{f(
		a:""""""
		b:"""abc"""
		c:"""\n\t" """
		d:"""
			foo
				bar
		"""
	)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenStrBlock),
		Token(gqlscan.TokenArg, "b"),
		Token(gqlscan.TokenStrBlock, "abc"),
		Token(gqlscan.TokenArg, "c"),
		Token(gqlscan.TokenStrBlock, `\n\t" `),
		Token(gqlscan.TokenArg, "d"),
		Token(gqlscan.TokenStrBlock,
			"\n\t\t\tfoo\n\t\t\t\tbar\n\t\t"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`subscription S{f}`,
		Token(gqlscan.TokenDefSub),
		Token(gqlscan.TokenOprName, "S"),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "f"),
		Token(gqlscan.TokenSelEnd),
	),
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
	Input(`query($v: Int = 12 @ok $v2: String) {x(a:$v)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenInt, "12"),
		Token(gqlscan.TokenDirName, "ok"),
		Token(gqlscan.TokenVarName, "v2"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query Int($v: Int = 12) {x(a:$v)}
		query String($v: String = "text") {x(a:$v)}
		query BooleanTrue($v: Boolean = true) {x(a:$v)}
		query BooleanFalse($v: Boolean = false) {x(a:$v)}
		query Array($v: [Int] = [12]) {x(a:$v)}
		`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Int"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenInt, "12"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "String"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenStr, "text"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "BooleanTrue"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "Boolean"),
		Token(gqlscan.TokenTrue, "true"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "BooleanFalse"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeName, "Boolean"),
		Token(gqlscan.TokenFalse, "false"),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),

		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenOprName, "Array"),
		Token(gqlscan.TokenVarList),
		Token(gqlscan.TokenVarName, "v"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenInt, "12"),
		Token(gqlscan.TokenArrEnd),
		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a"),
		Token(gqlscan.TokenVarRef, "v"),
		Token(gqlscan.TokenArgListEnd),
		Token(gqlscan.TokenSelEnd),
	),
	Input(`query(
		$v1: Boolean = false
		$v2: Boolean = true
		$v3: Int = 12
		$v4: Float = -3.14159265359
		$v5: String = "default value"
		$v6: String = ""
		$v7: Int = null
		$v8: [Int] = [1,null,3]
		$v9: [Int] = []
		$v10: Input = {foo: "bar"}
		$v11: Input = {faz: "baz" taz: """maz"""}
		$v12: String! = """block string"""
		$v13: String! = """"""
	) {x(
		a1:$v1
		a2:$v2
		a3:$v3
		a4:$v4
		a5:$v5
		a6:$v6
		a7:$v7
		a8:$v8
		a9:$v9
		a10:$v10
		a11:$v11
		a12:$v12
		a13:$v13
	)}`,
		Token(gqlscan.TokenDefQry),
		Token(gqlscan.TokenVarList),

		Token(gqlscan.TokenVarName, "v1"),
		Token(gqlscan.TokenVarTypeName, "Boolean"),
		Token(gqlscan.TokenFalse, "false"),

		Token(gqlscan.TokenVarName, "v2"),
		Token(gqlscan.TokenVarTypeName, "Boolean"),
		Token(gqlscan.TokenTrue, "true"),

		Token(gqlscan.TokenVarName, "v3"),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenInt, "12"),

		Token(gqlscan.TokenVarName, "v4"),
		Token(gqlscan.TokenVarTypeName, "Float"),
		Token(gqlscan.TokenFloat, "-3.14159265359"),

		Token(gqlscan.TokenVarName, "v5"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenStr, "default value"),

		Token(gqlscan.TokenVarName, "v6"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenStr, ""),

		Token(gqlscan.TokenVarName, "v7"),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenNull, "null"),

		Token(gqlscan.TokenVarName, "v8"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenInt, "1"),
		Token(gqlscan.TokenNull, "null"),
		Token(gqlscan.TokenInt, "3"),
		Token(gqlscan.TokenArrEnd),

		Token(gqlscan.TokenVarName, "v9"),
		Token(gqlscan.TokenVarTypeArr),
		Token(gqlscan.TokenVarTypeName, "Int"),
		Token(gqlscan.TokenVarTypeArrEnd),
		Token(gqlscan.TokenArr),
		Token(gqlscan.TokenArrEnd),

		Token(gqlscan.TokenVarName, "v10"),
		Token(gqlscan.TokenVarTypeName, "Input"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "foo"),
		Token(gqlscan.TokenStr, "bar"),
		Token(gqlscan.TokenObjEnd),

		Token(gqlscan.TokenVarName, "v11"),
		Token(gqlscan.TokenVarTypeName, "Input"),
		Token(gqlscan.TokenObj),
		Token(gqlscan.TokenObjField, "faz"),
		Token(gqlscan.TokenStr, "baz"),
		Token(gqlscan.TokenObjField, "taz"),
		Token(gqlscan.TokenStrBlock, "maz"),
		Token(gqlscan.TokenObjEnd),

		Token(gqlscan.TokenVarName, "v12"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenStrBlock, "block string"),

		Token(gqlscan.TokenVarName, "v13"),
		Token(gqlscan.TokenVarTypeName, "String"),
		Token(gqlscan.TokenVarTypeNotNull),
		Token(gqlscan.TokenStrBlock, ""),

		Token(gqlscan.TokenVarListEnd),
		Token(gqlscan.TokenSel),
		Token(gqlscan.TokenField, "x"),
		Token(gqlscan.TokenArgList),
		Token(gqlscan.TokenArg, "a1"),
		Token(gqlscan.TokenVarRef, "v1"),
		Token(gqlscan.TokenArg, "a2"),
		Token(gqlscan.TokenVarRef, "v2"),
		Token(gqlscan.TokenArg, "a3"),
		Token(gqlscan.TokenVarRef, "v3"),
		Token(gqlscan.TokenArg, "a4"),
		Token(gqlscan.TokenVarRef, "v4"),
		Token(gqlscan.TokenArg, "a5"),
		Token(gqlscan.TokenVarRef, "v5"),
		Token(gqlscan.TokenArg, "a6"),
		Token(gqlscan.TokenVarRef, "v6"),
		Token(gqlscan.TokenArg, "a7"),
		Token(gqlscan.TokenVarRef, "v7"),
		Token(gqlscan.TokenArg, "a8"),
		Token(gqlscan.TokenVarRef, "v8"),
		Token(gqlscan.TokenArg, "a9"),
		Token(gqlscan.TokenVarRef, "v9"),
		Token(gqlscan.TokenArg, "a10"),
		Token(gqlscan.TokenVarRef, "v10"),
		Token(gqlscan.TokenArg, "a11"),
		Token(gqlscan.TokenVarRef, "v11"),
		Token(gqlscan.TokenArg, "a12"),
		Token(gqlscan.TokenVarRef, "v12"),
		Token(gqlscan.TokenArg, "a13"),
		Token(gqlscan.TokenVarRef, "v13"),
		Token(gqlscan.TokenArgListEnd),
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
	InputErr( // Unexpected token as query.
		"q",
		"error at index 0 ('q'): unexpected token; expected definition",
	),
	InputErr( // Missing square bracket in type.
		"query($a: [A){f}",
		"error at index 11 ('A'): invalid type; "+
			"expected variable type",
	),
	InputErr( // Missing square bracket in type.
		"query($a: [[A]){f}",
		"error at index 13 (']'): invalid type; "+
			"expected variable type",
	),
	InputErr( // Unexpected square bracket in variable type.
		"query($a: A]){f}",
		"error at index 11 (']'): unexpected token; "+
			"expected variable",
	),
	InputErr( // Unexpected square bracket in variable type.
		"query($a: [[A]]]){f}",
		"error at index 15 (']'): unexpected token; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Missing query closing curly bracket.
		"{",
		"error at index 1: unexpected end of file; expected selection",
	),
	InputErr( // Invalid field name.
		"{1abc}",
		"error at index 1 ('1'): unexpected token; "+
			"expected field name or alias",
	),
	InputErr( // Trailing closing curly bracket.
		"{f}}",
		"error at index 3 ('}'): unexpected token; expected definition",
	),
	InputErr( // Query missing closing curly bracket.
		"{}",
		"error at index 1 ('}'): unexpected token; "+
			"expected field name or alias",
	),
	InputErr( // Empty args.
		"{f()}",
		"error at index 3 (')'): unexpected token; expected argument name",
	),
	InputErr( // Argument missing column.
		"{f(x null))}",
		"error at index 5 ('n'): "+
			"unexpected token; expected column after argument name",
	),
	InputErr( // Argument with trailing closing parenthesis.
		"{f(x:null))}",
		"error at index 10 (')'): unexpected token; "+
			"expected field name or alias",
	),
	InputErr( // Argument missing closing parenthesis.
		"{f(}",
		"error at index 3 ('}'): unexpected token; expected argument name",
	),
	InputErr( // Invalid argument.
		"{f(x:abc))}",
		"error at index 5 ('a'): invalid value; expected value",
	),
	InputErr( // String argument missing closing quotes.
		`{f(x:"))}`,
		"error at index 9: unexpected end of file; expected end of string",
	),
	InputErr( // Invalid negative number.
		`{f(x:-))}`,
		"error at index 6 (')'): invalid number value; expected value",
	),
	InputErr( // Number missing fraction.
		`{f(x:1.))}`,
		"error at index 7 (')'): invalid number value; expected value",
	),
	InputErr( // Number missing exponent.
		`{f(x:1.2e))}`,
		"error at index 9 (')'): invalid number value; expected value",
	),
	InputErr( // Number with leading zero.
		`{f(x:0123))}`,
		"error at index 6 ('1'): invalid number value; expected value",
	),

	// --- Unexpected EOF ---
	InputErr( // Unexpected EOF.
		"",
		"error at index 0: unexpected end of file; expected definition",
	),
	InputErr( // Unexpected EOF.
		"query",
		"error at index 5: unexpected end of file; "+
			"expected variable list or selection set",
	),
	InputErr( // Unexpected EOF.
		"query Name",
		"error at index 10: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"query Name ",
		"error at index 11: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"mutation Name",
		"error at index 13: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"mutation Name ",
		"error at index 14: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"subscription Name",
		"error at index 17: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"subscription Name ",
		"error at index 18: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"query(",
		"error at index 6: unexpected end of file; "+
			"expected variable",
	),
	InputErr( // Unexpected EOF.
		"query( ",
		"error at index 7: unexpected end of file; "+
			"expected variable",
	),
	InputErr( // Unexpected EOF.
		"query($",
		"error at index 7: unexpected end of file; "+
			"expected variable name",
	),
	InputErr( // Variable missing name.
		"query($ ",
		"error at index 8: unexpected end of file; "+
			"expected variable name",
	),
	InputErr( // Unexpected EOF.
		"query($v",
		"error at index 8: unexpected end of file; "+
			"expected column after variable name",
	),
	InputErr( // Unexpected EOF.
		"query($v ",
		"error at index 9: unexpected end of file; "+
			"expected column after variable name",
	),
	InputErr( // Unexpected EOF.
		"query($v:",
		"error at index 9: unexpected end of file; "+
			"expected variable type",
	),
	InputErr( // Unexpected EOF.
		"query($v: ",
		"error at index 10: unexpected end of file; "+
			"expected variable type",
	),
	InputErr( // Unexpected EOF.
		"query($v: T",
		"error at index 11: unexpected end of file; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected EOF.
		"query($v: T ",
		"error at index 12: unexpected end of file; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected EOF.
		"query($v: T)",
		"error at index 12: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"query($v: T) ",
		"error at index 13: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"{",
		"error at index 1: unexpected end of file; "+
			"expected selection",
	),
	InputErr( // Unexpected EOF.
		"{ ",
		"error at index 2: unexpected end of file; "+
			"expected selection",
	),
	InputErr( // Unexpected EOF.
		"{foo",
		"error at index 4: unexpected end of file; "+
			"expected field name or alias",
	),
	InputErr( // Unexpected EOF.
		"{foo ",
		"error at index 5: unexpected end of file; "+
			"expected field name or alias",
	),
	InputErr( // Unexpected EOF.
		"{foo(",
		"error at index 5: unexpected end of file; "+
			"expected argument name",
	),
	InputErr( // Unexpected EOF.
		"{foo( ",
		"error at index 6: unexpected end of file; "+
			"expected argument name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name",
		"error at index 9: unexpected end of file; "+
			"expected column after argument name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name ",
		"error at index 10: unexpected end of file; "+
			"expected column after argument name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name:",
		"error at index 10: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: ",
		"error at index 11: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: {",
		"error at index 12: unexpected end of file; "+
			"expected object field name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: { ",
		"error at index 13: unexpected end of file; "+
			"expected object field name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: {field",
		"error at index 17: unexpected end of file; "+
			"expected column after object field name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: {field ",
		"error at index 18: unexpected end of file; "+
			"expected column after object field name",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: {field:",
		"error at index 18: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		"{foo(name: {field: ",
		"error at index 19: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: "`,
		"error at index 12: unexpected end of file; "+
			"expected end of string",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: ""`,
		"error at index 13: unexpected end of file; "+
			"expected argument list closure or argument",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: f`,
		"error at index 12: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: t`,
		"error at index 12: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: n`,
		"error at index 12: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: 0`,
		"error at index 12: unexpected end of file; "+
			"expected argument list closure or argument",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: 0 `,
		"error at index 13: unexpected end of file; "+
			"expected argument list closure or argument",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: -`,
		"error at index 12: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: 0.`,
		"error at index 13: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{foo(name: 0.1e`,
		"error at index 15: unexpected end of file; expected value",
	),
	InputErr( // Unexpected EOF.
		`{.`,
		"error at index 2: unexpected end of file; expected fragment",
	),
	InputErr( // Unexpected EOF.
		`{..`,
		"error at index 3: unexpected end of file; expected fragment",
	),
	InputErr( // Unexpected EOF.
		`{...`,
		"error at index 4: unexpected end of file; expected fragment",
	),
	InputErr( // Unexpected EOF.
		`{... `,
		"error at index 5: unexpected end of file; expected fragment",
	),
	InputErr( // Unexpected EOF.
		`{... on`,
		"error at index 7: unexpected end of file; expected fragment",
	),
	InputErr( // Unexpected EOF.
		`{... on `,
		"error at index 8: unexpected end of file; "+
			"expected inlined fragment",
	),
	InputErr( // Unexpected EOF.
		"fragment",
		"error at index 8: unexpected end of file; "+
			"expected fragment name",
	),
	InputErr( // Unexpected EOF.
		"{x",
		"error at index 2: unexpected end of file; "+
			"expected field name or alias",
	),

	InputErr( // Invalid value.
		"{x(p:falsa",
		"error at index 5 ('f'): invalid value; "+
			"expected value",
	),
	InputErr( // Invalid value.
		"{x(p:truu",
		"error at index 5 ('t'): invalid value; "+
			"expected value",
	),
	InputErr( // Invalid value.
		"{x(p:nuli",
		"error at index 5 ('n'): invalid value; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		"{x(p:[",
		"error at index 6: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"query($x:T)x",
		"error at index 11 ('x'): unexpected token; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"mutation M",
		"error at index 10: unexpected end of file; "+
			"expected selection set",
	),
	InputErr( // Unexpected token.
		"query\x00",
		"error at index 5 (0x0): unexpected token; "+
			"expected operation name",
	),
	InputErr( // Unexpected token.
		"{x(y:12e)}",
		"error at index 8 (')'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:12.)}",
		"error at index 8 (')'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:12x)}",
		"error at index 7 ('x'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:12.12x)}",
		"error at index 10 ('x'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected EOF.
		"{x(y:12.12",
		"error at index 10: unexpected end of file; "+
			"expected argument list closure or argument",
	),
	InputErr( // Unexpected EOF.
		"{x(y:12.",
		"error at index 8: unexpected end of file; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:12e111x",
		"error at index 11 ('x'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:12ex",
		"error at index 8 ('x'): invalid number value; "+
			"expected value",
	),
	InputErr( // Unexpected token.
		"{x(y:{f})}",
		"error at index 7 ('}'): unexpected token; "+
			"expected column after object field name",
	),
	InputErr( // Unexpected token.
		"{x(\x00:1)}",
		"error at index 3 (0x0): unexpected token; "+
			"expected argument name",
	),
	InputErr( // Unexpected EOF.
		"{x(y\x00:1)}",
		"error at index 4 (0x0): unexpected token; "+
			"expected argument name",
	),
	InputErr( // Unexpected token.
		"query M [",
		"error at index 8 ('['): unexpected token; "+
			"expected selection set",
	),
	InputErr( // Unexpected token.
		"mutation M|",
		"error at index 10 ('|'): unexpected token; "+
			"expected selection set",
	),
	InputErr( // Unexpected EOF.
		"fragment f on",
		"error at index 13: unexpected end of file; "+
			"expected fragment type condition",
	),
	InputErr( // Unexpected token.
		"mutation\x00",
		"error at index 8 (0x0): unexpected token; "+
			"expected operation name",
	),
	InputErr( // Unexpected token.
		"subscription\x00",
		"error at index 12 (0x0): unexpected token; "+
			"expected operation name",
	),
	InputErr( // Unexpected token.
		"fragment\x00",
		"error at index 8 (0x0): unexpected token; "+
			"expected fragment name",
	),
	InputErr( // Unexpected EOF.
		"{x(y:$",
		"error at index 6: unexpected end of file; "+
			"expected referenced variable name",
	),
	InputErr( // Unexpected EOF.
		"mutation",
		"error at index 8: unexpected end of file; "+
			"expected variable list or selection set",
	),
	InputErr( // Unexpected EOF.
		"{x(y:null)",
		"error at index 10: unexpected end of file; "+
			"expected selection set or selection",
	),
	InputErr( // Unexpected token.
		"query($v |",
		"error at index 9 ('|'): unexpected token; "+
			"expected column after variable name",
	),
	InputErr( // Unexpected token.
		"query($v:[T] |)",
		"error at index 13 ('|'): unexpected token; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected token.
		"fragment X at",
		"error at index 11 ('a'): unexpected token; "+
			"expected keyword 'on'",
	),
	InputErr( // Unexpected EOF.
		"query($a:[A]",
		"error at index 12: unexpected end of file; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected EOF.
		"fragment f ",
		"error at index 11: unexpected end of file; "+
			"expected keyword 'on'",
	),
	InputErr( // Unexpected EOF.
		"{f{x} ",
		"error at index 6: unexpected end of file; "+
			"expected selection or end of selection set",
	),
	InputErr( // Unexpected token.
		"{f(x:\"abc\n\")}",
		"error at index 9 (0xa): unexpected token; "+
			"expected end of string",
	),
	InputErr( // Unexpected token.
		"{.f}",
		"error at index 2 ('f'): unexpected token; "+
			"expected fragment",
	),
	InputErr( // Unexpected token.
		"{..f}",
		"error at index 3 ('f'): unexpected token; "+
			"expected fragment",
	),
	InputErr( // Unexpected token.
		"query($v:T ! !){x(a:$v)}",
		"error at index 13 ('!'): unexpected token; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected token.
		"query($v: [ T ! ] ! ! ){x(a:$v)}",
		"error at index 20 ('!'): unexpected token; "+
			"expected variable list closure or variable name",
	),
	InputErr( // Unexpected token.
		"{alias : alias2 : x}",
		"error at index 16 (':'): unexpected token; "+
			"expected field name or alias",
	),
	InputErr( // Unexpected EOF.
		"{f:",
		"error at index 3: unexpected end of file; "+
			"expected field name",
	),
	InputErr( // Unexpected EOF.
		"{f: ",
		"error at index 4: unexpected end of file; "+
			"expected field name",
	),
	InputErr( // Invalid escape sequence.
		`{f(a:"\a")}`,
		"error at index 7 ('a'): unexpected token; "+
			"expected escaped sequence",
	),
	InputErr( // Invalid escape sequence.
		`{f(a:"\u")}`,
		"error at index 8 ('\"'): unexpected token; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Invalid escape sequence.
		`{f(a:"\u1")}`,
		"error at index 9 ('\"'): unexpected token; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Invalid escape sequence.
		`{f(a:"\u12")}`,
		"error at index 10 ('\"'): unexpected token; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Unexpected EOF.
		`{f(a:"\u`,
		"error at index 8: unexpected end of file; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Unexpected EOF.
		`{f(a:"\u1`,
		"error at index 9: unexpected end of file; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Unexpected EOF.
		`{f(a:"\u12`,
		"error at index 10: unexpected end of file; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Unexpected EOF.
		`{f(a:"\u123`,
		"error at index 11: unexpected end of file; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Invalid escape sequence.
		`{f(a:"\u123")}`,
		"error at index 11 ('\"'): unexpected token; "+
			"expected escaped unicode sequence",
	),
	InputErr( // Unexpected EOF.
		`{f(a:"""`,
		`error at index 8: unexpected end of file; `+
			"expected end of block string",
	),
	InputErr( // Unexpected EOF.
		`{f(a:""" `,
		"error at index 9: unexpected end of file; "+
			"expected end of block string",
	),
	InputErr( // Control character in string.
		`{f(a:"0123456`+string(rune(0x00))+`")}`,
		"error at index 13 (0x0): unexpected token; "+
			"expected end of string",
	),
	InputErr( // Unexpected EOF.
		`{f #c`,
		"error at index 5: unexpected end of file; "+
			"expected selection, selection set or end of selection set",
	),
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

	type ExpectLevel struct {
		Decl  string
		Type  gqlscan.Token
		Value string
		Level int
	}

	TokenLevel := func(
		t gqlscan.Token,
		level int,
		v ...string,
	) ExpectLevel {
		var val string
		if len(v) > 1 {
			panic("expected single value")
		} else if len(v) > 0 {
			val = v[0]
		}
		return ExpectLevel{
			Decl:  decl(2),
			Level: level,
			Type:  t,
			Value: val,
		}
	}

	expect := []ExpectLevel{
		// Query Q
		TokenLevel(gqlscan.TokenDefQry, 0, ""),
		TokenLevel(gqlscan.TokenOprName, 0, "Q"),
		TokenLevel(gqlscan.TokenVarList, 0, ""),
		TokenLevel(gqlscan.TokenVarName, 0, "variable"),
		TokenLevel(gqlscan.TokenVarTypeName, 0, "Foo"),
		TokenLevel(gqlscan.TokenVarName, 0, "v"),
		TokenLevel(gqlscan.TokenVarTypeArr, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeArr, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeName, 0, "Bar"),
		TokenLevel(gqlscan.TokenVarTypeArrEnd, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeArrEnd, 0, ""),
		TokenLevel(gqlscan.TokenVarListEnd, 0, ""),
		TokenLevel(gqlscan.TokenSel, 0, ""),
		TokenLevel(gqlscan.TokenFieldAlias, 1, "foo_alias"),
		TokenLevel(gqlscan.TokenField, 1, "foo"),
		TokenLevel(gqlscan.TokenArgList, 1, ""),
		TokenLevel(gqlscan.TokenArg, 1, "x"),
		TokenLevel(gqlscan.TokenNull, 1, "null"),
		TokenLevel(gqlscan.TokenArgListEnd, 1, ""),

		TokenLevel(gqlscan.TokenSel, 1, ""),
		TokenLevel(gqlscan.TokenFieldAlias, 2, "foobar_alias"),
		TokenLevel(gqlscan.TokenField, 2, "foo_bar"),
		TokenLevel(gqlscan.TokenSelEnd, 2, ""),
		TokenLevel(gqlscan.TokenField, 1, "bar"),
		TokenLevel(gqlscan.TokenField, 1, "baz"),
		TokenLevel(gqlscan.TokenSel, 1, ""),
		TokenLevel(gqlscan.TokenField, 2, "baz_fuzz"),
		TokenLevel(gqlscan.TokenSel, 2, ""),
		TokenLevel(gqlscan.TokenFragInline, 3, "A"),
		TokenLevel(gqlscan.TokenSel, 3, ""),
		TokenLevel(gqlscan.TokenField, 4, "baz_fuzz_taz_A"),
		TokenLevel(gqlscan.TokenFragRef, 4, "namedFragment1"),
		TokenLevel(gqlscan.TokenFragRef, 4, "namedFragment2"),
		TokenLevel(gqlscan.TokenSelEnd, 4, ""),
		TokenLevel(gqlscan.TokenFragInline, 3, "B"),
		TokenLevel(gqlscan.TokenSel, 3, ""),
		TokenLevel(gqlscan.TokenField, 4, "baz_fuzz_taz_B"),
		TokenLevel(gqlscan.TokenSelEnd, 4, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz1"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "bool"),
		TokenLevel(gqlscan.TokenFalse, 3, "false"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz2"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "bool"),
		TokenLevel(gqlscan.TokenTrue, 3, "true"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz3"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "string"),
		TokenLevel(gqlscan.TokenStr, 3, "okay"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz4"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "array"),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz5"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "variable"),
		TokenLevel(gqlscan.TokenVarRef, 3, "variable"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz6"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "variable"),
		TokenLevel(gqlscan.TokenVarRef, 3, "v"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz7"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "object"),
		TokenLevel(gqlscan.TokenObj, 3, ""),
		TokenLevel(gqlscan.TokenObjField, 3, "number0"),
		TokenLevel(gqlscan.TokenInt, 3, "0"),
		TokenLevel(gqlscan.TokenObjField, 3, "number1"),
		TokenLevel(gqlscan.TokenInt, 3, "2"),
		TokenLevel(gqlscan.TokenObjField, 3, "number2"),
		TokenLevel(gqlscan.TokenFloat, 3, "123456789.1234e2"),

		TokenLevel(gqlscan.TokenObjField, 3, "arr0"),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenObj, 3, ""),
		TokenLevel(gqlscan.TokenObjField, 3, "x"),
		TokenLevel(gqlscan.TokenNull, 3, "null"),
		TokenLevel(gqlscan.TokenObjEnd, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),

		TokenLevel(gqlscan.TokenObjEnd, 3, ""),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenSelEnd, 3, ""),
		TokenLevel(gqlscan.TokenSelEnd, 2, ""),
		TokenLevel(gqlscan.TokenSelEnd, 1, ""),

		// Mutation M
		TokenLevel(gqlscan.TokenDefMut, 0, ""),
		TokenLevel(gqlscan.TokenOprName, 0, "M"),
		TokenLevel(gqlscan.TokenVarList, 0, ""),
		TokenLevel(gqlscan.TokenVarName, 0, "variable"),
		TokenLevel(gqlscan.TokenVarTypeName, 0, "Foo"),
		TokenLevel(gqlscan.TokenVarName, 0, "v"),
		TokenLevel(gqlscan.TokenVarTypeArr, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeArr, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeName, 0, "Bar"),
		TokenLevel(gqlscan.TokenVarTypeArrEnd, 0, ""),
		TokenLevel(gqlscan.TokenVarTypeArrEnd, 0, ""),
		TokenLevel(gqlscan.TokenVarListEnd, 0, ""),
		TokenLevel(gqlscan.TokenSel, 0, ""),
		TokenLevel(gqlscan.TokenField, 1, "foo"),
		TokenLevel(gqlscan.TokenArgList, 1, ""),
		TokenLevel(gqlscan.TokenArg, 1, "x"),
		TokenLevel(gqlscan.TokenNull, 1, "null"),
		TokenLevel(gqlscan.TokenArgListEnd, 1, ""),
		TokenLevel(gqlscan.TokenSel, 1, ""),
		TokenLevel(gqlscan.TokenField, 2, "foo_bar"),
		TokenLevel(gqlscan.TokenSelEnd, 2, ""),
		TokenLevel(gqlscan.TokenField, 1, "bar"),
		TokenLevel(gqlscan.TokenField, 1, "baz"),
		TokenLevel(gqlscan.TokenSel, 1, ""),
		TokenLevel(gqlscan.TokenField, 2, "baz_fuzz"),
		TokenLevel(gqlscan.TokenSel, 2, ""),
		TokenLevel(gqlscan.TokenFragInline, 3, "A"),
		TokenLevel(gqlscan.TokenSel, 3, ""),
		TokenLevel(gqlscan.TokenField, 4, "baz_fuzz_taz_A"),
		TokenLevel(gqlscan.TokenFragRef, 4, "namedFragment1"),
		TokenLevel(gqlscan.TokenFragRef, 4, "namedFragment2"),
		TokenLevel(gqlscan.TokenSelEnd, 4, ""),
		TokenLevel(gqlscan.TokenFragInline, 3, "B"),
		TokenLevel(gqlscan.TokenSel, 3, ""),
		TokenLevel(gqlscan.TokenField, 4, "baz_fuzz_taz_B"),
		TokenLevel(gqlscan.TokenSelEnd, 4, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz1"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "bool"),
		TokenLevel(gqlscan.TokenFalse, 3, "false"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz2"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "bool"),
		TokenLevel(gqlscan.TokenTrue, 3, "true"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz3"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "string"),
		TokenLevel(gqlscan.TokenStr, 3, "okay"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz4"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "array"),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz5"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "variable"),
		TokenLevel(gqlscan.TokenVarRef, 3, "variable"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz6"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "variable"),
		TokenLevel(gqlscan.TokenVarRef, 3, "v"),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenField, 3, "baz_fuzz_taz7"),
		TokenLevel(gqlscan.TokenArgList, 3, ""),
		TokenLevel(gqlscan.TokenArg, 3, "object"),
		TokenLevel(gqlscan.TokenObj, 3, ""),
		TokenLevel(gqlscan.TokenObjField, 3, "number0"),
		TokenLevel(gqlscan.TokenInt, 3, "0"),
		TokenLevel(gqlscan.TokenObjField, 3, "number1"),
		TokenLevel(gqlscan.TokenInt, 3, "2"),
		TokenLevel(gqlscan.TokenObjField, 3, "number2"),
		TokenLevel(gqlscan.TokenFloat, 3, "123456789.1234e2"),

		TokenLevel(gqlscan.TokenObjField, 3, "arr0"),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArr, 3, ""),
		TokenLevel(gqlscan.TokenObj, 3, ""),
		TokenLevel(gqlscan.TokenObjField, 3, "x"),
		TokenLevel(gqlscan.TokenNull, 3, "null"),
		TokenLevel(gqlscan.TokenObjEnd, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),
		TokenLevel(gqlscan.TokenArrEnd, 3, ""),

		TokenLevel(gqlscan.TokenObjEnd, 3, ""),
		TokenLevel(gqlscan.TokenArgListEnd, 3, ""),
		TokenLevel(gqlscan.TokenSelEnd, 3, ""),
		TokenLevel(gqlscan.TokenSelEnd, 2, ""),
		TokenLevel(gqlscan.TokenSelEnd, 1, ""),

		// Fragment f1
		TokenLevel(gqlscan.TokenDefFrag, 0, ""),
		TokenLevel(gqlscan.TokenFragName, 0, "f1"),
		TokenLevel(gqlscan.TokenFragTypeCond, 0, "Query"),
		TokenLevel(gqlscan.TokenSel, 0, ""),
		TokenLevel(gqlscan.TokenField, 1, "todos"),
		TokenLevel(gqlscan.TokenSel, 1, ""),
		TokenLevel(gqlscan.TokenFragRef, 2, "f2"),
		TokenLevel(gqlscan.TokenSelEnd, 2, ""),
		TokenLevel(gqlscan.TokenSelEnd, 1, ""),

		// Query Todos
		TokenLevel(gqlscan.TokenDefQry, 0, ""),
		TokenLevel(gqlscan.TokenOprName, 0, "Todos"),
		TokenLevel(gqlscan.TokenSel, 0, ""),
		TokenLevel(gqlscan.TokenFragRef, 1, "f1"),
		TokenLevel(gqlscan.TokenSelEnd, 1, ""),

		// Fragment f2
		TokenLevel(gqlscan.TokenDefFrag, 0, ""),
		TokenLevel(gqlscan.TokenFragName, 0, "f2"),
		TokenLevel(gqlscan.TokenFragTypeCond, 0, "Todo"),
		TokenLevel(gqlscan.TokenSel, 0, ""),
		TokenLevel(gqlscan.TokenField, 1, "id"),
		TokenLevel(gqlscan.TokenField, 1, "text"),
		TokenLevel(gqlscan.TokenArgList, 1, ""),
		TokenLevel(gqlscan.TokenArg, 1, "foo"),
		TokenLevel(gqlscan.TokenInt, 1, "2"),
		TokenLevel(gqlscan.TokenArg, 1, "bar"),
		TokenLevel(gqlscan.TokenStr, 1, "ok"),
		TokenLevel(gqlscan.TokenArg, 1, "baz"),
		TokenLevel(gqlscan.TokenNull, 1, "null"),
		TokenLevel(gqlscan.TokenArgListEnd, 1, ""),
		TokenLevel(gqlscan.TokenField, 1, "done"),
		TokenLevel(gqlscan.TokenSelEnd, 1, ""),
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

type ExpectBlockStr struct {
	Decl         string
	Input        string
	TokenIndex   int
	Buffer       []byte
	ExpectWrites [][]byte
}

var testdataBlockStrings = []ExpectBlockStr{
	TokenBlockStr(
		`{f(a:"0")}`,
		nil,
		// No writes
	),
	TokenBlockStr(
		`{f(a:"0")}`,
		make([]byte, 8),
		"0",
	),
	TokenBlockStr(
		`{f(a:"01234567")}`,
		make([]byte, 8),
		"01234567",
	),
	TokenBlockStr(
		`{f(a:"0123456789ab")}`,
		make([]byte, 8),
		"01234567", "89ab",
	),
	TokenBlockStr(
		`{f(a:"""""")}`,
		make([]byte, 8),
		// No writes
	),
	TokenBlockStr(
		`{f(a:"""abc""")}`,
		make([]byte, 8),
		"abc",
	),
	TokenBlockStr(
		`{f(a:"""\n\t" """)}`,
		make([]byte, 10),
		`\n\t" `,
	),
	TokenBlockStr(
		`{f(a:"""
			


			1234567
			12345678
			


		""")}`,
		make([]byte, 8),
		"1234567\n", "12345678",
	),
	TokenBlockStr(
		`{f(a:"""
			first line
			 second\t\tline
		 """)}`,
		make([]byte, 8),
		"first li", "ne\n seco", `nd\t\tli`, `ne`,
	),
	TokenBlockStr(
		`{f(a:"""\"""""")}`,
		make([]byte, 3),
		`"""`,
	),
	TokenBlockStr(
		`{f(a:"""
			a
			 b
			"
			\
			\"""
		""")}`,
		make([]byte, 1),
		"a", "\n",
		" ", "b", "\n",
		`"`, "\n",
		`\`, "\n",
		`"`, `"`, `"`,
	),
	TokenBlockStr(
		// Three non-empty and two empty lines.
		// The second non-empty line consists of two spaces.
		`{f(a:"""`+
			"\n   a\n   \n     \n   \n   b\n"+
			`""")}`,
		make([]byte, 8),
		"a\n\n  \n\nb",
	),
	TokenBlockStr(
		blockstring_2747b,
		make([]byte, 4096), // 4 KiB buffer
		blockstring_2747b_expect,
	),
}

func TestScanInterpreted(t *testing.T) {
	for _, td := range testdataBlockStrings {
		t.Run(td.Decl, func(t *testing.T) {
			require := require.New(t)
			writes, c := [][]byte{}, 0
			err := gqlscan.Scan(
				[]byte(td.Input),
				func(i *gqlscan.Iterator) (err bool) {
					if c != td.TokenIndex {
						c++
						return false
					}
					i.ScanInterpreted(td.Buffer, func(b []byte) (stop bool) {
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
			require.Len(writes, len(td.ExpectWrites))
			for i, e := range td.ExpectWrites {
				require.Equal(
					string(e), string(writes[i]),
					"unexpected write at index %d (%s)", i, td.Decl,
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

func decl(skipFrames int) string {
	_, filename, line, _ := runtime.Caller(skipFrames)
	return fmt.Sprintf("%s:%d", filepath.Base(filename), line)
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

func TokenBlockStr(
	input string,
	buffer []byte,
	expectWrites ...string,
) ExpectBlockStr {
	e := make([][]byte, len(expectWrites))
	for i := range expectWrites {
		e[i] = []byte(expectWrites[i])
	}
	return ExpectBlockStr{
		TokenIndex:   5,
		Decl:         decl(2),
		Input:        input,
		Buffer:       buffer,
		ExpectWrites: e,
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
