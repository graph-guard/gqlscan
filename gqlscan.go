package gqlscan

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Iterator is a GraphQL iterator for lexical analysis
type Iterator struct {
	// stack holds either TokenArr or TokenObj
	// and is reset for every argument.
	stack []Token

	expect Expect
	token  Token

	// str holds the original source
	str []byte

	// tail and head represent the iterator tail and head indexes
	tail, head int
	levelSel   int

	// errc holds the recent error code
	errc ErrorCode
}

func (i *Iterator) stackReset() {
	i.stack = i.stack[:0]
}

func (i *Iterator) stackLen() int {
	return len(i.stack)
}

// stackPush pushes a new token onto the stack.
func (i *Iterator) stackPush(t Token) {
	i.stack = append(i.stack, t)
}

// stackPop pops the top element of the stack returning it.
// Returns 0 if the stack was empty.
func (i *Iterator) stackPop() {
	if l := len(i.stack); l > 0 {
		i.stack = i.stack[:l-1]
	}
}

// stackTop returns the last pushed token.
func (i *Iterator) stackTop() Token {
	if l := len(i.stack); l > 0 {
		return i.stack[l-1]
	}
	return 0
}

var iteratorPool = sync.Pool{
	New: func() interface{} {
		return &Iterator{
			stack: make([]Token, 64),
		}
	},
}

func acquireIterator(str []byte) *Iterator {
	i := iteratorPool.Get().(*Iterator)
	i.stackReset()
	i.expect = ExpectDef
	i.tail, i.head = -1, 0
	i.str = str
	i.levelSel = 0
	i.errc = 0
	return i
}

// LevelSelect returns the current selector level.
func (i *Iterator) LevelSelect() int {
	return i.levelSel
}

// IndexHead returns the current head index.
func (i *Iterator) IndexHead() int {
	return i.head
}

// IndexTail returns the current tail index.
func (i *Iterator) IndexTail() int {
	return i.tail
}

// Token returns the current token type.
func (i *Iterator) Token() Token {
	return i.token
}

// Value returns the raw value of the current token.
// For TokenStrBlock it's the raw uninterpreted body of the string.
func (i *Iterator) Value() []byte {
	if i.tail < 0 {
		return nil
	}
	return i.str[i.tail:i.head]
}

// ScanInterpreted calls fn writing the interpreted part of
// the value to buffer as long as fn doesn't return true and
// the scan didn't reach the end of the interpreted value.
func (i *Iterator) ScanInterpreted(
	buffer []byte,
	fn func(buffer []byte) (stop bool),
) {
	if len(buffer) < 1 {
		return
	}
	if i.token != TokenStrBlock {
		offset := 0
		for offset < len(i.Value()) {
			b := buffer
			v := i.Value()[offset:]
			if len(v) > len(b) {
				v = v[:len(b)]
			} else {
				b = b[:len(v)]
			}
			copy(b, v)
			if fn(b) {
				return
			}
			offset += len(v)
		}
		return
	}

	// Determine block prefix
	shortestPrefixLen := 0
	v := i.Value()
	start, end := 0, len(v)
	{
		lastLineBreak := 0
		for i := range v {
			if v[i] == '\n' {
				lastLineBreak = i
			}
			if v[i] != '\n' && v[i] != ' ' && v[i] != '\t' {
				start = lastLineBreak
				break
			}
		}
	FIND_END:
		for i := len(v) - 1; i >= 0; i-- {
			if v[i] == '\n' {
				for ; i >= 0; i-- {
					if v[i] != '\n' && v[i] != ' ' && v[i] != '\t' {
						end = i + 1
						break FIND_END
					}
				}
			}
		}
		v = v[start:end]
	COUNT_LOOP:
		for len(v) > 0 {
			if v[0] == '\n' {
				// Count prefix length
				l := 0
				for v = v[1:]; ; l++ {
					if l >= len(v) {
						break COUNT_LOOP
					} else if v[l] != ' ' && v[l] != '\t' {
						v = v[l:]
						if shortestPrefixLen == 0 || shortestPrefixLen > l {
							shortestPrefixLen = l
						}
						break
					}
				}
				continue
			}
			v = v[1:]
		}
	}

	{
		v, bi := i.Value()[start:end], 0

		write := func(b byte) (stop bool) {
			buffer[bi] = b
			bi++
			if bi >= len(buffer) {
				bi = 0
				return fn(buffer)
			}
			return false
		}

		for i := 0; i < len(v); {
			if v[i] == '\n' {
				if i != 0 {
					if write(v[i]) {
						return
					}
				}
				// Ignore prefix
				if i+shortestPrefixLen+1 <= len(v) {
					i += shortestPrefixLen + 1
				}
				if v[i] == '\n' {
					continue
				}
			}
			if v[i] == '\\' && i+3 <= len(v) &&
				v[i+3] == '"' &&
				v[i+2] == '"' &&
				v[i+1] == '"' {
				if write('"') {
					return
				}
				if write('"') {
					return
				}
				if write('"') {
					return
				}
				i += 4
				continue
			}
			if write(v[i]) {
				return
			}
			i++
		}
		if b := buffer[:bi]; len(b) > 0 {
			if fn(buffer[:bi]) {
				return
			}
		}
	}
}

// skipSTNRC advances the iterator until the end of a sequence of spaces,
// tabs, line-feeds, carriage-returns and commas.
func (i *Iterator) skipSTNRC() {
	for i.head < len(i.str) {
		switch i.str[i.head] {
		case ',', ' ', '\n', '\t', '\r':
			i.head++
		default:
			return
		}
	}
}

// isHeadNameBody returns true if the current head is a
// legal name body character, otherwise returns false.
func (i *Iterator) isHeadNameBody() bool {
	return i.str[i.head] == '_' ||
		(i.str[i.head] >= '0' && i.str[i.head] <= '9') ||
		(i.str[i.head] >= 'a' && i.str[i.head] <= 'z') ||
		(i.str[i.head] >= 'A' && i.str[i.head] <= 'Z')
}

// isHeadSNTRC returns true if the current head is a space, line-feed,
// horizontal tab, carriage-return or comma, otherwise returns false.
func (i *Iterator) isHeadSNTRC() bool {
	return i.str[i.head] == ' ' ||
		i.str[i.head] == '\n' ||
		i.str[i.head] == '\r' ||
		i.str[i.head] == '\t' ||
		i.str[i.head] == ','
}

// isHeadNotNameStart returns true if the current head is not
// a legal name start character, otherwise returns false.
func (i *Iterator) isHeadNotNameStart() bool {
	return i.str[i.head] != '_' &&
		(i.str[i.head] < 'a' || i.str[i.head] > 'z') &&
		(i.str[i.head] < 'A' || i.str[i.head] > 'Z')
}

// isHeadCtrl returns true if the current head is a control character,
// otherwise returns false.
func (i *Iterator) isHeadCtrl() bool {
	return i.str[i.head] < 0x20
}

// isHeadDigit returns true if the current head is
// a number start character, otherwise returns false.
func (i *Iterator) isHeadDigit() bool {
	switch i.str[i.head] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isHeadDigit returns true if the current head is
// a number start character, otherwise returns false.
func (i *Iterator) isHeadHexDigit() bool {
	switch i.str[i.head] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f',
		'A', 'B', 'C', 'D', 'E', 'F':
		return true
	}
	return false
}

// isHeadNumEnd returns true if the current head is a space, line-feed,
// horizontal tab, carriage-return, comma, right parenthesis or
// right curly brace, otherwise returns false.
func (i *Iterator) isHeadNumEnd() bool {
	switch i.str[i.head] {
	case ' ', '\t', '\r', '\n', ',', ')', '}', ']', '#':
		return true
	}
	return false
}

// isHeadKeywordQuery returns true if the current head equals 'query'.
func (i *Iterator) isHeadKeywordQuery() bool {
	return i.head+4 < len(i.str) &&
		i.str[i.head+4] == 'y' &&
		i.str[i.head+3] == 'r' &&
		i.str[i.head+2] == 'e' &&
		i.str[i.head+1] == 'u' &&
		i.str[i.head] == 'q'
}

// isHeadKeywordMutation returns true if the current head equals 'mutation'.
func (i *Iterator) isHeadKeywordMutation() bool {
	return i.head+7 < len(i.str) &&
		i.str[i.head+7] == 'n' &&
		i.str[i.head+6] == 'o' &&
		i.str[i.head+5] == 'i' &&
		i.str[i.head+4] == 't' &&
		i.str[i.head+3] == 'a' &&
		i.str[i.head+2] == 't' &&
		i.str[i.head+1] == 'u' &&
		i.str[i.head] == 'm'
}

// isHeadKeywordFragment returns true if the current head equals 'fragment'.
func (i *Iterator) isHeadKeywordFragment() bool {
	return i.head+7 < len(i.str) &&
		i.str[i.head+7] == 't' &&
		i.str[i.head+6] == 'n' &&
		i.str[i.head+5] == 'e' &&
		i.str[i.head+4] == 'm' &&
		i.str[i.head+3] == 'g' &&
		i.str[i.head+2] == 'a' &&
		i.str[i.head+1] == 'r' &&
		i.str[i.head] == 'f'
}

// Expect defines an expectation
type Expect int

// Expectations
const (
	_ Expect = iota
	ExpectVal
	ExpectDef
	ExpectQryName
	ExpectMutName
	ExpectSelSet
	ExpectArgName
	ExpectEscapedSequence
	ExpectEscapedUnicodeSequence
	ExpectEndOfString
	ExpectEndOfBlockString
	ExpectColumnAfterArg
	ExpectFieldNameOrAlias
	ExpectFieldName
	ExpectSel
	ExpectVarName
	ExpectVarRefName
	ExpectVarType
	ExpectColumnAfterVar
	ExpectObjFieldName
	ExpectColObjFieldName
	ExpectFragTypeCond
	ExpectFragKeywordOn
	ExpectFragName
	ExpectFrag
	ExpectFragRef
	ExpectFragInlined
	ExpectAfterFieldName
	ExpectAfterSelection
	ExpectAfterValue
	ExpectAfterArgList
	ExpectAfterKeywordQuery
	ExpectAfterKeywordMutation
	ExpectAfterVarType
	ExpectAfterVarTypeName
)

func (e Expect) String() string {
	switch e {
	case ExpectDef:
		return "definition"
	case ExpectQryName:
		return "query name"
	case ExpectMutName:
		return "mutation name"
	case ExpectVal:
		return "value"
	case ExpectSelSet:
		return "selection set"
	case ExpectArgName:
		return "argument name"
	case ExpectEscapedSequence:
		return "escaped sequence"
	case ExpectEscapedUnicodeSequence:
		return "escaped unicode sequence"
	case ExpectEndOfString:
		return "end of string"
	case ExpectEndOfBlockString:
		return "end of block string"
	case ExpectColumnAfterArg:
		return "column after argument name"
	case ExpectFieldNameOrAlias:
		return "field name or alias"
	case ExpectFieldName:
		return "field name"
	case ExpectSel:
		return "selection"
	case ExpectVarName:
		return "variable name"
	case ExpectVarRefName:
		return "referenced variable name"
	case ExpectVarType:
		return "variable type"
	case ExpectColumnAfterVar:
		return "column after variable name"
	case ExpectObjFieldName:
		return "object field name"
	case ExpectColObjFieldName:
		return "column after object field name"
	case ExpectFragTypeCond:
		return "fragment type condition"
	case ExpectFragKeywordOn:
		return "keyword 'on'"
	case ExpectFragName:
		return "fragment name"
	case ExpectFrag:
		return "fragment"
	case ExpectFragRef:
		return "fragment reference"
	case ExpectFragInlined:
		return "inlined fragment"
	case ExpectAfterFieldName:
		return "selection, selection set or end of selection set"
	case ExpectAfterSelection:
		return "selection or end of selection set"
	case ExpectAfterValue:
		return "argument list closure or argument"
	case ExpectAfterArgList:
		return "selection set or selection"
	case ExpectAfterKeywordQuery:
		return "variable list or selection set"
	case ExpectAfterKeywordMutation:
		return "variable list or selection set"
	case ExpectAfterVarType:
		return "variable list closure or variable name"
	case ExpectAfterVarTypeName:
		return "variable list closure or variable name"
	}
	return ""
}

// Token defines the type of a token.
type Token int

// Token types
const (
	_ Token = iota
	TokenDefQry
	TokenDefMut
	TokenDefFrag
	TokenQryName
	TokenMutName
	TokenVarList
	TokenVarListEnd
	TokenArgList
	TokenArgListEnd
	TokenSel
	TokenSelEnd
	TokenFragTypeCond
	TokenFragName
	TokenFragInline
	TokenFragRef
	TokenFieldAlias
	TokenField
	TokenArg
	TokenArr
	TokenArrEnd
	TokenStr
	TokenStrBlock
	TokenNum
	TokenTrue
	TokenFalse
	TokenNull
	TokenVarName
	TokenVarTypeName
	TokenVarTypeArr
	TokenVarTypeArrEnd
	TokenVarTypeNotNull
	TokenVarRef
	TokenObj
	TokenObjEnd
	TokenObjField
)

func (t Token) String() string {
	switch t {
	case TokenDefQry:
		return "query definition"
	case TokenDefMut:
		return "mutation definition"
	case TokenDefFrag:
		return "fragment definition"
	case TokenQryName:
		return "query name"
	case TokenMutName:
		return "mutation name"
	case TokenVarList:
		return "variable list"
	case TokenVarListEnd:
		return "variable list end"
	case TokenArgList:
		return "argument list"
	case TokenArgListEnd:
		return "argument list end"
	case TokenSel:
		return "selection set"
	case TokenSelEnd:
		return "selection set end"
	case TokenFragTypeCond:
		return "fragment type condition"
	case TokenFragName:
		return "fragment name"
	case TokenFragInline:
		return "fragment inline"
	case TokenFragRef:
		return "fragment reference"
	case TokenFieldAlias:
		return "field alias"
	case TokenField:
		return "field"
	case TokenArg:
		return "argument"
	case TokenArr:
		return "array"
	case TokenArrEnd:
		return "array end"
	case TokenStr:
		return "string"
	case TokenStrBlock:
		return "block string"
	case TokenNum:
		return "number"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenNull:
		return "null"
	case TokenVarName:
		return "variable name"
	case TokenVarTypeName:
		return "variable type name"
	case TokenVarTypeArr:
		return "variable array type"
	case TokenVarTypeArrEnd:
		return "variable array type end"
	case TokenVarTypeNotNull:
		return "variable type not null"
	case TokenVarRef:
		return "variable reference"
	case TokenObj:
		return "object"
	case TokenObjEnd:
		return "object end"
	case TokenObjField:
		return "object field"
	}
	return ""
}

// ErrorCode defines the type of an error.
type ErrorCode int

const (
	_ ErrorCode = iota
	ErrCallbackFn
	ErrUnexpToken
	ErrUnexpEOF
	ErrInvalNum
	ErrInvalVal
	ErrInvalType
)

// Error is a GraphQL lexical scan error.
type Error struct {
	Index       int
	AtIndex     rune
	Code        ErrorCode
	Expectation Expect
}

// IsErr returns true if there is an error, otherwise returns false.
func (e Error) IsErr() bool {
	return e.Code != 0
}

func (e Error) Error() string {
	if e.Code == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("error at index ")
	b.WriteString(strconv.Itoa(e.Index))
	if e.Code != ErrUnexpEOF {
		if e.AtIndex < 0x20 {
			b.WriteString(" (")
			b.WriteString(fmt.Sprintf("0x%x", e.AtIndex))
			b.WriteString(")")
		} else {
			b.WriteString(" ('")
			b.WriteRune(e.AtIndex)
			b.WriteString("')")
		}
	}
	switch e.Code {
	case ErrCallbackFn:
		b.WriteString(": callback function returned error")
	case ErrUnexpToken:
		b.WriteString(": unexpected token")
	case ErrInvalNum:
		b.WriteString(": invalid number value")
	case ErrInvalVal:
		b.WriteString(": invalid value")
	case ErrInvalType:
		b.WriteString(": invalid type")
	case ErrUnexpEOF:
		b.WriteString(": unexpected end of file")
	}
	if e.Expectation != 0 {
		b.WriteString("; expected ")
		b.WriteString(e.Expectation.String())
	}
	return b.String()
}
