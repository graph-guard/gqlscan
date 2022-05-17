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
		{13, gqlscan.TokenFieldAlias, "foo_alias"},
		{14, gqlscan.TokenField, "foo"},
		{15, gqlscan.TokenArgList, ""},
		{16, gqlscan.TokenArg, "x"},
		{17, gqlscan.TokenNull, "null"},
		{18, gqlscan.TokenArgListEnd, ""},
		{19, gqlscan.TokenSel, ""},
		{20, gqlscan.TokenFieldAlias, "foobar_alias"},
		{21, gqlscan.TokenField, "foo_bar"},
		{22, gqlscan.TokenSelEnd, ""},
		{23, gqlscan.TokenField, "bar"},
		{24, gqlscan.TokenField, "baz"},
		{25, gqlscan.TokenSel, ""},
		{26, gqlscan.TokenField, "baz_fuzz"},
		{27, gqlscan.TokenSel, ""},
		{28, gqlscan.TokenFragInline, "A"},
		{29, gqlscan.TokenSel, ""},
		{30, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{31, gqlscan.TokenFragRef, "namedFragment1"},
		{32, gqlscan.TokenFragRef, "namedFragment2"},
		{33, gqlscan.TokenSelEnd, ""},
		{34, gqlscan.TokenFragInline, "B"},
		{35, gqlscan.TokenSel, ""},
		{36, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{37, gqlscan.TokenSelEnd, ""},
		{38, gqlscan.TokenField, "baz_fuzz_taz1"},
		{39, gqlscan.TokenArgList, ""},
		{40, gqlscan.TokenArg, "bool"},
		{41, gqlscan.TokenFalse, "false"},
		{42, gqlscan.TokenArgListEnd, ""},
		{43, gqlscan.TokenField, "baz_fuzz_taz2"},
		{44, gqlscan.TokenArgList, ""},
		{45, gqlscan.TokenArg, "bool"},
		{46, gqlscan.TokenTrue, "true"},
		{47, gqlscan.TokenArgListEnd, ""},
		{48, gqlscan.TokenField, "baz_fuzz_taz3"},
		{49, gqlscan.TokenArgList, ""},
		{50, gqlscan.TokenArg, "string"},
		{51, gqlscan.TokenStr, "okay"},
		{52, gqlscan.TokenArgListEnd, ""},
		{53, gqlscan.TokenField, "baz_fuzz_taz4"},
		{54, gqlscan.TokenArgList, ""},
		{55, gqlscan.TokenArg, "array"},
		{56, gqlscan.TokenArr, ""},
		{57, gqlscan.TokenArrEnd, ""},
		{58, gqlscan.TokenArgListEnd, ""},
		{59, gqlscan.TokenField, "baz_fuzz_taz5"},
		{60, gqlscan.TokenArgList, ""},
		{61, gqlscan.TokenArg, "variable"},
		{62, gqlscan.TokenVarRef, "variable"},
		{63, gqlscan.TokenArgListEnd, ""},
		{64, gqlscan.TokenField, "baz_fuzz_taz6"},
		{65, gqlscan.TokenArgList, ""},
		{66, gqlscan.TokenArg, "variable"},
		{67, gqlscan.TokenVarRef, "v"},
		{68, gqlscan.TokenArgListEnd, ""},
		{69, gqlscan.TokenField, "baz_fuzz_taz7"},
		{70, gqlscan.TokenArgList, ""},
		{71, gqlscan.TokenArg, "object"},
		{72, gqlscan.TokenObj, ""},
		{73, gqlscan.TokenObjField, "number0"},
		{74, gqlscan.TokenNum, "0"},
		{75, gqlscan.TokenObjField, "number1"},
		{76, gqlscan.TokenNum, "2"},
		{77, gqlscan.TokenObjField, "number2"},
		{78, gqlscan.TokenNum, "123456789.1234e2"},

		{79, gqlscan.TokenObjField, "arr0"},
		{80, gqlscan.TokenArr, ""},
		{81, gqlscan.TokenArr, ""},
		{82, gqlscan.TokenArrEnd, ""},
		{83, gqlscan.TokenArr, ""},
		{84, gqlscan.TokenObj, ""},
		{85, gqlscan.TokenObjField, "x"},
		{86, gqlscan.TokenNull, "null"},
		{87, gqlscan.TokenObjEnd, ""},
		{88, gqlscan.TokenArrEnd, ""},
		{89, gqlscan.TokenArrEnd, ""},

		{90, gqlscan.TokenObjEnd, ""},
		{91, gqlscan.TokenArgListEnd, ""},
		{92, gqlscan.TokenSelEnd, ""},
		{93, gqlscan.TokenSelEnd, ""},
		{94, gqlscan.TokenSelEnd, ""},

		// Mutation M
		{95, gqlscan.TokenDefMut, ""},
		{96, gqlscan.TokenMutName, "M"},
		{97, gqlscan.TokenVarList, ""},
		{98, gqlscan.TokenVarName, "variable"},
		{99, gqlscan.TokenVarTypeName, "Foo"},
		{100, gqlscan.TokenVarName, "v"},
		{101, gqlscan.TokenVarTypeArr, ""},
		{102, gqlscan.TokenVarTypeArr, ""},
		{103, gqlscan.TokenVarTypeName, "Bar"},
		{104, gqlscan.TokenVarTypeArrEnd, ""},
		{105, gqlscan.TokenVarTypeArrEnd, ""},
		{106, gqlscan.TokenVarListEnd, ""},
		{107, gqlscan.TokenSel, ""},
		{108, gqlscan.TokenField, "foo"},
		{109, gqlscan.TokenArgList, ""},
		{110, gqlscan.TokenArg, "x"},
		{111, gqlscan.TokenNull, "null"},
		{112, gqlscan.TokenArgListEnd, ""},
		{113, gqlscan.TokenSel, ""},
		{114, gqlscan.TokenField, "foo_bar"},
		{115, gqlscan.TokenSelEnd, ""},
		{116, gqlscan.TokenField, "bar"},
		{117, gqlscan.TokenField, "baz"},
		{118, gqlscan.TokenSel, ""},
		{119, gqlscan.TokenField, "baz_fuzz"},
		{120, gqlscan.TokenSel, ""},
		{121, gqlscan.TokenFragInline, "A"},
		{122, gqlscan.TokenSel, ""},
		{123, gqlscan.TokenField, "baz_fuzz_taz_A"},
		{124, gqlscan.TokenFragRef, "namedFragment1"},
		{125, gqlscan.TokenFragRef, "namedFragment2"},
		{126, gqlscan.TokenSelEnd, ""},
		{127, gqlscan.TokenFragInline, "B"},
		{128, gqlscan.TokenSel, ""},
		{129, gqlscan.TokenField, "baz_fuzz_taz_B"},
		{130, gqlscan.TokenSelEnd, ""},
		{131, gqlscan.TokenField, "baz_fuzz_taz1"},
		{132, gqlscan.TokenArgList, ""},
		{133, gqlscan.TokenArg, "bool"},
		{134, gqlscan.TokenFalse, "false"},
		{135, gqlscan.TokenArgListEnd, ""},
		{136, gqlscan.TokenField, "baz_fuzz_taz2"},
		{137, gqlscan.TokenArgList, ""},
		{138, gqlscan.TokenArg, "bool"},
		{139, gqlscan.TokenTrue, "true"},
		{140, gqlscan.TokenArgListEnd, ""},
		{141, gqlscan.TokenField, "baz_fuzz_taz3"},
		{142, gqlscan.TokenArgList, ""},
		{143, gqlscan.TokenArg, "string"},
		{144, gqlscan.TokenStr, "okay"},
		{145, gqlscan.TokenArgListEnd, ""},
		{146, gqlscan.TokenField, "baz_fuzz_taz4"},
		{147, gqlscan.TokenArgList, ""},
		{148, gqlscan.TokenArg, "array"},
		{149, gqlscan.TokenArr, ""},
		{150, gqlscan.TokenArrEnd, ""},
		{151, gqlscan.TokenArgListEnd, ""},
		{152, gqlscan.TokenField, "baz_fuzz_taz5"},
		{153, gqlscan.TokenArgList, ""},
		{154, gqlscan.TokenArg, "variable"},
		{155, gqlscan.TokenVarRef, "variable"},
		{156, gqlscan.TokenArgListEnd, ""},
		{157, gqlscan.TokenField, "baz_fuzz_taz6"},
		{158, gqlscan.TokenArgList, ""},
		{159, gqlscan.TokenArg, "variable"},
		{160, gqlscan.TokenVarRef, "v"},
		{161, gqlscan.TokenArgListEnd, ""},
		{162, gqlscan.TokenField, "baz_fuzz_taz7"},
		{163, gqlscan.TokenArgList, ""},
		{164, gqlscan.TokenArg, "object"},
		{165, gqlscan.TokenObj, ""},
		{166, gqlscan.TokenObjField, "number0"},
		{167, gqlscan.TokenNum, "0"},
		{168, gqlscan.TokenObjField, "number1"},
		{169, gqlscan.TokenNum, "2"},
		{170, gqlscan.TokenObjField, "number2"},
		{171, gqlscan.TokenNum, "123456789.1234e2"},

		{172, gqlscan.TokenObjField, "arr0"},
		{173, gqlscan.TokenArr, ""},
		{174, gqlscan.TokenArr, ""},
		{175, gqlscan.TokenArrEnd, ""},
		{176, gqlscan.TokenArr, ""},
		{177, gqlscan.TokenObj, ""},
		{178, gqlscan.TokenObjField, "x"},
		{179, gqlscan.TokenNull, "null"},
		{180, gqlscan.TokenObjEnd, ""},
		{181, gqlscan.TokenArrEnd, ""},
		{182, gqlscan.TokenArrEnd, ""},

		{183, gqlscan.TokenObjEnd, ""},
		{184, gqlscan.TokenArgListEnd, ""},
		{185, gqlscan.TokenSelEnd, ""},
		{186, gqlscan.TokenSelEnd, ""},
		{187, gqlscan.TokenSelEnd, ""},

		// Fragment f1
		{188, gqlscan.TokenDefFrag, ""},
		{189, gqlscan.TokenFragName, "f1"},
		{190, gqlscan.TokenFragTypeCond, "Query"},
		{191, gqlscan.TokenSel, ""},
		{192, gqlscan.TokenField, "todos"},
		{193, gqlscan.TokenSel, ""},
		{194, gqlscan.TokenFragRef, "f2"},
		{195, gqlscan.TokenSelEnd, ""},
		{196, gqlscan.TokenSelEnd, ""},

		// Query Todos
		{197, gqlscan.TokenDefQry, ""},
		{198, gqlscan.TokenQryName, "Todos"},
		{199, gqlscan.TokenSel, ""},
		{200, gqlscan.TokenFragRef, "f1"},
		{201, gqlscan.TokenSelEnd, ""},

		// Fragment f2
		{202, gqlscan.TokenDefFrag, ""},
		{203, gqlscan.TokenFragName, "f2"},
		{204, gqlscan.TokenFragTypeCond, "Todo"},
		{205, gqlscan.TokenSel, ""},
		{206, gqlscan.TokenField, "id"},
		{207, gqlscan.TokenField, "text"},
		{208, gqlscan.TokenArgList, ""},
		{209, gqlscan.TokenArg, "foo"},
		{210, gqlscan.TokenNum, "2"},
		{211, gqlscan.TokenArg, "bar"},
		{212, gqlscan.TokenStr, "ok"},
		{213, gqlscan.TokenArg, "baz"},
		{214, gqlscan.TokenNull, "null"},
		{215, gqlscan.TokenArgListEnd, ""},
		{216, gqlscan.TokenField, "done"},
		{217, gqlscan.TokenSelEnd, ""},
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
	{41, "{f#comment\n{f2}}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "f2"},
		{5, gqlscan.TokenSelEnd, ""},
		{6, gqlscan.TokenSelEnd, ""},
	}},

	// String escape
	{42, `{x(s:"\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\"`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{43, `{x(s:"\\")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\\`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{44, `{x(s:"\\\"")}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "s"},
		{5, gqlscan.TokenStr, `\\\"`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},

	{45, `{x(y:1e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `1e8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{46, `{x(y:0e8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{47, `{x(y:0e+8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e+8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{48, `{x(y:0e-8)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "y"},
		{5, gqlscan.TokenNum, `0e-8`},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{49, `mutation{x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "x"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{50, `mutation($x:T){x}`, []Expect{
		{0, gqlscan.TokenDefMut, ""},
		{1, gqlscan.TokenVarList, ""},
		{2, gqlscan.TokenVarName, "x"},
		{3, gqlscan.TokenVarTypeName, "T"},
		{4, gqlscan.TokenVarListEnd, ""},
		{5, gqlscan.TokenSel, ""},
		{6, gqlscan.TokenField, "x"},
		{7, gqlscan.TokenSelEnd, ""},
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
	{53, `{f(a:[0])}`, []Expect{
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
	{54, `query($v:T ! ){x(a:$v)}`, []Expect{
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
	{55, `query ($v: [ [ T ! ] ! ] ! ) {x(a:$v)}`, []Expect{
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
	{56, `{ bob : alice }`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenFieldAlias, "bob"},
		{3, gqlscan.TokenField, "alice"},
		{4, gqlscan.TokenSelEnd, ""},
	}},
	{57, `query # This is a test with many comments
	# sample comment text line
	{  #sample comment text line
		# sample comment text line
		a  #sample comment text line
		# sample comment text line
		{  #sample comment text line
			# sample comment text line
			b  #sample comment text line
			# sample comment text line
			(  #sample comment text line
				# sample comment text line
				x  #sample comment text line
				# sample comment text line
				:  #sample comment text line
				# sample comment text line
				1  #sample comment text line
			# sample comment text line
			)  #sample comment text line
			# sample comment text line
			{  #sample comment text line
				# sample comment text line
				c  #sample comment text line
				# sample comment text line
				d  #sample comment text line
			# sample comment text line
			}  #sample comment text line
		# sample comment text line
		}  #sample comment text line
	# sample comment text line
	}  #sample comment text line
	`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "a"},
		{3, gqlscan.TokenSel, ""},
		{4, gqlscan.TokenField, "b"},
		{5, gqlscan.TokenArgList, ""},
		{6, gqlscan.TokenArg, "x"},
		{7, gqlscan.TokenNum, "1"},
		{8, gqlscan.TokenArgListEnd, ""},
		{9, gqlscan.TokenSel, ""},
		{10, gqlscan.TokenField, "c"},
		{11, gqlscan.TokenField, "d"},
		{12, gqlscan.TokenSelEnd, ""},
		{13, gqlscan.TokenSelEnd, ""},
		{14, gqlscan.TokenSelEnd, ""},
	}},
	{58, `{f}
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
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenSelEnd, ""},
	}},
	{59, `{f(a:
		"\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000"
		b:123456789
	)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "a"},
		{5, gqlscan.TokenStr,
			`\b\t\r\n\f\/\"\u1234\u5678\u9abc\udefA\uBCDE\uF000`},
		{6, gqlscan.TokenArg, "b"},
		{7, gqlscan.TokenNum, "123456789"},
		{8, gqlscan.TokenArgListEnd, ""},
		{9, gqlscan.TokenSelEnd, ""},
	}},
	{60, "{f(a:" + string_2695b + ")}", []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "a"},
		{5, gqlscan.TokenStr, string_2695b[1 : len(string_2695b)-1]},
		{6, gqlscan.TokenArgListEnd, ""},
		{7, gqlscan.TokenSelEnd, ""},
	}},
	{61, `{f(
		a:""""""
		b:"""abc"""
		c:"""\n\t" """
		d:"""
			foo
				bar
		"""
	)}`, []Expect{
		{0, gqlscan.TokenDefQry, ""},
		{1, gqlscan.TokenSel, ""},
		{2, gqlscan.TokenField, "f"},
		{3, gqlscan.TokenArgList, ""},
		{4, gqlscan.TokenArg, "a"},
		{5, gqlscan.TokenStrBlock, ""},
		{6, gqlscan.TokenArg, "b"},
		{7, gqlscan.TokenStrBlock, "abc"},
		{8, gqlscan.TokenArg, "c"},
		{9, gqlscan.TokenStrBlock, `\n\t" `},
		{10, gqlscan.TokenArg, "d"},
		{11, gqlscan.TokenStrBlock,
			"\n\t\t\tfoo\n\t\t\t\tbar\n\t\t"},
		{12, gqlscan.TokenArgListEnd, ""},
		{13, gqlscan.TokenSelEnd, ""},
	}},
}

//go:embed string_2695b.txt
var string_2695b string

func TestScan(t *testing.T) {
	for ti, td := range testdata {
		t.Run("", func(t *testing.T) {
			t.Run("Scan", func(t *testing.T) {
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

			t.Run("ScanAll", func(t *testing.T) {
				require.Equal(t, ti, td.index)
				tName := t.Name()
				t.Log(tName)
				j := 0
				prevHead := 0
				err := gqlscan.ScanAll(
					[]byte(td.input),
					func(i *gqlscan.Iterator) {
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
		"error at index 1 ('1'): unexpected token; " +
			"expected field name or alias",
	},
	{7,
		"trailing closing curly bracket",
		"{f}}",
		"error at index 3 ('}'): unexpected token; expected definition",
	},
	{8,
		"query missing closing curly bracket",
		"{}",
		"error at index 1 ('}'): unexpected token; " +
			"expected field name or alias",
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
		"error at index 10 (')'): unexpected token; " +
			"expected field name or alias",
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
			"expected field name or alias",
	},
	{38,
		"unexpected EOF",
		"{foo ",
		"error at index 5: unexpected end of file; " +
			"expected field name or alias",
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
			"expected field name or alias",
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
	{106,
		"unexpected token",
		"{alias : alias2 : x}",
		"error at index 16 (':'): unexpected token; " +
			"expected field name or alias",
	},
	{107,
		"unexpected EOF",
		"{f:",
		"error at index 3: unexpected end of file; " +
			"expected field name",
	},
	{108,
		"unexpected EOF",
		"{f: ",
		"error at index 4: unexpected end of file; " +
			"expected field name",
	},
	{109,
		"invalid escape sequence",
		`{f(a:"\a")}`,
		"error at index 7 ('a'): unexpected token; " +
			"expected escaped sequence",
	},
	{110,
		"invalid escape sequence",
		`{f(a:"\u")}`,
		"error at index 8 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{111,
		"invalid escape sequence",
		`{f(a:"\u1")}`,
		"error at index 9 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{112,
		"invalid escape sequence",
		`{f(a:"\u12")}`,
		"error at index 10 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{113,
		"unexpected EOF",
		`{f(a:"\u`,
		"error at index 8: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{114,
		"unexpected EOF",
		`{f(a:"\u1`,
		"error at index 9: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{115,
		"unexpected EOF",
		`{f(a:"\u12`,
		"error at index 10: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{116,
		"unexpected EOF",
		`{f(a:"\u123`,
		"error at index 11: unexpected end of file; " +
			"expected escaped unicode sequence",
	},
	{117,
		"invalid escape sequence",
		`{f(a:"\u123")}`,
		"error at index 11 ('\"'): unexpected token; " +
			"expected escaped unicode sequence",
	},
	{118,
		"unexpected EOF",
		`{f(a:"""`,
		`error at index 8: unexpected end of file; ` +
			"expected end of block string",
	},
	{119,
		"unexpected EOF",
		`{f(a:""" `,
		"error at index 9: unexpected end of file; " +
			"expected end of block string",
	},
	{120,
		"control character in string",
		`{f(a:"0123456` + string(rune(0x00)) + `")}`,
		"error at index 13 (0x0): unexpected token; " +
			"expected end of string",
	},
	{121,
		"unexpected EOF",
		`{f #c`,
		"error at index 5: unexpected end of file; " +
			"expected selection, selection set or end of selection set",
	},
}

func TestScanErr(t *testing.T) {
	for ti, td := range testdataErr {
		t.Run("", func(t *testing.T) {
			t.Run("Scan", func(t *testing.T) {
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

			t.Run("ScanAll", func(t *testing.T) {
				require.Equal(t, ti, td.index)
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
		{13, gqlscan.TokenFieldAlias, "foo_alias", 1},
		{14, gqlscan.TokenField, "foo", 1},
		{15, gqlscan.TokenArgList, "", 1},
		{16, gqlscan.TokenArg, "x", 1},
		{17, gqlscan.TokenNull, "null", 1},
		{18, gqlscan.TokenArgListEnd, "", 1},

		{19, gqlscan.TokenSel, "", 1},
		{20, gqlscan.TokenFieldAlias, "foobar_alias", 2},
		{21, gqlscan.TokenField, "foo_bar", 2},
		{22, gqlscan.TokenSelEnd, "", 2},
		{23, gqlscan.TokenField, "bar", 1},
		{24, gqlscan.TokenField, "baz", 1},
		{25, gqlscan.TokenSel, "", 1},
		{26, gqlscan.TokenField, "baz_fuzz", 2},
		{27, gqlscan.TokenSel, "", 2},
		{28, gqlscan.TokenFragInline, "A", 3},
		{29, gqlscan.TokenSel, "", 3},
		{30, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{31, gqlscan.TokenFragRef, "namedFragment1", 4},
		{32, gqlscan.TokenFragRef, "namedFragment2", 4},
		{33, gqlscan.TokenSelEnd, "", 4},
		{34, gqlscan.TokenFragInline, "B", 3},
		{35, gqlscan.TokenSel, "", 3},
		{36, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{37, gqlscan.TokenSelEnd, "", 4},
		{38, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{39, gqlscan.TokenArgList, "", 3},
		{40, gqlscan.TokenArg, "bool", 3},
		{41, gqlscan.TokenFalse, "false", 3},
		{42, gqlscan.TokenArgListEnd, "", 3},
		{43, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{44, gqlscan.TokenArgList, "", 3},
		{45, gqlscan.TokenArg, "bool", 3},
		{46, gqlscan.TokenTrue, "true", 3},
		{47, gqlscan.TokenArgListEnd, "", 3},
		{48, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{49, gqlscan.TokenArgList, "", 3},
		{50, gqlscan.TokenArg, "string", 3},
		{51, gqlscan.TokenStr, "okay", 3},
		{52, gqlscan.TokenArgListEnd, "", 3},
		{53, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{54, gqlscan.TokenArgList, "", 3},
		{55, gqlscan.TokenArg, "array", 3},
		{56, gqlscan.TokenArr, "", 3},
		{57, gqlscan.TokenArrEnd, "", 3},
		{58, gqlscan.TokenArgListEnd, "", 3},
		{59, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{60, gqlscan.TokenArgList, "", 3},
		{61, gqlscan.TokenArg, "variable", 3},
		{62, gqlscan.TokenVarRef, "variable", 3},
		{63, gqlscan.TokenArgListEnd, "", 3},
		{64, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{65, gqlscan.TokenArgList, "", 3},
		{66, gqlscan.TokenArg, "variable", 3},
		{67, gqlscan.TokenVarRef, "v", 3},
		{68, gqlscan.TokenArgListEnd, "", 3},
		{69, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{70, gqlscan.TokenArgList, "", 3},
		{71, gqlscan.TokenArg, "object", 3},
		{72, gqlscan.TokenObj, "", 3},
		{73, gqlscan.TokenObjField, "number0", 3},
		{74, gqlscan.TokenNum, "0", 3},
		{75, gqlscan.TokenObjField, "number1", 3},
		{76, gqlscan.TokenNum, "2", 3},
		{77, gqlscan.TokenObjField, "number2", 3},
		{78, gqlscan.TokenNum, "123456789.1234e2", 3},

		{79, gqlscan.TokenObjField, "arr0", 3},
		{80, gqlscan.TokenArr, "", 3},
		{81, gqlscan.TokenArr, "", 3},
		{82, gqlscan.TokenArrEnd, "", 3},
		{83, gqlscan.TokenArr, "", 3},
		{84, gqlscan.TokenObj, "", 3},
		{85, gqlscan.TokenObjField, "x", 3},
		{86, gqlscan.TokenNull, "null", 3},
		{87, gqlscan.TokenObjEnd, "", 3},
		{88, gqlscan.TokenArrEnd, "", 3},
		{89, gqlscan.TokenArrEnd, "", 3},

		{90, gqlscan.TokenObjEnd, "", 3},
		{91, gqlscan.TokenArgListEnd, "", 3},
		{92, gqlscan.TokenSelEnd, "", 3},
		{93, gqlscan.TokenSelEnd, "", 2},
		{94, gqlscan.TokenSelEnd, "", 1},

		// Mutation M
		{95, gqlscan.TokenDefMut, "", 0},
		{96, gqlscan.TokenMutName, "M", 0},
		{97, gqlscan.TokenVarList, "", 0},
		{98, gqlscan.TokenVarName, "variable", 0},
		{99, gqlscan.TokenVarTypeName, "Foo", 0},
		{100, gqlscan.TokenVarName, "v", 0},
		{101, gqlscan.TokenVarTypeArr, "", 0},
		{102, gqlscan.TokenVarTypeArr, "", 0},
		{103, gqlscan.TokenVarTypeName, "Bar", 0},
		{104, gqlscan.TokenVarTypeArrEnd, "", 0},
		{105, gqlscan.TokenVarTypeArrEnd, "", 0},
		{106, gqlscan.TokenVarListEnd, "", 0},
		{107, gqlscan.TokenSel, "", 0},
		{108, gqlscan.TokenField, "foo", 1},
		{109, gqlscan.TokenArgList, "", 1},
		{110, gqlscan.TokenArg, "x", 1},
		{111, gqlscan.TokenNull, "null", 1},
		{112, gqlscan.TokenArgListEnd, "", 1},
		{113, gqlscan.TokenSel, "", 1},
		{114, gqlscan.TokenField, "foo_bar", 2},
		{115, gqlscan.TokenSelEnd, "", 2},
		{116, gqlscan.TokenField, "bar", 1},
		{117, gqlscan.TokenField, "baz", 1},
		{118, gqlscan.TokenSel, "", 1},
		{119, gqlscan.TokenField, "baz_fuzz", 2},
		{120, gqlscan.TokenSel, "", 2},
		{121, gqlscan.TokenFragInline, "A", 3},
		{122, gqlscan.TokenSel, "", 3},
		{123, gqlscan.TokenField, "baz_fuzz_taz_A", 4},
		{124, gqlscan.TokenFragRef, "namedFragment1", 4},
		{125, gqlscan.TokenFragRef, "namedFragment2", 4},
		{126, gqlscan.TokenSelEnd, "", 4},
		{127, gqlscan.TokenFragInline, "B", 3},
		{128, gqlscan.TokenSel, "", 3},
		{129, gqlscan.TokenField, "baz_fuzz_taz_B", 4},
		{130, gqlscan.TokenSelEnd, "", 4},
		{131, gqlscan.TokenField, "baz_fuzz_taz1", 3},
		{132, gqlscan.TokenArgList, "", 3},
		{133, gqlscan.TokenArg, "bool", 3},
		{134, gqlscan.TokenFalse, "false", 3},
		{135, gqlscan.TokenArgListEnd, "", 3},
		{136, gqlscan.TokenField, "baz_fuzz_taz2", 3},
		{137, gqlscan.TokenArgList, "", 3},
		{138, gqlscan.TokenArg, "bool", 3},
		{139, gqlscan.TokenTrue, "true", 3},
		{140, gqlscan.TokenArgListEnd, "", 3},
		{141, gqlscan.TokenField, "baz_fuzz_taz3", 3},
		{142, gqlscan.TokenArgList, "", 3},
		{143, gqlscan.TokenArg, "string", 3},
		{144, gqlscan.TokenStr, "okay", 3},
		{145, gqlscan.TokenArgListEnd, "", 3},
		{146, gqlscan.TokenField, "baz_fuzz_taz4", 3},
		{147, gqlscan.TokenArgList, "", 3},
		{148, gqlscan.TokenArg, "array", 3},
		{149, gqlscan.TokenArr, "", 3},
		{150, gqlscan.TokenArrEnd, "", 3},
		{151, gqlscan.TokenArgListEnd, "", 3},
		{152, gqlscan.TokenField, "baz_fuzz_taz5", 3},
		{153, gqlscan.TokenArgList, "", 3},
		{154, gqlscan.TokenArg, "variable", 3},
		{155, gqlscan.TokenVarRef, "variable", 3},
		{156, gqlscan.TokenArgListEnd, "", 3},
		{157, gqlscan.TokenField, "baz_fuzz_taz6", 3},
		{158, gqlscan.TokenArgList, "", 3},
		{159, gqlscan.TokenArg, "variable", 3},
		{160, gqlscan.TokenVarRef, "v", 3},
		{161, gqlscan.TokenArgListEnd, "", 3},
		{162, gqlscan.TokenField, "baz_fuzz_taz7", 3},
		{163, gqlscan.TokenArgList, "", 3},
		{164, gqlscan.TokenArg, "object", 3},
		{165, gqlscan.TokenObj, "", 3},
		{166, gqlscan.TokenObjField, "number0", 3},
		{167, gqlscan.TokenNum, "0", 3},
		{168, gqlscan.TokenObjField, "number1", 3},
		{169, gqlscan.TokenNum, "2", 3},
		{170, gqlscan.TokenObjField, "number2", 3},
		{171, gqlscan.TokenNum, "123456789.1234e2", 3},

		{172, gqlscan.TokenObjField, "arr0", 3},
		{173, gqlscan.TokenArr, "", 3},
		{174, gqlscan.TokenArr, "", 3},
		{175, gqlscan.TokenArrEnd, "", 3},
		{176, gqlscan.TokenArr, "", 3},
		{177, gqlscan.TokenObj, "", 3},
		{178, gqlscan.TokenObjField, "x", 3},
		{179, gqlscan.TokenNull, "null", 3},
		{180, gqlscan.TokenObjEnd, "", 3},
		{181, gqlscan.TokenArrEnd, "", 3},
		{182, gqlscan.TokenArrEnd, "", 3},

		{183, gqlscan.TokenObjEnd, "", 3},
		{184, gqlscan.TokenArgListEnd, "", 3},
		{185, gqlscan.TokenSelEnd, "", 3},
		{186, gqlscan.TokenSelEnd, "", 2},
		{187, gqlscan.TokenSelEnd, "", 1},

		// Fragment f1
		{188, gqlscan.TokenDefFrag, "", 0},
		{189, gqlscan.TokenFragName, "f1", 0},
		{190, gqlscan.TokenFragTypeCond, "Query", 0},
		{191, gqlscan.TokenSel, "", 0},
		{192, gqlscan.TokenField, "todos", 1},
		{193, gqlscan.TokenSel, "", 1},
		{194, gqlscan.TokenFragRef, "f2", 2},
		{195, gqlscan.TokenSelEnd, "", 2},
		{196, gqlscan.TokenSelEnd, "", 1},

		// Query Todos
		{197, gqlscan.TokenDefQry, "", 0},
		{198, gqlscan.TokenQryName, "Todos", 0},
		{199, gqlscan.TokenSel, "", 0},
		{200, gqlscan.TokenFragRef, "f1", 1},
		{201, gqlscan.TokenSelEnd, "", 1},

		// Fragment f2
		{202, gqlscan.TokenDefFrag, "", 0},
		{203, gqlscan.TokenFragName, "f2", 0},
		{204, gqlscan.TokenFragTypeCond, "Todo", 0},
		{205, gqlscan.TokenSel, "", 0},
		{206, gqlscan.TokenField, "id", 1},
		{207, gqlscan.TokenField, "text", 1},
		{208, gqlscan.TokenArgList, "", 1},
		{209, gqlscan.TokenArg, "foo", 1},
		{210, gqlscan.TokenNum, "2", 1},
		{211, gqlscan.TokenArg, "bar", 1},
		{212, gqlscan.TokenStr, "ok", 1},
		{213, gqlscan.TokenArg, "baz", 1},
		{214, gqlscan.TokenNull, "null", 1},
		{215, gqlscan.TokenArgListEnd, "", 1},
		{216, gqlscan.TokenField, "done", 1},
		{217, gqlscan.TokenSelEnd, "", 1},
	}

	t.Run("Scan", func(t *testing.T) {
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
	})

	t.Run("ScanAll", func(t *testing.T) {
		j := 0
		err := gqlscan.ScanAll(
			[]byte(input),
			func(i *gqlscan.Iterator) {
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
	index        int
	input        string
	tokenIndex   int
	buffer       []byte
	expectWrites [][]byte
}{
	{index: 0,
		input:        `{f(a:"0")}`,
		tokenIndex:   5,
		buffer:       nil,
		expectWrites: [][]byte{},
	},
	{index: 1,
		input:      `{f(a:"0")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("0"),
		},
	},
	{index: 2,
		input:      `{f(a:"01234567")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("01234567"),
		},
	},
	{index: 3,
		input:      `{f(a:"0123456789ab")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("01234567"),
			[]byte("89ab"),
		},
	},
	{index: 4,
		input:        `{f(a:"""""")}`,
		tokenIndex:   5,
		buffer:       make([]byte, 8),
		expectWrites: [][]byte{},
	},
	{index: 5,
		input:      `{f(a:"""abc""")}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("abc"),
		},
	},
	{index: 6,
		input:      `{f(a:"""\n\t" """)}`,
		tokenIndex: 5,
		buffer:     make([]byte, 10),
		expectWrites: [][]byte{
			[]byte(`\\n\\t\" `),
		},
	},
	{index: 7,
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
	{index: 8,
		input: `{f(a:"""
					first line
					 second\tline
				 """)}`,
		tokenIndex: 5,
		buffer:     make([]byte, 8),
		expectWrites: [][]byte{
			[]byte("first li"),
			[]byte("ne\n seco"),
			[]byte(`nd\\tlin`),
			[]byte("e"),
		},
	},
	{index: 9,
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
			[]byte(`\`), []byte(`"`), []byte("\n"),
			[]byte(`\`), []byte(`\`), []byte("\n"),
			[]byte(`\`), []byte(`"`),
			[]byte(`\`), []byte(`"`),
			[]byte(`\`), []byte(`"`),
		},
	},
	{index: 10,
		input:        `{f(a:"""\"""""")}`,
		tokenIndex:   5,
		buffer:       make([]byte, 6),
		expectWrites: [][]byte{[]byte(`\"\"\"`)},
	},
}

func TestScanInterpreted(t *testing.T) {
	for ti, td := range testdataBlockStrings {
		t.Run("", func(t *testing.T) {
			require := require.New(t)
			require.Equal(ti, td.index)
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
					"unexpected write at index %d", i,
				)
			}
		})
	}
}

func TestScanInterpretedStop(t *testing.T) {
	require := require.New(t)
	const s = "\n\t\t\tfirst line\\\"\"\"\n\t\t\t second\\tline\n\t\t"
	/*
		`{f(a:"""
			first line\"""
			 second\tline
		""")}`
	*/

	const q = `{f(a:"""` + s + `""")}`
	in := []byte(q)

	t.Run("Scan", func(t *testing.T) {
		c := -1
		err := gqlscan.Scan(in, func(i *gqlscan.Iterator) (err bool) {
			c++
			if c != 5 {
				return false
			}
			const bufLen = 8
			for stopAt := 0; stopAt < len(s)/bufLen; stopAt++ {
				buf, callCount := make([]byte, bufLen), 0
				i.ScanInterpreted(buf, func(buffer []byte) (stop bool) {
					callCount++
					return callCount > stopAt
				})
				require.Equal(stopAt+1, callCount)
			}
			return false
		})
		require.False(err.IsErr())
		require.Equal(q, string(in), "making sure input isn't mutated")
	})

	t.Run("ScanAll", func(t *testing.T) {
		c := -1
		err := gqlscan.ScanAll(in, func(i *gqlscan.Iterator) {
			c++
			if c != 5 {
				return
			}
			const bufLen = 8
			for stopAt := 0; stopAt < len(s)/bufLen; stopAt++ {
				buf, callCount := make([]byte, bufLen), 0
				i.ScanInterpreted(buf, func(buffer []byte) (stop bool) {
					callCount++
					return callCount > stopAt
				})
				require.Equal(stopAt+1, callCount)
			}
		})
		require.False(err.IsErr())
		require.Equal(q, string(in), "making sure input isn't mutated")
	})
}
