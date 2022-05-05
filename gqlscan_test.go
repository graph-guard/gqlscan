package gqlscan_test

import (
	_ "embed"
	"testing"

	"github.com/graph-guard/gqlscan"

	"github.com/stretchr/testify/require"
)

type Expect struct {
	Index int
	Type  gqlscan.Token
	Value string
}

var testdata = []struct {
	index  int // Helps visually finding the dataset
	input  string
	expect []Expect
}{
	{0, `{foo}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "foo"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{1, `query {foo}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "foo"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{2, `{f(f: {foo: false})}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenObj, ""},
		{5, gqlscan.TokenObjField, "foo"},
		{6, gqlscan.TokenFalse, "false"},
		{7, gqlscan.TokenObjEnd, ""},
		{8, gqlscan.TokenSelEnd, ""},
	}},
	{3, `{f(f: false)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenFalse, "false"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{4, `{f(f: true)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenTrue, "true"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{5, `{f(f: null)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNull, "null"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{6, `{f(f: [])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenArr, ""},
		{5, gqlscan.TokenArrEnd, ""},
		{6, gqlscan.TokenSelEnd, ""},
	}},
	{7, `{f(f: [[]])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenArr, ""},
		{5, gqlscan.TokenArr, ""},
		{6, gqlscan.TokenArrEnd, ""},
		{7, gqlscan.TokenArrEnd, ""},
		{8, gqlscan.TokenSelEnd, ""},
	}},
	{8, `{f(f: [null "" [[] 42 false] true])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenArr, ""},
		{5, gqlscan.TokenNull, "null"},
		{6, gqlscan.TokenStr, ""},
		{7, gqlscan.TokenArr, ""},
		{8, gqlscan.TokenArr, ""},
		{9, gqlscan.TokenArrEnd, ""},
		{10, gqlscan.TokenNum, "42"},
		{11, gqlscan.TokenFalse, "false"},
		{12, gqlscan.TokenArrEnd, ""},
		{13, gqlscan.TokenTrue, "true"},
		{14, gqlscan.TokenArrEnd, ""},
		{15, gqlscan.TokenSelEnd, ""},
	}},
	{9, `{f(f: 0)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{10, `{f(f: 0.0)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "0.0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{11, `{f(f: 42)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "42"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{12, `{f(f: -42)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "-42"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{13, `{f(f: -42.5678)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "-42.5678"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{14, `{f(f: -42.5678e2)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenNum, "-42.5678e2"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{15, `{ f (f: {x: 2}) }`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "f"},
		{4, gqlscan.TokenObj, ""},
		{5, gqlscan.TokenObjField, "x"},
		{6, gqlscan.TokenNum, "2"},
		{7, gqlscan.TokenObjEnd, ""},
		{8, gqlscan.TokenSelEnd, ""},
	}},
	{16, `fragment f1 on Query { todos { ...f2 } }
	query Todos { ...f1 }
	fragment f2 on Todo { id text done }`, []Expect{
		// Fragment f1
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f1"},
		{2, gqlscan.TokenFragTypeCond, "Query"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "todos"},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenFragRef, "f2"},
		{7, gqlscan.TokenSelEnd, ""},
		{8, gqlscan.TokenSelEnd, ""},

		// Query Todos
		{9, gqlscan.TokenDefQry, ""},
		{10, gqlscan.TokenQryName, "Todos"},
		{11, gqlscan.TokenSel, ""},
		{12, gqlscan.TokenFragRef, "f1"},
		{13, gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{14, gqlscan.TokenDefFrag, ""},
		{15, gqlscan.TokenFragName, "f2"},
		{16, gqlscan.TokenFragTypeCond, "Todo"},
		{17, gqlscan.TokenSel, ""},
		{18, gqlscan.TokenField, "id"},
		{19, gqlscan.TokenField, "text"},
		{20, gqlscan.TokenField, "done"},
		{21, gqlscan.TokenSelEnd, ""},
	}},
	{17, `query Q($variable: Foo, $v: [ [ Bar ] ]) {
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
	) done }`, []Expect{
		// Query Q
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenQryName, "Q"},
		{2, gqlscan.TokenVarName, "variable"},
		{3, gqlscan.TokenVarTypeName, "Foo"},
		{4, gqlscan.TokenVarName, "v"},
		{5, gqlscan.TokenVarTypeArr, ""},
		{6, gqlscan.TokenVarTypeArr, ""},
		{7, gqlscan.TokenVarTypeName, "Bar"},
		{8, gqlscan.TokenVarTypeArrEnd, ""},
		{9, gqlscan.TokenVarTypeArrEnd, ""},
		{10, gqlscan.TokenSel, ""},
		{11, gqlscan.TokenField, "foo"},
		{12, gqlscan.TokenArg, "x"},
		{13, gqlscan.TokenNull, "null"},
		{14, gqlscan.TokenSel, ""},
		{15, gqlscan.TokenField, "foo_bar"},
		{16, gqlscan.TokenSelEnd, ""},
		{17, gqlscan.TokenField, "bar"},
		{18, gqlscan.TokenField, "baz"},
		{19, gqlscan.TokenSel, ""},
		{20, gqlscan.TokenField, "baz_fuzz"},
		{21, gqlscan.TokenSel, ""},
		{22, gqlscan.TokenFragInline, "A"},
		{23, gqlscan.TokenSel, ""},
		{24, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{25, gqlscan.TokenFragRef, "namedFragment1"},
		{26, gqlscan.TokenFragRef, "namedFragment2"},
		{27, gqlscan.TokenSelEnd, ""},
		{28, gqlscan.TokenFragInline, "B"},
		{29, gqlscan.TokenSel, ""},
		{30, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{31, gqlscan.TokenSelEnd, ""},
		{32, gqlscan.TokenField, "baz_fuzz_taz1"},
		{33, gqlscan.TokenArg, "bool"},
		{34, gqlscan.TokenFalse, "false"},
		{35, gqlscan.TokenField, "baz_fuzz_taz2"},
		{36, gqlscan.TokenArg, "bool"},
		{37, gqlscan.TokenTrue, "true"},
		{38, gqlscan.TokenField, "baz_fuzz_taz3"},
		{39, gqlscan.TokenArg, "string"},
		{40, gqlscan.TokenStr, "okay"},
		{41, gqlscan.TokenField, "baz_fuzz_taz4"},
		{42, gqlscan.TokenArg, "array"},
		{43, gqlscan.TokenArr, ""},
		{44, gqlscan.TokenArrEnd, ""},
		{45, gqlscan.TokenField, "baz_fuzz_taz5"},
		{46, gqlscan.TokenArg, "variable"},
		{47, gqlscan.TokenVarRef, "variable"},
		{48, gqlscan.TokenField, "baz_fuzz_taz6"},
		{49, gqlscan.TokenArg, "variable"},
		{50, gqlscan.TokenVarRef, "v"},
		{51, gqlscan.TokenField, "baz_fuzz_taz7"},
		{52, gqlscan.TokenArg, "object"},
		{53, gqlscan.TokenObj, ""},
		{54, gqlscan.TokenObjField, "number0"},
		{55, gqlscan.TokenNum, "0"},
		{56, gqlscan.TokenObjField, "number1"},
		{57, gqlscan.TokenNum, "2"},
		{58, gqlscan.TokenObjField, "number2"},
		{59, gqlscan.TokenNum, "123456789.1234e2"},

		{60, gqlscan.TokenObjField, "arr0"},
		{61, gqlscan.TokenArr, ""},
		{62, gqlscan.TokenArr, ""},
		{63, gqlscan.TokenArrEnd, ""},
		{64, gqlscan.TokenArr, ""},
		{65, gqlscan.TokenObj, ""},
		{66, gqlscan.TokenObjField, "x"},
		{67, gqlscan.TokenNull, "null"},
		{68, gqlscan.TokenObjEnd, ""},
		{69, gqlscan.TokenArrEnd, ""},
		{70, gqlscan.TokenArrEnd, ""},

		{71, gqlscan.TokenObjEnd, ""},
		{72, gqlscan.TokenSelEnd, ""},
		{73, gqlscan.TokenSelEnd, ""},
		{74, gqlscan.TokenSelEnd, ""},

		// Mutation M
		{75, gqlscan.TokenDefMut, ""},
		{76, gqlscan.TokenMutName, "M"},
		{77, gqlscan.TokenVarName, "variable"},
		{78, gqlscan.TokenVarTypeName, "Foo"},
		{79, gqlscan.TokenVarName, "v"},
		{80, gqlscan.TokenVarTypeArr, ""},
		{81, gqlscan.TokenVarTypeArr, ""},
		{82, gqlscan.TokenVarTypeName, "Bar"},
		{83, gqlscan.TokenVarTypeArrEnd, ""},
		{84, gqlscan.TokenVarTypeArrEnd, ""},
		{85, gqlscan.TokenSel, ""},
		{86, gqlscan.TokenField, "foo"},
		{87, gqlscan.TokenArg, "x"},
		{88, gqlscan.TokenNull, "null"},
		{89, gqlscan.TokenSel, ""},
		{90, gqlscan.TokenField, "foo_bar"},
		{91, gqlscan.TokenSelEnd, ""},
		{92, gqlscan.TokenField, "bar"},
		{93, gqlscan.TokenField, "baz"},
		{94, gqlscan.TokenSel, ""},
		{95, gqlscan.TokenField, "baz_fuzz"},
		{96, gqlscan.TokenSel, ""},
		{97, gqlscan.TokenFragInline, "A"},
		{98, gqlscan.TokenSel, ""},
		{99, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{100, gqlscan.TokenFragRef, "namedFragment1"},
		{101, gqlscan.TokenFragRef, "namedFragment2"},
		{102, gqlscan.TokenSelEnd, ""},
		{103, gqlscan.TokenFragInline, "B"},
		{104, gqlscan.TokenSel, ""},
		{105, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{106, gqlscan.TokenSelEnd, ""},
		{107, gqlscan.TokenField, "baz_fuzz_taz1"},
		{108, gqlscan.TokenArg, "bool"},
		{109, gqlscan.TokenFalse, "false"},
		{110, gqlscan.TokenField, "baz_fuzz_taz2"},
		{111, gqlscan.TokenArg, "bool"},
		{112, gqlscan.TokenTrue, "true"},
		{113, gqlscan.TokenField, "baz_fuzz_taz3"},
		{114, gqlscan.TokenArg, "string"},
		{115, gqlscan.TokenStr, "okay"},
		{116, gqlscan.TokenField, "baz_fuzz_taz4"},
		{117, gqlscan.TokenArg, "array"},
		{118, gqlscan.TokenArr, ""},
		{119, gqlscan.TokenArrEnd, ""},
		{120, gqlscan.TokenField, "baz_fuzz_taz5"},
		{121, gqlscan.TokenArg, "variable"},
		{122, gqlscan.TokenVarRef, "variable"},
		{123, gqlscan.TokenField, "baz_fuzz_taz6"},
		{124, gqlscan.TokenArg, "variable"},
		{125, gqlscan.TokenVarRef, "v"},
		{126, gqlscan.TokenField, "baz_fuzz_taz7"},
		{127, gqlscan.TokenArg, "object"},
		{128, gqlscan.TokenObj, ""},
		{129, gqlscan.TokenObjField, "number0"},
		{130, gqlscan.TokenNum, "0"},
		{131, gqlscan.TokenObjField, "number1"},
		{132, gqlscan.TokenNum, "2"},
		{133, gqlscan.TokenObjField, "number2"},
		{134, gqlscan.TokenNum, "123456789.1234e2"},

		{135, gqlscan.TokenObjField, "arr0"},
		{136, gqlscan.TokenArr, ""},
		{137, gqlscan.TokenArr, ""},
		{138, gqlscan.TokenArrEnd, ""},
		{139, gqlscan.TokenArr, ""},
		{140, gqlscan.TokenObj, ""},
		{141, gqlscan.TokenObjField, "x"},
		{142, gqlscan.TokenNull, "null"},
		{143, gqlscan.TokenObjEnd, ""},
		{144, gqlscan.TokenArrEnd, ""},
		{145, gqlscan.TokenArrEnd, ""},

		{146, gqlscan.TokenObjEnd, ""},
		{147, gqlscan.TokenSelEnd, ""},
		{148, gqlscan.TokenSelEnd, ""},
		{149, gqlscan.TokenSelEnd, ""},

		// Fragment f1
		{150, gqlscan.TokenDefFrag, ""},
		{151, gqlscan.TokenFragName, "f1"},
		{152, gqlscan.TokenFragTypeCond, "Query"},
		{153, gqlscan.TokenSel, ""},
		{154, gqlscan.TokenField, "todos"},
		{155, gqlscan.TokenSel, ""},
		{156, gqlscan.TokenFragRef, "f2"},
		{157, gqlscan.TokenSelEnd, ""},
		{158, gqlscan.TokenSelEnd, ""},

		// Query Todos
		{159, gqlscan.TokenDefQry, ""},
		{160, gqlscan.TokenQryName, "Todos"},
		{161, gqlscan.TokenSel, ""},
		{162, gqlscan.TokenFragRef, "f1"},
		{163, gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{164, gqlscan.TokenDefFrag, ""},
		{165, gqlscan.TokenFragName, "f2"},
		{166, gqlscan.TokenFragTypeCond, "Todo"},
		{167, gqlscan.TokenSel, ""},
		{168, gqlscan.TokenField, "id"},
		{169, gqlscan.TokenField, "text"},
		{170, gqlscan.TokenArg, "foo"},
		{171, gqlscan.TokenNum, "2"},
		{172, gqlscan.TokenArg, "bar"},
		{173, gqlscan.TokenStr, "ok"},
		{174, gqlscan.TokenArg, "baz"},
		{175, gqlscan.TokenNull, "null"},
		{176, gqlscan.TokenField, "done"},
		{177, gqlscan.TokenSelEnd, ""},
	}},

	// Comments
	{18, "  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{19, "{  #comment1\n  #comment2\n  x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{20, "{x  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{21, "{x}  #comment1\n  #comment2\n", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{22, "{x(  #comment1\n  #comment2\n  y:0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{23, "{x(y  #comment1\n  #comment2\n  :0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{24, "{x(y:  #comment1\n  #comment2\n  0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{25, "{x(y:0  #comment1\n  #comment2\n  )}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{26, "{x(y:0)  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, "0"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{27, "query  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{28, "mutation  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{29, "fragment  #comment1\n  #comment2\n  f on X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{30, "fragment f  #comment1\n  #comment2\n  on X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{31, "fragment f on  #comment1\n  #comment2\n  X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{32, "fragment f on X  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{33, "{  ...  #comment1\n  #comment2\n  f  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenFragRef, "f"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{34, "{  ...  f  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenFragRef, "f"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{35, "query(  #comment1\n  #comment2\n  $x: T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeName, "T"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{36, "query($x  #comment1\n  #comment2\n  : T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeName, "T"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{37, "query($x:  #comment1\n  #comment2\n  T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeName, "T"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{38, "query($x:[  #comment1\n  #comment2\n  T]){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeArr, ""},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarTypeArrEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{39, "query($x:[T  #comment1\n  #comment2\n  ]){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeArr, ""},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarTypeArrEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{40, "query($x:[T]  #comment1\n  #comment2\n  ){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeArr, ""},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarTypeArrEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{41, "query($x:[T])  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeArr, ""},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarTypeArrEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},

	// String escape
	{42, `{x(s:"\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "s"},
		{4, gqlscan.TokenStr, `\"`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{43, `{x(s:"\\")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "s"},
		{4, gqlscan.TokenStr, `\\`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{44, `{x(s:"\\\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "s"},
		{4, gqlscan.TokenStr, `\\\"`},
		{5, gqlscan.TokenSelEnd, ""},
	}},

	{45, `{x(y:1e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, `1e8`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{46, `{x(y:0e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, `0e8`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{47, `{x(y:0e+8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, `0e+8`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{48, `{x(y:0e-8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArg, "y"},
		{4, gqlscan.TokenNum, `0e-8`},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{49, `mutation{x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{50, `mutation($x:T){x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenVarName, "x"},
		{2, gqlscan.TokenVarTypeName, "T"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{51, `mutation M{x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenMutName, "M"},
		{2, gqlscan.TokenSel, ""},
		{3, gqlscan.TokenField, "x"},
		{4, gqlscan.TokenSelEnd, ""},
	}},
	{52, `{f(o:{o2:{x:[]}})}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "o"},
		{4, gqlscan.TokenObj, ""},
		{5, gqlscan.TokenObjField, "o2"},
		{6, gqlscan.TokenObj, ""},
		{7, gqlscan.TokenObjField, "x"},
		{8, gqlscan.TokenArr, ""},
		{9, gqlscan.TokenArrEnd, ""},
		{10, gqlscan.TokenObjEnd, ""},
		{11, gqlscan.TokenObjEnd, ""},
		{12, gqlscan.TokenSelEnd, ""},
	}},
	{53, `{f(a:[0])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArg, "a"},
		{4, gqlscan.TokenArr, ""},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArrEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
}

func TestScan(t *testing.T) {
	for ti, td := range testdata {
		t.Run("", func(t *testing.T) {
			require.Equal(t, ti, td.index)
			tName := t.Name()
			t.Log(tName)
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
						"unexpected type at index %d", j,
					)
					require.Equal(
						t, td.expect[j].Value, string(i.Value()),
						"unexpected value at index %d", j,
					)
					require.GreaterOrEqual(t, i.IndexHead(), prevHead)
					require.GreaterOrEqual(t, i.IndexHead(), i.IndexTail())
					i.Value()
					require.Equal(t, j, td.expect[j].Index)
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
	}
}

var testdataErr = []struct {
	index     int
	name      string
	input     string
	expectErr string
}{
	{0,
		"unexpected token as query",
		"q",
		"error at index 0 ('q'): unexpected token; expected definition",
	},
	{1,
		"missing square bracket in type",
		"query($a: [A){f}",
		"error at index 11 ('A'): invalid type; " +
			"expected variable type",
	},
	{2,
		"missing square bracket in type",
		"query($a: [[A]){f}",
		"error at index 13 (']'): invalid type; " +
			"expected variable type",
	},
	{3,
		"unexpected square bracket in variable type",
		"query($a: A]){f}",
		"error at index 11 (']'): unexpected token; " +
			"expected variable name",
	},
	{4,
		"unexpected square bracket in variable type",
		"query($a: [[A]]]){f}",
		"error at index 15 (']'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{5,
		"missing query closing curly bracket",
		"{",
		"error at index 1: unexpected end of file; expected selection",
	},
	{6,
		"invalid field name",
		"{1abc}",
		"error at index 1 ('1'): unexpected token; expected field name",
	},
	{7,
		"trailing closing curly bracket",
		"{f}}",
		"error at index 3 ('}'): unexpected token; expected definition",
	},
	{8,
		"query missing closing curly bracket",
		"{}",
		"error at index 1 ('}'): unexpected token; expected field name",
	},
	{9,
		"variable missing name",
		"query($ ",
		"error at index 7 (' '): unexpected token; " +
			"expected variable name",
	},
	{10,
		"empty args",
		"{f()}",
		"error at index 3 (')'): unexpected token; expected argument name",
	},
	{11,
		"argument missing column",
		"{f(x null))}",
		"error at index 5 ('n'): " +
			"unexpected token; expected column after argument name",
	},
	{12,
		"argument with trailing closing parenthesis",
		"{f(x:null))}",
		"error at index 10 (')'): unexpected token; expected field name",
	},
	{13,
		"argument missing closing parenthesis",
		"{f(}",
		"error at index 3 ('}'): unexpected token; expected argument name",
	},
	{14,
		"invalid argument",
		"{f(x:abc))}",
		"error at index 5 ('a'): invalid value; expected value",
	},
	{15,
		"string argument missing closing quotes",
		`{f(x:"))}`,
		"error at index 9: unexpected end of file; expected end of string",
	},
	{16,
		"invalid negative number",
		`{f(x:-))}`,
		"error at index 6 (')'): invalid number value; expected value",
	},
	{17,
		"number missing fraction",
		`{f(x:1.))}`,
		"error at index 7 (')'): invalid number value; expected value",
	},
	{18,
		"number missing exponent",
		`{f(x:1.2e))}`,
		"error at index 9 (')'): invalid number value; expected value",
	},
	{19,
		"number with leading zero",
		`{f(x:0123))}`,
		"error at index 6 ('1'): invalid number value; expected value",
	},

	// --- Unexpected EOF ---
	{20,
		"unexpected EOF",
		"",
		"error at index 0: unexpected end of file; expected definition",
	},
	{21,
		"unexpected EOF",
		"query",
		"error at index 5: unexpected end of file; " +
			"expected variable list or selection set",
	},
	{22,
		"unexpected EOF",
		"query Name",
		"error at index 10: unexpected end of file; " +
			"expected selection set",
	},
	{23,
		"unexpected EOF",
		"query Name ",
		"error at index 11: unexpected end of file; " +
			"expected selection set",
	},
	{24,
		"unexpected EOF",
		"query(",
		"error at index 6: unexpected end of file; " +
			"expected variable name",
	},
	{25,
		"unexpected EOF",
		"query( ",
		"error at index 7: unexpected end of file; " +
			"expected variable name",
	},
	{26,
		"unexpected EOF",
		"query($",
		"error at index 6: unexpected end of file; " +
			"expected variable name",
	},
	{27,
		"unexpected EOF",
		"query($v",
		"error at index 8: unexpected end of file; " +
			"expected column after variable name",
	},
	{28,
		"unexpected EOF",
		"query($v ",
		"error at index 9: unexpected end of file; " +
			"expected column after variable name",
	},
	{29,
		"unexpected EOF",
		"query($v:",
		"error at index 9: unexpected end of file; " +
			"expected variable type",
	},
	{30,
		"unexpected EOF",
		"query($v: ",
		"error at index 10: unexpected end of file; " +
			"expected variable type",
	},
	{31,
		"unexpected EOF",
		"query($v: T",
		"error at index 11: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{32,
		"unexpected EOF",
		"query($v: T ",
		"error at index 12: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{33,
		"unexpected EOF",
		"query($v: T)",
		"error at index 12: unexpected end of file; " +
			"expected selection set",
	},
	{34,
		"unexpected EOF",
		"query($v: T) ",
		"error at index 13: unexpected end of file; " +
			"expected selection set",
	},
	{35,
		"unexpected EOF",
		"{",
		"error at index 1: unexpected end of file; " +
			"expected selection",
	},
	{36,
		"unexpected EOF",
		"{ ",
		"error at index 2: unexpected end of file; " +
			"expected selection",
	},
	{37,
		"unexpected EOF",
		"{foo",
		"error at index 4: unexpected end of file; " +
			"expected field name",
	},
	{38,
		"unexpected EOF",
		"{foo ",
		"error at index 5: unexpected end of file; " +
			"expected field name",
	},
	{39,
		"unexpected EOF",
		"{foo(",
		"error at index 5: unexpected end of file; " +
			"expected argument name",
	},
	{40,
		"unexpected EOF",
		"{foo( ",
		"error at index 6: unexpected end of file; " +
			"expected argument name",
	},
	{41,
		"unexpected EOF",
		"{foo(name",
		"error at index 9: unexpected end of file; " +
			"expected column after argument name",
	},
	{42,
		"unexpected EOF",
		"{foo(name ",
		"error at index 10: unexpected end of file; " +
			"expected column after argument name",
	},
	{43,
		"unexpected EOF",
		"{foo(name:",
		"error at index 10: unexpected end of file; " +
			"expected value",
	},
	{44,
		"unexpected EOF",
		"{foo(name: ",
		"error at index 11: unexpected end of file; " +
			"expected value",
	},
	{45,
		"unexpected EOF",
		"{foo(name: {",
		"error at index 12: unexpected end of file; " +
			"expected object field name",
	},
	{46,
		"unexpected EOF",
		"{foo(name: { ",
		"error at index 13: unexpected end of file; " +
			"expected object field name",
	},
	{47,
		"unexpected EOF",
		"{foo(name: {field",
		"error at index 17: unexpected end of file; " +
			"expected column after object field name",
	},
	{48,
		"unexpected EOF",
		"{foo(name: {field ",
		"error at index 18: unexpected end of file; " +
			"expected column after object field name",
	},
	{49,
		"unexpected EOF",
		"{foo(name: {field:",
		"error at index 18: unexpected end of file; " +
			"expected value",
	},
	{50,
		"unexpected EOF",
		"{foo(name: {field: ",
		"error at index 19: unexpected end of file; " +
			"expected value",
	},
	{51,
		"unexpected EOF",
		`{foo(name: "`,
		"error at index 12: unexpected end of file; " +
			"expected end of string",
	},
	{52,
		"unexpected EOF",
		`{foo(name: ""`,
		"error at index 13: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{53,
		"unexpected EOF",
		`{foo(name: f`,
		"error at index 12: unexpected end of file; expected value",
	},
	{54,
		"unexpected EOF",
		`{foo(name: t`,
		"error at index 12: unexpected end of file; expected value",
	},
	{55,
		"unexpected EOF",
		`{foo(name: n`,
		"error at index 12: unexpected end of file; expected value",
	},
	{56,
		"unexpected EOF",
		`{foo(name: 0`,
		"error at index 12: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{57,
		"unexpected EOF",
		`{foo(name: 0 `,
		"error at index 13: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{58,
		"unexpected EOF",
		`{foo(name: -`,
		"error at index 12: unexpected end of file; expected value",
	},
	{59,
		"unexpected EOF",
		`{foo(name: 0.`,
		"error at index 13: unexpected end of file; expected value",
	},
	{60,
		"unexpected EOF",
		`{foo(name: 0.1e`,
		"error at index 15: unexpected end of file; expected value",
	},
	{61,
		"unexpected EOF",
		`{.`,
		"error at index 2: unexpected end of file; expected fragment",
	},
	{62,
		"unexpected EOF",
		`{..`,
		"error at index 3: unexpected end of file; expected fragment",
	},
	{63,
		"unexpected EOF",
		`{...`,
		"error at index 4: unexpected end of file; expected fragment",
	},
	{64,
		"unexpected EOF",
		`{... `,
		"error at index 5: unexpected end of file; expected fragment",
	},
	{65,
		"unexpected EOF",
		`{... on`,
		"error at index 7: unexpected end of file; expected fragment",
	},
	{66,
		"unexpected EOF",
		`{... on `,
		"error at index 8: unexpected end of file; " +
			"expected inlined fragment",
	},
	{67,
		"unexpected EOF",
		"fragment",
		"error at index 8: unexpected end of file; " +
			"expected fragment name",
	},
	{68,
		"unexpected EOF",
		"{x",
		"error at index 2: unexpected end of file; " +
			"expected field name",
	},

	{69,
		"invalid value",
		"{x(p:falsa",
		"error at index 5 ('f'): invalid value; " +
			"expected value",
	},
	{70,
		"invalid value",
		"{x(p:truu",
		"error at index 5 ('t'): invalid value; " +
			"expected value",
	},
	{71,
		"invalid value",
		"{x(p:nuli",
		"error at index 5 ('n'): invalid value; " +
			"expected value",
	},
	{72,
		"unexpected EOF",
		"{x(p:[",
		"error at index 6: unexpected end of file; " +
			"expected value",
	},
	{73,
		"unexpected token",
		"query($x:T)x",
		"error at index 11 ('x'): unexpected token; " +
			"expected selection set",
	},
	{74,
		"unexpected EOF",
		"mutation M",
		"error at index 10: unexpected end of file; " +
			"expected selection set",
	},
	{75,
		"unexpected token",
		"query\x00",
		"error at index 5 (0x0): unexpected token; " +
			"expected query name",
	},
	{76,
		"unexpected token",
		"{x(y:12e)}",
		"error at index 8 (')'): invalid number value; " +
			"expected value",
	},
	{77,
		"unexpected token",
		"{x(y:12.)}",
		"error at index 8 (')'): invalid number value; " +
			"expected value",
	},
	{78,
		"unexpected token",
		"{x(y:12x)}",
		"error at index 7 ('x'): invalid number value; " +
			"expected value",
	},
	{79,
		"unexpected token",
		"{x(y:12.12x)}",
		"error at index 10 ('x'): invalid number value; " +
			"expected value",
	},
	{80,
		"unexpected EOF",
		"{x(y:12.12",
		"error at index 10: unexpected end of file; " +
			"expected argument list closure or argument",
	},
	{81,
		"unexpected EOF",
		"{x(y:12.",
		"error at index 8: unexpected end of file; " +
			"expected value",
	},
	{82,
		"unexpected token",
		"{x(y:12e111x",
		"error at index 11 ('x'): invalid number value; " +
			"expected value",
	},
	{83,
		"unexpected token",
		"{x(y:12ex",
		"error at index 8 ('x'): invalid number value; " +
			"expected value",
	},
	{84,
		"unexpected token",
		"{x(y:{f})}",
		"error at index 7 ('}'): unexpected token; " +
			"expected column after object field name",
	},
	{85,
		"unexpected token",
		"{x(\x00:1)}",
		"error at index 3 (0x0): unexpected token; " +
			"expected argument name",
	},
	{86,
		"unexpected EOF",
		"{x(y\x00:1)}",
		"error at index 4 (0x0): unexpected token; " +
			"expected argument name",
	},
	{87,
		"unexpected token",
		"query M [",
		"error at index 8 ('['): unexpected token; " +
			"expected selection set",
	},
	{88,
		"unexpected token",
		"mutation M|",
		"error at index 10 ('|'): unexpected token; " +
			"expected selection set",
	},
	{89,
		"unexpected EOF",
		"fragment f on",
		"error at index 13: unexpected end of file; " +
			"expected fragment type condition",
	},
	{90,
		"unexpected token",
		"mutation\x00",
		"error at index 8 (0x0): unexpected token; " +
			"expected mutation name",
	},
	{91,
		"unexpected token",
		"fragment\x00",
		"error at index 8 (0x0): unexpected token; " +
			"expected fragment name",
	},
	{92,
		"unexpected EOF",
		"{x(y:$",
		"error at index 6: unexpected end of file; " +
			"expected referenced variable name",
	},
	{93,
		"unexpected EOF",
		"mutation",
		"error at index 8: unexpected end of file; " +
			"expected variable list or selection set",
	},
	{94,
		"unexpected EOF",
		"{x(y:null)",
		"error at index 10: unexpected end of file; " +
			"expected selection set or selection",
	},
	{95,
		"unexpected token",
		"query($v |",
		"error at index 9 ('|'): unexpected token; " +
			"expected column after variable name",
	},
	{96,
		"unexpected token",
		"query($v:[T] |)",
		"error at index 13 ('|'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{97,
		"unexpected token",
		"fragment X at",
		"error at index 11 ('a'): unexpected token; " +
			"expected keyword 'on'",
	},
	{98,
		"unexpected EOF",
		"query($a:[A]",
		"error at index 12: unexpected end of file; " +
			"expected variable list closure or variable name",
	},
	{99,
		"unexpected EOF",
		"fragment f ",
		"error at index 11: unexpected end of file; " +
			"expected keyword 'on'",
	},
	{100,
		"unexpected EOF",
		"{f{x} ",
		"error at index 6: unexpected end of file; " +
			"expected selection or end of selection set",
	},
	{101,
		"unexpected token",
		"{f(x:\"abc\n\")}",
		"error at index 9 (0xa): unexpected token; " +
			"expected end of string",
	},
	{102,
		"unexpected token",
		"{.f}",
		"error at index 2 ('f'): unexpected token; " +
			"expected fragment",
	},
	{103,
		"unexpected token",
		"{..f}",
		"error at index 3 ('f'): unexpected token; " +
			"expected fragment",
	},
	// {104,
	// 	"unexpected token",
	// 	"{f(o:{o2:{x:[]}})}",
	// 	"error at index 3 ('f'): unexpected token; " +
	// 		"expected fragment",
	// },
}

func TestScanErr(t *testing.T) {
	for ti, td := range testdataErr {
		t.Run("", func(t *testing.T) {
			require.Equal(t, ti, td.index)
			err := gqlscan.Scan(
				[]byte(td.input),
				func(*gqlscan.Iterator) (err bool) {
					return false
				},
			)
			require.Equal(t, td.expectErr, err.Error())
			require.True(t, err.IsErr())
		})
	}
}

func TestScanFuncErr(t *testing.T) {
	const input = `
		{x}
		query Q($variable: Foo, $v: [ [ Bar ] ]) {
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

	const NumOfTokensInInput = 159
	for i := 0; i <= NumOfTokensInInput; i++ {
		c := 0
		err := gqlscan.Scan(
			[]byte(input),
			func(*gqlscan.Iterator) (err bool) {
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
		Index int
		Type  gqlscan.Token
		Value string
		Level int
	}{
		// Query Q
		{0, gqlscan.TokenDefQry, "", 0},
		{1, gqlscan.TokenQryName, "Q", 0},
		{2, gqlscan.TokenVarName, "variable", 0},
		{3, gqlscan.TokenVarTypeName, "Foo", 0},
		{4, gqlscan.TokenVarName, "v", 0},
		{5, gqlscan.TokenVarTypeArr, "", 0},
		{6, gqlscan.TokenVarTypeArr, "", 0},
		{7, gqlscan.TokenVarTypeName, "Bar", 0},
		{8, gqlscan.TokenVarTypeArrEnd, "", 0},
		{9, gqlscan.TokenVarTypeArrEnd, "", 0},
		{10, gqlscan.TokenSel, "", 0},
		{11, gqlscan.TokenField, "foo", 1},
		{12, gqlscan.TokenArg, "x", 1},
		{13, gqlscan.TokenNull, "null", 1},
		{14, gqlscan.TokenSel, "", 1},
		{15, gqlscan.TokenField, "foo_bar", 2},
		{16, gqlscan.TokenSelEnd, "", 2},
		{17, gqlscan.TokenField, "bar", 1},
		{18, gqlscan.TokenField, "baz", 1},
		{19, gqlscan.TokenSel, "", 1},
		{20, gqlscan.TokenField, "baz_fuzz", 2},
		{21, gqlscan.TokenSel, "", 2},
		{22, gqlscan.TokenFragInline, "A", 3},
		{23, gqlscan.TokenSel, "", 3},
		{24, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{25, gqlscan.TokenFragRef, "namedFragment1", 4},
		{26, gqlscan.TokenFragRef, "namedFragment2", 4},
		{27, gqlscan.TokenSelEnd, "", 4},
		{28, gqlscan.TokenFragInline, "B", 3},
		{29, gqlscan.TokenSel, "", 3},
		{30, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{31, gqlscan.TokenSelEnd, "", 4},
		{32, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{33, gqlscan.TokenArg, "bool", 3},
		{34, gqlscan.TokenFalse, "false", 3},
		{35, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{36, gqlscan.TokenArg, "bool", 3},
		{37, gqlscan.TokenTrue, "true", 3},
		{38, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{39, gqlscan.TokenArg, "string", 3},
		{40, gqlscan.TokenStr, "okay", 3},
		{41, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{42, gqlscan.TokenArg, "array", 3},
		{43, gqlscan.TokenArr, "", 3},
		{44, gqlscan.TokenArrEnd, "", 3},
		{45, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{46, gqlscan.TokenArg, "variable", 3},
		{47, gqlscan.TokenVarRef, "variable", 3},
		{48, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{49, gqlscan.TokenArg, "variable", 3},
		{50, gqlscan.TokenVarRef, "v", 3},
		{51, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{52, gqlscan.TokenArg, "object", 3},
		{53, gqlscan.TokenObj, "", 3},
		{54, gqlscan.TokenObjField, "number0", 3},
		{55, gqlscan.TokenNum, "0", 3},
		{56, gqlscan.TokenObjField, "number1", 3},
		{57, gqlscan.TokenNum, "2", 3},
		{58, gqlscan.TokenObjField, "number2", 3},
		{59, gqlscan.TokenNum, "123456789.1234e2", 3},

		{60, gqlscan.TokenObjField, "arr0", 3},
		{61, gqlscan.TokenArr, "", 3},
		{62, gqlscan.TokenArr, "", 3},
		{63, gqlscan.TokenArrEnd, "", 3},
		{64, gqlscan.TokenArr, "", 3},
		{65, gqlscan.TokenObj, "", 3},
		{66, gqlscan.TokenObjField, "x", 3},
		{67, gqlscan.TokenNull, "null", 3},
		{68, gqlscan.TokenObjEnd, "", 3},
		{69, gqlscan.TokenArrEnd, "", 3},
		{70, gqlscan.TokenArrEnd, "", 3},

		{71, gqlscan.TokenObjEnd, "", 3},
		{72, gqlscan.TokenSelEnd, "", 3},
		{73, gqlscan.TokenSelEnd, "", 2},
		{74, gqlscan.TokenSelEnd, "", 1},

		// Mutation M
		{75, gqlscan.TokenDefMut, "", 0},
		{76, gqlscan.TokenMutName, "M", 0},
		{77, gqlscan.TokenVarName, "variable", 0},
		{78, gqlscan.TokenVarTypeName, "Foo", 0},
		{79, gqlscan.TokenVarName, "v", 0},
		{80, gqlscan.TokenVarTypeArr, "", 0},
		{81, gqlscan.TokenVarTypeArr, "", 0},
		{82, gqlscan.TokenVarTypeName, "Bar", 0},
		{83, gqlscan.TokenVarTypeArrEnd, "", 0},
		{84, gqlscan.TokenVarTypeArrEnd, "", 0},
		{85, gqlscan.TokenSel, "", 0},
		{86, gqlscan.TokenField, "foo", 1},
		{87, gqlscan.TokenArg, "x", 1},
		{88, gqlscan.TokenNull, "null", 1},
		{89, gqlscan.TokenSel, "", 1},
		{90, gqlscan.TokenField, "foo_bar", 2},
		{91, gqlscan.TokenSelEnd, "", 2},
		{92, gqlscan.TokenField, "bar", 1},
		{93, gqlscan.TokenField, "baz", 1},
		{94, gqlscan.TokenSel, "", 1},
		{95, gqlscan.TokenField, "baz_fuzz", 2},
		{96, gqlscan.TokenSel, "", 2},
		{97, gqlscan.TokenFragInline, "A", 3},
		{98, gqlscan.TokenSel, "", 3},
		{99, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{100, gqlscan.TokenFragRef, "namedFragment1", 4},
		{101, gqlscan.TokenFragRef, "namedFragment2", 4},
		{102, gqlscan.TokenSelEnd, "", 4},
		{103, gqlscan.TokenFragInline, "B", 3},
		{104, gqlscan.TokenSel, "", 3},
		{105, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{106, gqlscan.TokenSelEnd, "", 4},
		{107, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{108, gqlscan.TokenArg, "bool", 3},
		{109, gqlscan.TokenFalse, "false", 3},
		{110, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{111, gqlscan.TokenArg, "bool", 3},
		{112, gqlscan.TokenTrue, "true", 3},
		{113, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{114, gqlscan.TokenArg, "string", 3},
		{115, gqlscan.TokenStr, "okay", 3},
		{116, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{117, gqlscan.TokenArg, "array", 3},
		{118, gqlscan.TokenArr, "", 3},
		{119, gqlscan.TokenArrEnd, "", 3},
		{120, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{121, gqlscan.TokenArg, "variable", 3},
		{122, gqlscan.TokenVarRef, "variable", 3},
		{123, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{124, gqlscan.TokenArg, "variable", 3},
		{125, gqlscan.TokenVarRef, "v", 3},
		{126, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{127, gqlscan.TokenArg, "object", 3},
		{128, gqlscan.TokenObj, "", 3},
		{129, gqlscan.TokenObjField, "number0", 3},
		{130, gqlscan.TokenNum, "0", 3},
		{131, gqlscan.TokenObjField, "number1", 3},
		{132, gqlscan.TokenNum, "2", 3},
		{133, gqlscan.TokenObjField, "number2", 3},
		{134, gqlscan.TokenNum, "123456789.1234e2", 3},

		{135, gqlscan.TokenObjField, "arr0", 3},
		{136, gqlscan.TokenArr, "", 3},
		{137, gqlscan.TokenArr, "", 3},
		{138, gqlscan.TokenArrEnd, "", 3},
		{139, gqlscan.TokenArr, "", 3},
		{140, gqlscan.TokenObj, "", 3},
		{141, gqlscan.TokenObjField, "x", 3},
		{142, gqlscan.TokenNull, "null", 3},
		{143, gqlscan.TokenObjEnd, "", 3},
		{144, gqlscan.TokenArrEnd, "", 3},
		{145, gqlscan.TokenArrEnd, "", 3},

		{146, gqlscan.TokenObjEnd, "", 3},
		{147, gqlscan.TokenSelEnd, "", 3},
		{148, gqlscan.TokenSelEnd, "", 2},
		{149, gqlscan.TokenSelEnd, "", 1},

		// Fragment f1
		{150, gqlscan.TokenDefFrag, "", 0},
		{151, gqlscan.TokenFragName, "f1", 0},
		{152, gqlscan.TokenFragTypeCond, "Query", 0},
		{153, gqlscan.TokenSel, "", 0},
		{154, gqlscan.TokenField, "todos", 1},
		{155, gqlscan.TokenSel, "", 1},
		{156, gqlscan.TokenFragRef, "f2", 2},
		{157, gqlscan.TokenSelEnd, "", 2},
		{158, gqlscan.TokenSelEnd, "", 1},

		// Query Todos
		{159, gqlscan.TokenDefQry, "", 0},
		{160, gqlscan.TokenQryName, "Todos", 0},
		{161, gqlscan.TokenSel, "", 0},
		{162, gqlscan.TokenFragRef, "f1", 1},
		{163, gqlscan.TokenSelEnd, "", 1},

		// Fragment f2
		{164, gqlscan.TokenDefFrag, "", 0},
		{165, gqlscan.TokenFragName, "f2", 0},
		{166, gqlscan.TokenFragTypeCond, "Todo", 0},
		{167, gqlscan.TokenSel, "", 0},
		{168, gqlscan.TokenField, "id", 1},
		{169, gqlscan.TokenField, "text", 1},
		{170, gqlscan.TokenArg, "foo", 1},
		{171, gqlscan.TokenNum, "2", 1},
		{172, gqlscan.TokenArg, "bar", 1},
		{173, gqlscan.TokenStr, "ok", 1},
		{174, gqlscan.TokenArg, "baz", 1},
		{175, gqlscan.TokenNull, "null", 1},
		{176, gqlscan.TokenField, "done", 1},
		{177, gqlscan.TokenSelEnd, "", 1},
	}

	j := 0
	err := gqlscan.Scan(
		[]byte(input),
		func(i *gqlscan.Iterator) (err bool) {
			require.True(
				t, j < len(expect),
				"exceeding expectation set at: %d {T: %s; V: %s}",
				j, i.Token().String(), i.Value(),
			)
			require.Equal(
				t, expect[j].Type.String(), i.Token().String(),
				"unexpected type at index %d", j,
			)
			require.Equal(
				t, expect[j].Value, string(i.Value()),
				"unexpected value at index %d", j,
			)
			require.Equal(
				t, expect[j].Level, i.LevelSelect(),
				"unexpected selection level at index %d", j,
			)
			i.Value()
			require.Equal(t, j, expect[j].Index)
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
}

func TestZeroValueToString(t *testing.T) {
	var expect gqlscan.Expect
	require.Zero(t, expect.String())

	var token gqlscan.Token
	require.Zero(t, token.String())
}
