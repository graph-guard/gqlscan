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
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenObj, ""},
		{6, gqlscan.TokenObjField, "foo"},
		{7, gqlscan.TokenFalse, "false"},
		{8, gqlscan.TokenObjEnd, ""},
		{9, gqlscan.TokenArgListEnd, ""},
		{10, gqlscan.TokenSelEnd, ""},
	}},
	{3, `{f(f: false)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenFalse, "false"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{4, `{f(f: true)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenTrue, "true"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{5, `{f(f: null)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNull, "null"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{6, `{f(f: [])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenArr, ""},
		{6, gqlscan.TokenArrEnd, ""},
		{7, gqlscan.TokenArgListEnd, ""},
		{8, gqlscan.TokenSelEnd, ""},
	}},
	{7, `{f(f: [[]])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenArr, ""},
		{6, gqlscan.TokenArr, ""},
		{7, gqlscan.TokenArrEnd, ""},
		{8, gqlscan.TokenArrEnd, ""},
		{9, gqlscan.TokenArgListEnd, ""},
		{10, gqlscan.TokenSelEnd, ""},
	}},
	{8, `{f(f: 0)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{9, `{f(f: 0.0)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "0.0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{10, `{f(f: 42)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "42"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{11, `{f(f: -42)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "-42"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{12, `{f(f: -42.5678)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "-42.5678"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{13, `{f(f: -42.5678e2)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenNum, "-42.5678e2"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{14, `{ f (f: {x: 2}) }`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "f"},
		{5, gqlscan.TokenObj, ""},
		{6, gqlscan.TokenObjField, "x"},
		{7, gqlscan.TokenNum, "2"},
		{8, gqlscan.TokenObjEnd, ""},
		{9, gqlscan.TokenArgListEnd, ""},
		{10, gqlscan.TokenSelEnd, ""},
	}},
	{15, `fragment f1 on Query { todos { ...f2 } }
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
	{16, `query Q($variable: Foo, $v: [ [ Bar ] ]) {
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
		{2, gqlscan.TokenVarList, ""},
		{3, gqlscan.TokenVarName, "variable"},
		{4, gqlscan.TokenVarTypeName, "Foo"},
		{5, gqlscan.TokenVarName, "v"},
		{6, gqlscan.TokenVarTypeArr, ""},
		{7, gqlscan.TokenVarTypeArr, ""},
		{8, gqlscan.TokenVarTypeName, "Bar"},
		{9, gqlscan.TokenVarTypeArrEnd, ""},
		{10, gqlscan.TokenVarTypeArrEnd, ""},
		{11, gqlscan.TokenVarListEnd, ""},
		{12, gqlscan.TokenSel, ""},
		{13, gqlscan.TokenField, "foo"},
		{14, gqlscan.TokenArgList, ""},
		{15, gqlscan.TokenArg, "x"},
		{16, gqlscan.TokenNull, "null"},
		{17, gqlscan.TokenArgListEnd, ""},
		{18, gqlscan.TokenSel, ""},
		{19, gqlscan.TokenField, "foo_bar"},
		{20, gqlscan.TokenSelEnd, ""},
		{21, gqlscan.TokenField, "bar"},
		{22, gqlscan.TokenField, "baz"},
		{23, gqlscan.TokenSel, ""},
		{24, gqlscan.TokenField, "baz_fuzz"},
		{25, gqlscan.TokenSel, ""},
		{26, gqlscan.TokenFragInline, "A"},
		{27, gqlscan.TokenSel, ""},
		{28, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{29, gqlscan.TokenFragRef, "namedFragment1"},
		{30, gqlscan.TokenFragRef, "namedFragment2"},
		{31, gqlscan.TokenSelEnd, ""},
		{32, gqlscan.TokenFragInline, "B"},
		{33, gqlscan.TokenSel, ""},
		{34, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{35, gqlscan.TokenSelEnd, ""},
		{36, gqlscan.TokenField, "baz_fuzz_taz1"},
		{37, gqlscan.TokenArgList, ""},
		{38, gqlscan.TokenArg, "bool"},
		{39, gqlscan.TokenFalse, "false"},
		{40, gqlscan.TokenArgListEnd, ""},
		{41, gqlscan.TokenField, "baz_fuzz_taz2"},
		{42, gqlscan.TokenArgList, ""},
		{43, gqlscan.TokenArg, "bool"},
		{44, gqlscan.TokenTrue, "true"},
		{45, gqlscan.TokenArgListEnd, ""},
		{46, gqlscan.TokenField, "baz_fuzz_taz3"},
		{47, gqlscan.TokenArgList, ""},
		{48, gqlscan.TokenArg, "string"},
		{49, gqlscan.TokenStr, "okay"},
		{50, gqlscan.TokenArgListEnd, ""},
		{51, gqlscan.TokenField, "baz_fuzz_taz4"},
		{52, gqlscan.TokenArgList, ""},
		{53, gqlscan.TokenArg, "array"},
		{54, gqlscan.TokenArr, ""},
		{55, gqlscan.TokenArrEnd, ""},
		{56, gqlscan.TokenArgListEnd, ""},
		{57, gqlscan.TokenField, "baz_fuzz_taz5"},
		{58, gqlscan.TokenArgList, ""},
		{59, gqlscan.TokenArg, "variable"},
		{60, gqlscan.TokenVarRef, "variable"},
		{61, gqlscan.TokenArgListEnd, ""},
		{62, gqlscan.TokenField, "baz_fuzz_taz6"},
		{63, gqlscan.TokenArgList, ""},
		{64, gqlscan.TokenArg, "variable"},
		{65, gqlscan.TokenVarRef, "v"},
		{66, gqlscan.TokenArgListEnd, ""},
		{67, gqlscan.TokenField, "baz_fuzz_taz7"},
		{68, gqlscan.TokenArgList, ""},
		{69, gqlscan.TokenArg, "object"},
		{70, gqlscan.TokenObj, ""},
		{71, gqlscan.TokenObjField, "number0"},
		{72, gqlscan.TokenNum, "0"},
		{73, gqlscan.TokenObjField, "number1"},
		{74, gqlscan.TokenNum, "2"},
		{75, gqlscan.TokenObjField, "number2"},
		{76, gqlscan.TokenNum, "123456789.1234e2"},

		{77, gqlscan.TokenObjField, "arr0"},
		{78, gqlscan.TokenArr, ""},
		{79, gqlscan.TokenArr, ""},
		{80, gqlscan.TokenArrEnd, ""},
		{81, gqlscan.TokenArr, ""},
		{82, gqlscan.TokenObj, ""},
		{83, gqlscan.TokenObjField, "x"},
		{84, gqlscan.TokenNull, "null"},
		{85, gqlscan.TokenObjEnd, ""},
		{86, gqlscan.TokenArrEnd, ""},
		{87, gqlscan.TokenArrEnd, ""},

		{88, gqlscan.TokenObjEnd, ""},
		{89, gqlscan.TokenArgListEnd, ""},
		{90, gqlscan.TokenSelEnd, ""},
		{91, gqlscan.TokenSelEnd, ""},
		{92, gqlscan.TokenSelEnd, ""},

		// Mutation M
		{93, gqlscan.TokenDefMut, ""},
		{94, gqlscan.TokenMutName, "M"},
		{95, gqlscan.TokenVarList, ""},
		{96, gqlscan.TokenVarName, "variable"},
		{97, gqlscan.TokenVarTypeName, "Foo"},
		{98, gqlscan.TokenVarName, "v"},
		{99, gqlscan.TokenVarTypeArr, ""},
		{100, gqlscan.TokenVarTypeArr, ""},
		{101, gqlscan.TokenVarTypeName, "Bar"},
		{102, gqlscan.TokenVarTypeArrEnd, ""},
		{103, gqlscan.TokenVarTypeArrEnd, ""},
		{104, gqlscan.TokenVarListEnd, ""},
		{105, gqlscan.TokenSel, ""},
		{106, gqlscan.TokenField, "foo"},
		{107, gqlscan.TokenArgList, ""},
		{108, gqlscan.TokenArg, "x"},
		{109, gqlscan.TokenNull, "null"},
		{110, gqlscan.TokenArgListEnd, ""},
		{111, gqlscan.TokenSel, ""},
		{112, gqlscan.TokenField, "foo_bar"},
		{113, gqlscan.TokenSelEnd, ""},
		{114, gqlscan.TokenField, "bar"},
		{115, gqlscan.TokenField, "baz"},
		{116, gqlscan.TokenSel, ""},
		{117, gqlscan.TokenField, "baz_fuzz"},
		{118, gqlscan.TokenSel, ""},
		{119, gqlscan.TokenFragInline, "A"},
		{120, gqlscan.TokenSel, ""},
		{121, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{122, gqlscan.TokenFragRef, "namedFragment1"},
		{123, gqlscan.TokenFragRef, "namedFragment2"},
		{124, gqlscan.TokenSelEnd, ""},
		{125, gqlscan.TokenFragInline, "B"},
		{126, gqlscan.TokenSel, ""},
		{127, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{128, gqlscan.TokenSelEnd, ""},
		{129, gqlscan.TokenField, "baz_fuzz_taz1"},
		{130, gqlscan.TokenArgList, ""},
		{131, gqlscan.TokenArg, "bool"},
		{132, gqlscan.TokenFalse, "false"},
		{133, gqlscan.TokenArgListEnd, ""},
		{134, gqlscan.TokenField, "baz_fuzz_taz2"},
		{135, gqlscan.TokenArgList, ""},
		{136, gqlscan.TokenArg, "bool"},
		{137, gqlscan.TokenTrue, "true"},
		{138, gqlscan.TokenArgListEnd, ""},
		{139, gqlscan.TokenField, "baz_fuzz_taz3"},
		{140, gqlscan.TokenArgList, ""},
		{141, gqlscan.TokenArg, "string"},
		{142, gqlscan.TokenStr, "okay"},
		{143, gqlscan.TokenArgListEnd, ""},
		{144, gqlscan.TokenField, "baz_fuzz_taz4"},
		{145, gqlscan.TokenArgList, ""},
		{146, gqlscan.TokenArg, "array"},
		{147, gqlscan.TokenArr, ""},
		{148, gqlscan.TokenArrEnd, ""},
		{149, gqlscan.TokenArgListEnd, ""},
		{150, gqlscan.TokenField, "baz_fuzz_taz5"},
		{151, gqlscan.TokenArgList, ""},
		{152, gqlscan.TokenArg, "variable"},
		{153, gqlscan.TokenVarRef, "variable"},
		{154, gqlscan.TokenArgListEnd, ""},
		{155, gqlscan.TokenField, "baz_fuzz_taz6"},
		{156, gqlscan.TokenArgList, ""},
		{157, gqlscan.TokenArg, "variable"},
		{158, gqlscan.TokenVarRef, "v"},
		{159, gqlscan.TokenArgListEnd, ""},
		{160, gqlscan.TokenField, "baz_fuzz_taz7"},
		{161, gqlscan.TokenArgList, ""},
		{162, gqlscan.TokenArg, "object"},
		{163, gqlscan.TokenObj, ""},
		{164, gqlscan.TokenObjField, "number0"},
		{165, gqlscan.TokenNum, "0"},
		{166, gqlscan.TokenObjField, "number1"},
		{167, gqlscan.TokenNum, "2"},
		{168, gqlscan.TokenObjField, "number2"},
		{169, gqlscan.TokenNum, "123456789.1234e2"},

		{170, gqlscan.TokenObjField, "arr0"},
		{171, gqlscan.TokenArr, ""},
		{172, gqlscan.TokenArr, ""},
		{173, gqlscan.TokenArrEnd, ""},
		{174, gqlscan.TokenArr, ""},
		{175, gqlscan.TokenObj, ""},
		{176, gqlscan.TokenObjField, "x"},
		{177, gqlscan.TokenNull, "null"},
		{178, gqlscan.TokenObjEnd, ""},
		{179, gqlscan.TokenArrEnd, ""},
		{180, gqlscan.TokenArrEnd, ""},

		{181, gqlscan.TokenObjEnd, ""},
		{182, gqlscan.TokenArgListEnd, ""},
		{183, gqlscan.TokenSelEnd, ""},
		{184, gqlscan.TokenSelEnd, ""},
		{185, gqlscan.TokenSelEnd, ""},

		// Fragment f1
		{186, gqlscan.TokenDefFrag, ""},
		{187, gqlscan.TokenFragName, "f1"},
		{188, gqlscan.TokenFragTypeCond, "Query"},
		{189, gqlscan.TokenSel, ""},
		{190, gqlscan.TokenField, "todos"},
		{191, gqlscan.TokenSel, ""},
		{192, gqlscan.TokenFragRef, "f2"},
		{193, gqlscan.TokenSelEnd, ""},
		{194, gqlscan.TokenSelEnd, ""},

		// Query Todos
		{195, gqlscan.TokenDefQry, ""},
		{196, gqlscan.TokenQryName, "Todos"},
		{197, gqlscan.TokenSel, ""},
		{198, gqlscan.TokenFragRef, "f1"},
		{199, gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{200, gqlscan.TokenDefFrag, ""},
		{201, gqlscan.TokenFragName, "f2"},
		{202, gqlscan.TokenFragTypeCond, "Todo"},
		{203, gqlscan.TokenSel, ""},
		{204, gqlscan.TokenField, "id"},
		{205, gqlscan.TokenField, "text"},
		{206, gqlscan.TokenArgList, ""},
		{207, gqlscan.TokenArg, "foo"},
		{208, gqlscan.TokenNum, "2"},
		{209, gqlscan.TokenArg, "bar"},
		{210, gqlscan.TokenStr, "ok"},
		{211, gqlscan.TokenArg, "baz"},
		{212, gqlscan.TokenNull, "null"},
		{213, gqlscan.TokenArgListEnd, ""},
		{214, gqlscan.TokenField, "done"},
		{215, gqlscan.TokenSelEnd, ""},
	}},

	// Comments
	{17, "  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{18, "{  #comment1\n  #comment2\n  x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{19, "{x  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{20, "{x}  #comment1\n  #comment2\n", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{21, "{x(  #comment1\n  #comment2\n  y:0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{22, "{x(y  #comment1\n  #comment2\n  :0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{23, "{x(y:  #comment1\n  #comment2\n  0)}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{24, "{x(y:0  #comment1\n  #comment2\n  )}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{25, "{x(y:0)  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, "0"},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{26, "query  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{27, "mutation  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{28, "fragment  #comment1\n  #comment2\n  f on X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{29, "fragment f  #comment1\n  #comment2\n  on X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{30, "fragment f on  #comment1\n  #comment2\n  X{x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{31, "fragment f on X  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefFrag, ""},
		{1, gqlscan.TokenFragName, "f"},
		{2, gqlscan.TokenFragTypeCond, "X"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "x"},
		{5, gqlscan.TokenSelEnd, ""},
	}},
	{32, "{  ...  #comment1\n  #comment2\n  f  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenFragRef, "f"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{33, "{  ...  f  #comment1\n  #comment2\n  }", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenFragRef, "f"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{34, "query(  #comment1\n  #comment2\n  $x: T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarListEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{35, "query($x  #comment1\n  #comment2\n  : T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarListEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{36, "query($x:  #comment1\n  #comment2\n  T){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarListEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{37, "query($x:[  #comment1\n  #comment2\n  T]){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeArr, ""},
		{4, gqlscan.TokenVarTypeName, "T"},
		{5, gqlscan.TokenVarTypeArrEnd, ""},
		{6, gqlscan.TokenVarListEnd, ""},
		{7, gqlscan.TokenSel, ""},
		{8, gqlscan.TokenField, "x"},
		{9, gqlscan.TokenSelEnd, ""},
	}},
	{38, "query($x:[T  #comment1\n  #comment2\n  ]){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeArr, ""},
		{4, gqlscan.TokenVarTypeName, "T"},
		{5, gqlscan.TokenVarTypeArrEnd, ""},
		{6, gqlscan.TokenVarListEnd, ""},
		{7, gqlscan.TokenSel, ""},
		{8, gqlscan.TokenField, "x"},
		{9, gqlscan.TokenSelEnd, ""},
	}},
	{39, "query($x:[T]  #comment1\n  #comment2\n  ){x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeArr, ""},
		{4, gqlscan.TokenVarTypeName, "T"},
		{5, gqlscan.TokenVarTypeArrEnd, ""},
		{6, gqlscan.TokenVarListEnd, ""},
		{7, gqlscan.TokenSel, ""},
		{8, gqlscan.TokenField, "x"},
		{9, gqlscan.TokenSelEnd, ""},
	}},
	{40, "query($x:[T])  #comment1\n  #comment2\n  {x}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeArr, ""},
		{4, gqlscan.TokenVarTypeName, "T"},
		{5, gqlscan.TokenVarTypeArrEnd, ""},
		{6, gqlscan.TokenVarListEnd, ""},
		{7, gqlscan.TokenSel, ""},
		{8, gqlscan.TokenField, "x"},
		{9, gqlscan.TokenSelEnd, ""},
	}},

	// String escape
	{41, `{x(s:"\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\"`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{42, `{x(s:"\\")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\\`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{43, `{x(s:"\\\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\\\"`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},

	{44, `{x(y:1e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `1e8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{45, `{x(y:0e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{46, `{x(y:0e+8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e+8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{47, `{x(y:0e-8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e-8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{48, `mutation{x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{49, `mutation($x:T){x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarListEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{50, `mutation M{x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenMutName, "M"},
		{2, gqlscan.TokenSel, ""},
		{3, gqlscan.TokenField, "x"},
		{4, gqlscan.TokenSelEnd, ""},
	}},
	{51, `{f(o:{o2:{x:[]}})}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "o"},
		{5, gqlscan.TokenObj, ""},
		{6, gqlscan.TokenObjField, "o2"},
		{7, gqlscan.TokenObj, ""},
		{8, gqlscan.TokenObjField, "x"},
		{9, gqlscan.TokenArr, ""},
		{10, gqlscan.TokenArrEnd, ""},
		{11, gqlscan.TokenObjEnd, ""},
		{12, gqlscan.TokenObjEnd, ""},
		{13, gqlscan.TokenArgListEnd, ""},
		{14, gqlscan.TokenSelEnd, ""},
	}},
	{52, `{f(a:[0])}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "a"},
		{5, gqlscan.TokenArr, ""},
		{6, gqlscan.TokenNum, "0"},
		{7, gqlscan.TokenArrEnd, ""},
		{8, gqlscan.TokenArgListEnd, ""},
		{9, gqlscan.TokenSelEnd, ""},
	}},
	{53, `query($v:T ! ){x(a:$v)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "v"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarTypeNotNull, ""},
		{5, gqlscan.TokenVarListEnd, ""},
		{6, gqlscan.TokenSel, ""},
		{7, gqlscan.TokenField, "x"},
		{8, gqlscan.TokenArgList, ""},
		{9, gqlscan.TokenArg, "a"},
		{10, gqlscan.TokenVarRef, "v"},
		{11, gqlscan.TokenArgListEnd, ""},
		{12, gqlscan.TokenSelEnd, ""},
	}},
	{54, `query ($v: [ [ T ! ] ! ] ! ) {x(a:$v)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "v"},
		{3, gqlscan.TokenVarTypeArr, ""},
		{4, gqlscan.TokenVarTypeArr, ""},
		{5, gqlscan.TokenVarTypeName, "T"},
		{6, gqlscan.TokenVarTypeNotNull, ""},
		{7, gqlscan.TokenVarTypeArrEnd, ""},
		{8, gqlscan.TokenVarTypeNotNull, ""},
		{9, gqlscan.TokenVarTypeArrEnd, ""},
		{10, gqlscan.TokenVarTypeNotNull, ""},
		{11, gqlscan.TokenVarListEnd, ""},
		{12, gqlscan.TokenSel, ""},
		{13, gqlscan.TokenField, "x"},
		{14, gqlscan.TokenArgList, ""},
		{15, gqlscan.TokenArg, "a"},
		{16, gqlscan.TokenVarRef, "v"},
		{17, gqlscan.TokenArgListEnd, ""},
		{18, gqlscan.TokenSelEnd, ""},
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
	{104,
		"unexpected token",
		"query($v:T ! !){x(a:$v)}",
		"error at index 13 ('!'): unexpected token; " +
			"expected variable list closure or variable name",
	},
	{105,
		"unexpected token",
		"query($v: [ T ! ] ! ! ){x(a:$v)}",
		"error at index 20 ('!'): unexpected token; " +
			"expected variable list closure or variable name",
	},
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
		query($v: [T!]!) {x}
		mutation($v: [T!]!) {x}
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
		{2, gqlscan.TokenVarList, "", 0},
		{3, gqlscan.TokenVarName, "variable", 0},
		{4, gqlscan.TokenVarTypeName, "Foo", 0},
		{5, gqlscan.TokenVarName, "v", 0},
		{6, gqlscan.TokenVarTypeArr, "", 0},
		{7, gqlscan.TokenVarTypeArr, "", 0},
		{8, gqlscan.TokenVarTypeName, "Bar", 0},
		{9, gqlscan.TokenVarTypeArrEnd, "", 0},
		{10, gqlscan.TokenVarTypeArrEnd, "", 0},
		{11, gqlscan.TokenVarListEnd, "", 0},
		{12, gqlscan.TokenSel, "", 0},
		{13, gqlscan.TokenField, "foo", 1},
		{14, gqlscan.TokenArgList, "", 1},
		{15, gqlscan.TokenArg, "x", 1},
		{16, gqlscan.TokenNull, "null", 1},
		{17, gqlscan.TokenArgListEnd, "", 1},

		{18, gqlscan.TokenSel, "", 1},
		{19, gqlscan.TokenField, "foo_bar", 2},
		{20, gqlscan.TokenSelEnd, "", 2},
		{21, gqlscan.TokenField, "bar", 1},
		{22, gqlscan.TokenField, "baz", 1},
		{23, gqlscan.TokenSel, "", 1},
		{24, gqlscan.TokenField, "baz_fuzz", 2},
		{25, gqlscan.TokenSel, "", 2},
		{26, gqlscan.TokenFragInline, "A", 3},
		{27, gqlscan.TokenSel, "", 3},
		{28, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{29, gqlscan.TokenFragRef, "namedFragment1", 4},
		{30, gqlscan.TokenFragRef, "namedFragment2", 4},
		{31, gqlscan.TokenSelEnd, "", 4},
		{32, gqlscan.TokenFragInline, "B", 3},
		{33, gqlscan.TokenSel, "", 3},
		{34, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{35, gqlscan.TokenSelEnd, "", 4},
		{36, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{37, gqlscan.TokenArgList, "", 3},
		{38, gqlscan.TokenArg, "bool", 3},
		{39, gqlscan.TokenFalse, "false", 3},
		{40, gqlscan.TokenArgListEnd, "", 3},
		{41, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{42, gqlscan.TokenArgList, "", 3},
		{43, gqlscan.TokenArg, "bool", 3},
		{44, gqlscan.TokenTrue, "true", 3},
		{45, gqlscan.TokenArgListEnd, "", 3},
		{46, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{47, gqlscan.TokenArgList, "", 3},
		{48, gqlscan.TokenArg, "string", 3},
		{49, gqlscan.TokenStr, "okay", 3},
		{50, gqlscan.TokenArgListEnd, "", 3},
		{51, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{52, gqlscan.TokenArgList, "", 3},
		{53, gqlscan.TokenArg, "array", 3},
		{54, gqlscan.TokenArr, "", 3},
		{55, gqlscan.TokenArrEnd, "", 3},
		{56, gqlscan.TokenArgListEnd, "", 3},
		{57, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{58, gqlscan.TokenArgList, "", 3},
		{59, gqlscan.TokenArg, "variable", 3},
		{60, gqlscan.TokenVarRef, "variable", 3},
		{61, gqlscan.TokenArgListEnd, "", 3},
		{62, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{63, gqlscan.TokenArgList, "", 3},
		{64, gqlscan.TokenArg, "variable", 3},
		{65, gqlscan.TokenVarRef, "v", 3},
		{66, gqlscan.TokenArgListEnd, "", 3},
		{67, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{68, gqlscan.TokenArgList, "", 3},
		{69, gqlscan.TokenArg, "object", 3},
		{70, gqlscan.TokenObj, "", 3},
		{71, gqlscan.TokenObjField, "number0", 3},
		{72, gqlscan.TokenNum, "0", 3},
		{73, gqlscan.TokenObjField, "number1", 3},
		{74, gqlscan.TokenNum, "2", 3},
		{75, gqlscan.TokenObjField, "number2", 3},
		{76, gqlscan.TokenNum, "123456789.1234e2", 3},

		{77, gqlscan.TokenObjField, "arr0", 3},
		{78, gqlscan.TokenArr, "", 3},
		{79, gqlscan.TokenArr, "", 3},
		{80, gqlscan.TokenArrEnd, "", 3},
		{81, gqlscan.TokenArr, "", 3},
		{82, gqlscan.TokenObj, "", 3},
		{83, gqlscan.TokenObjField, "x", 3},
		{84, gqlscan.TokenNull, "null", 3},
		{85, gqlscan.TokenObjEnd, "", 3},
		{86, gqlscan.TokenArrEnd, "", 3},
		{87, gqlscan.TokenArrEnd, "", 3},

		{88, gqlscan.TokenObjEnd, "", 3},
		{89, gqlscan.TokenArgListEnd, "", 3},
		{90, gqlscan.TokenSelEnd, "", 3},
		{91, gqlscan.TokenSelEnd, "", 2},
		{92, gqlscan.TokenSelEnd, "", 1},

		// Mutation M
		{93, gqlscan.TokenDefMut, "", 0},
		{94, gqlscan.TokenMutName, "M", 0},
		{95, gqlscan.TokenVarList, "", 0},
		{96, gqlscan.TokenVarName, "variable", 0},
		{97, gqlscan.TokenVarTypeName, "Foo", 0},
		{98, gqlscan.TokenVarName, "v", 0},
		{99, gqlscan.TokenVarTypeArr, "", 0},
		{100, gqlscan.TokenVarTypeArr, "", 0},
		{101, gqlscan.TokenVarTypeName, "Bar", 0},
		{102, gqlscan.TokenVarTypeArrEnd, "", 0},
		{103, gqlscan.TokenVarTypeArrEnd, "", 0},
		{104, gqlscan.TokenVarListEnd, "", 0},
		{105, gqlscan.TokenSel, "", 0},
		{106, gqlscan.TokenField, "foo", 1},
		{107, gqlscan.TokenArgList, "", 1},
		{108, gqlscan.TokenArg, "x", 1},
		{109, gqlscan.TokenNull, "null", 1},
		{110, gqlscan.TokenArgListEnd, "", 1},
		{111, gqlscan.TokenSel, "", 1},
		{112, gqlscan.TokenField, "foo_bar", 2},
		{113, gqlscan.TokenSelEnd, "", 2},
		{114, gqlscan.TokenField, "bar", 1},
		{115, gqlscan.TokenField, "baz", 1},
		{116, gqlscan.TokenSel, "", 1},
		{117, gqlscan.TokenField, "baz_fuzz", 2},
		{118, gqlscan.TokenSel, "", 2},
		{119, gqlscan.TokenFragInline, "A", 3},
		{120, gqlscan.TokenSel, "", 3},
		{121, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{122, gqlscan.TokenFragRef, "namedFragment1", 4},
		{123, gqlscan.TokenFragRef, "namedFragment2", 4},
		{124, gqlscan.TokenSelEnd, "", 4},
		{125, gqlscan.TokenFragInline, "B", 3},
		{126, gqlscan.TokenSel, "", 3},
		{127, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{128, gqlscan.TokenSelEnd, "", 4},
		{129, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{130, gqlscan.TokenArgList, "", 3},
		{131, gqlscan.TokenArg, "bool", 3},
		{132, gqlscan.TokenFalse, "false", 3},
		{133, gqlscan.TokenArgListEnd, "", 3},
		{134, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{135, gqlscan.TokenArgList, "", 3},
		{136, gqlscan.TokenArg, "bool", 3},
		{137, gqlscan.TokenTrue, "true", 3},
		{138, gqlscan.TokenArgListEnd, "", 3},
		{139, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{140, gqlscan.TokenArgList, "", 3},
		{141, gqlscan.TokenArg, "string", 3},
		{142, gqlscan.TokenStr, "okay", 3},
		{143, gqlscan.TokenArgListEnd, "", 3},
		{144, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{145, gqlscan.TokenArgList, "", 3},
		{146, gqlscan.TokenArg, "array", 3},
		{147, gqlscan.TokenArr, "", 3},
		{148, gqlscan.TokenArrEnd, "", 3},
		{149, gqlscan.TokenArgListEnd, "", 3},
		{150, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{151, gqlscan.TokenArgList, "", 3},
		{152, gqlscan.TokenArg, "variable", 3},
		{153, gqlscan.TokenVarRef, "variable", 3},
		{154, gqlscan.TokenArgListEnd, "", 3},
		{155, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{156, gqlscan.TokenArgList, "", 3},
		{157, gqlscan.TokenArg, "variable", 3},
		{158, gqlscan.TokenVarRef, "v", 3},
		{159, gqlscan.TokenArgListEnd, "", 3},
		{160, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{161, gqlscan.TokenArgList, "", 3},
		{162, gqlscan.TokenArg, "object", 3},
		{163, gqlscan.TokenObj, "", 3},
		{164, gqlscan.TokenObjField, "number0", 3},
		{165, gqlscan.TokenNum, "0", 3},
		{166, gqlscan.TokenObjField, "number1", 3},
		{167, gqlscan.TokenNum, "2", 3},
		{168, gqlscan.TokenObjField, "number2", 3},
		{169, gqlscan.TokenNum, "123456789.1234e2", 3},

		{170, gqlscan.TokenObjField, "arr0", 3},
		{171, gqlscan.TokenArr, "", 3},
		{172, gqlscan.TokenArr, "", 3},
		{173, gqlscan.TokenArrEnd, "", 3},
		{174, gqlscan.TokenArr, "", 3},
		{175, gqlscan.TokenObj, "", 3},
		{176, gqlscan.TokenObjField, "x", 3},
		{177, gqlscan.TokenNull, "null", 3},
		{178, gqlscan.TokenObjEnd, "", 3},
		{179, gqlscan.TokenArrEnd, "", 3},
		{180, gqlscan.TokenArrEnd, "", 3},

		{181, gqlscan.TokenObjEnd, "", 3},
		{182, gqlscan.TokenArgListEnd, "", 3},
		{183, gqlscan.TokenSelEnd, "", 3},
		{184, gqlscan.TokenSelEnd, "", 2},
		{185, gqlscan.TokenSelEnd, "", 1},

		// Fragment f1
		{186, gqlscan.TokenDefFrag, "", 0},
		{187, gqlscan.TokenFragName, "f1", 0},
		{188, gqlscan.TokenFragTypeCond, "Query", 0},
		{189, gqlscan.TokenSel, "", 0},
		{190, gqlscan.TokenField, "todos", 1},
		{191, gqlscan.TokenSel, "", 1},
		{192, gqlscan.TokenFragRef, "f2", 2},
		{193, gqlscan.TokenSelEnd, "", 2},
		{194, gqlscan.TokenSelEnd, "", 1},

		// Query Todos
		{195, gqlscan.TokenDefQry, "", 0},
		{196, gqlscan.TokenQryName, "Todos", 0},
		{197, gqlscan.TokenSel, "", 0},
		{198, gqlscan.TokenFragRef, "f1", 1},
		{199, gqlscan.TokenSelEnd, "", 1},

		// Fragment f2
		{200, gqlscan.TokenDefFrag, "", 0},
		{201, gqlscan.TokenFragName, "f2", 0},
		{202, gqlscan.TokenFragTypeCond, "Todo", 0},
		{203, gqlscan.TokenSel, "", 0},
		{204, gqlscan.TokenField, "id", 1},
		{205, gqlscan.TokenField, "text", 1},
		{206, gqlscan.TokenArgList, "", 1},
		{207, gqlscan.TokenArg, "foo", 1},
		{208, gqlscan.TokenNum, "2", 1},
		{209, gqlscan.TokenArg, "bar", 1},
		{210, gqlscan.TokenStr, "ok", 1},
		{211, gqlscan.TokenArg, "baz", 1},
		{212, gqlscan.TokenNull, "null", 1},
		{213, gqlscan.TokenArgListEnd, "", 1},
		{214, gqlscan.TokenField, "done", 1},
		{215, gqlscan.TokenSelEnd, "", 1},
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
