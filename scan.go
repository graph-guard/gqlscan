package gqlscan

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
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

// Value returns value of the current token.
// For TokenQry, TokenMut, TokenSel, TokenArr, TokenTrue,
// TokenFalse and TokenNull it the textual representation of the token.
// For TokenField it's the name of the field.
// For TokenStr it's the body of the string.
// For TokenNum it's the number.
func (i *Iterator) Value() []byte {
	if i.tail < 0 {
		return nil
	}
	return i.str[i.tail:i.head]
}

// skipSTNRC advances the iterator until the end of a sequence of spaces,
// tabs, line-feeds, carriage-returns and commas.
func (i *Iterator) skipSTNRC() {
	if i.head >= len(i.str) {
		return
	}
	for ; i.head < len(i.str) && (i.str[i.head] == ',' ||
		i.str[i.head] == ' ' ||
		i.str[i.head] == '\n' ||
		i.str[i.head] == '\t' ||
		i.str[i.head] == '\r'); i.head++ {
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

// isHeadNumberStart returns true if the current head is
// a number start character, otherwise returns false.
func (i *Iterator) isHeadNumberStart() bool {
	return i.str[i.head] == '+' ||
		i.str[i.head] == '-' ||
		(i.str[i.head] >= '0' && i.str[i.head] <= '9')
}

// isHeadDigit returns true if the current head is
// a number start character, otherwise returns false.
func (i *Iterator) isHeadDigit() bool {
	return i.str[i.head] >= '0' && i.str[i.head] <= '9'
}

// isHeadNumEnd returns true if the current head is a space, line-feed,
// horizontal tab, carriage-return, comma, right parenthesis or
// right curly brace, otherwise returns false.
func (i *Iterator) isHeadNumEnd() bool {
	return i.str[i.head] == ' ' ||
		i.str[i.head] == '\t' ||
		i.str[i.head] == '\r' ||
		i.str[i.head] == '\n' ||
		i.str[i.head] == ',' ||
		i.str[i.head] == ')' ||
		i.str[i.head] == '}' ||
		i.str[i.head] == ']' ||
		i.str[i.head] == '#'
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

// Scan calls fn for every token it scans in str.
// If fn returns true then an error with code ErrCallbackFn is returned.
// *Iterator passed to fn should never be aliased and used after Scan returns.
func Scan(str []byte, fn func(*Iterator) (err bool)) Error {
	i := acquireIterator(str)
	defer iteratorPool.Put(i)

	var typeArrLvl int

	i.skipSTNRC()

	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		i.expect = ExpectDef
		goto ERROR
	}

DEFINITION:
	if i.head >= len(i.str) {
		goto DEFINITION_END
	} else if i.str[i.head] == '#' {
		i.expect = ExpectDef
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token = TokenDefQry
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.isHeadKeywordQuery() {
		// Query
		i.token = TokenDefQry
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head += len("query")
		i.expect = ExpectAfterKeywordQuery
		goto AFTER_KEYWORD_QUERY
	} else if i.isHeadKeywordMutation() {
		// Mutation
		i.token = TokenDefMut
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head += len("mutation")
		i.expect = ExpectAfterKeywordMutation
		goto AFTER_KEYWORD_MUTATION
	} else if i.isHeadKeywordFragment() {
		// Fragment
		i.tail = -1
		i.token = TokenDefFrag
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head += len("fragment")
		i.expect = ExpectFragName
		goto AFTER_KEYWORD_FRAGMENT
	}

	i.errc = ErrUnexpToken
	i.expect = ExpectDef
	goto ERROR

AFTER_KEYWORD_QUERY:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '(' {
		// Query variable list
		i.head++
		i.expect = ExpectVarName
		goto QUERY_VAR
	}
	i.expect = ExpectQryName
	goto NAME

AFTER_KEYWORD_MUTATION:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '(' {
		// Mutation variable list
		i.head++
		i.expect = ExpectVarName
		goto QUERY_VAR
	}
	i.expect = ExpectMutName
	goto NAME

AFTER_KEYWORD_FRAGMENT:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	goto NAME

QUERY_VAR:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}

	// Variable name
	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] != '$' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	goto NAME

AFTER_VAR_TYPE:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if typeArrLvl != 0 {
		i.head--
		i.errc = ErrInvalType
		i.expect = ExpectVarType
		goto ERROR
	} else if i.str[i.head] != ')' {
		i.expect = ExpectAfterVarType
		goto QUERY_VAR
	}
	// End of query variable list
	i.head++
	i.expect = ExpectSelSet
	goto SELECTION_SET

SELECTION_SET:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != '{' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.tail = -1
	i.token = TokenSel
	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}
	i.levelSel++
	i.head++
	i.expect = ExpectSel
	goto SELECTION

AFTER_SELECTION:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '}' {
		i.tail = -1
		i.token = TokenSelEnd
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.levelSel--
		i.head++
		i.skipSTNRC()
		if i.levelSel < 1 {
			goto DEFINITION_END
		}
		goto AFTER_SELECTION
	}
	i.expect = ExpectSel
	goto SELECTION

VALUE:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		// Object begin
		i.tail = -1
		// Callback for argument
		i.token = TokenObj
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.stackPush(TokenObj)
		i.head++
		i.skipSTNRC()

		i.expect = ExpectObjFieldName
		goto NAME
	} else if i.str[i.head] == '[' {
		i.tail = -1
		// Callback for argument
		i.token = TokenArr
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		i.skipSTNRC()

		// Lookahead
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] == ']' {
			i.token = TokenArrEnd
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.expect = ExpectAfterValue
			goto AFTER_VALUE_COMMENT
		}
		i.stackPush(TokenArr)
		i.expect = ExpectAfterValue
		goto AFTER_VALUE_COMMENT

	} else if i.str[i.head] == '"' {
		// String value
		i.tail = i.head + 1
		escaped := false
		for i.head++; i.head < len(i.str); i.head++ {
			if i.isHeadCtrl() {
				i.errc = ErrUnexpToken
				i.expect = ExpectEndOfString
				goto ERROR
			}
			if i.str[i.head] == '"' {
				if escaped {
					escaped = false
				} else {
					goto AFTER_STR_VAL
				}
			} else if i.str[i.head] == '\\' {
				escaped = !escaped
			}
		}
		i.errc = ErrUnexpEOF
		i.expect = ExpectEndOfString
		goto ERROR

	AFTER_STR_VAL:
		// Callback for argument
		i.token = TokenStr
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		// Advance head index to include the closing double-quotes
		i.head++
	} else if i.str[i.head] == '$' {
		// Variable reference
		i.head++

		// Variable name
		i.expect = ExpectVarRefName
		goto NAME

	} else if i.str[i.head] == 'n' {
		// Null
		if i.head+3 >= len(i.str) {
			i.head += len(i.str) - i.head
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] != 'n' ||
			i.str[i.head+1] != 'u' ||
			i.str[i.head+2] != 'l' ||
			i.str[i.head+3] != 'l' {
			i.errc = ErrInvalVal
			i.expect = ExpectVal
			goto ERROR
		}
		i.tail = i.head
		i.head += len("null")

		// Callback for argument
		i.token = TokenNull
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
	} else if i.str[i.head] == 't' {
		// Boolean true
		if i.head+3 >= len(i.str) {
			i.head += len(i.str) - i.head
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] != 't' ||
			i.str[i.head+1] != 'r' ||
			i.str[i.head+2] != 'u' ||
			i.str[i.head+3] != 'e' {
			i.errc = ErrInvalVal
			i.expect = ExpectVal
			goto ERROR
		}
		i.tail = i.head
		i.head += len("true")

		// Callback for argument
		i.token = TokenTrue
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
	} else if i.str[i.head] == 'f' {
		// Boolean false
		if i.head+4 >= len(i.str) {
			i.head += len(i.str) - i.head
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] != 'f' ||
			i.str[i.head+1] != 'a' ||
			i.str[i.head+2] != 'l' ||
			i.str[i.head+3] != 's' ||
			i.str[i.head+4] != 'e' {
			i.errc = ErrInvalVal
			i.expect = ExpectVal
			goto ERROR
		}
		i.tail = i.head
		i.head += len("false")

		// Callback for argument
		i.token = TokenFalse
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
	} else if i.isHeadNumberStart() {
		// Number
		i.tail = i.head

		var s int

		switch i.str[i.head] {
		case '-':
			// Signed
			i.head++
			if i.head >= len(i.str) {
				// Expected at least one digit
				i.errc = ErrUnexpEOF
				i.expect = ExpectVal
				goto ERROR
			}
		case '0':
			// Leading zero
			i.head++
			if len(i.str) > i.head {
				if i.str[i.head] == '.' {
					i.head++
					goto FRACTION
				} else if i.str[i.head] == 'e' || i.str[i.head] == 'E' {
					i.head++
					goto EXPONENT_SIGN
				} else if i.isHeadNumEnd() {
					goto ON_NUM_VAL
				} else {
					i.errc = ErrInvalNum
					i.expect = ExpectVal
					goto ERROR
				}
			}
		}

		// Integer
		for s = i.head; i.head < len(i.str); i.head++ {
			if i.isHeadDigit() {
				continue
			} else if i.str[i.head] == '.' {
				i.head++
				goto FRACTION
			} else if i.isHeadNumEnd() {
				if i.head == s {
					// Expected at least one digit
					i.errc = ErrInvalNum
					i.expect = ExpectVal
					goto ERROR
				}
				// Integer
				goto ON_NUM_VAL
			} else if i.str[i.head] == 'e' || i.str[i.head] == 'E' {
				i.head++
				goto EXPONENT_SIGN
			}

			// Unexpected rune
			i.errc = ErrInvalNum
			i.expect = ExpectVal
			goto ERROR
		}

		if i.head >= len(i.str) {
			// Integer without exponent
			goto ON_NUM_VAL
		}
		// Continue to fraction

	FRACTION:

		for s = i.head; i.head < len(i.str); i.head++ {
			if i.isHeadDigit() {
				continue
			} else if i.isHeadNumEnd() {
				if i.head == s {
					// Expected at least one digit
					i.errc = ErrInvalNum
					i.expect = ExpectVal
					goto ERROR
				}
				// Number with fraction
				goto ON_NUM_VAL
			} else if i.str[i.head] == 'e' || i.str[i.head] == 'E' {
				i.head++
				goto EXPONENT_SIGN
			}

			// Unexpected rune
			i.errc = ErrInvalNum
			i.expect = ExpectVal
			goto ERROR
		}
		if s == i.head {
			// Unexpected end of number
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		}

		if i.head >= len(i.str) {
			// Number (with fraction but) without exponent
			goto ON_NUM_VAL
		}

	EXPONENT_SIGN:
		if i.head >= len(i.str) {
			// Missing exponent value
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		}
		if i.str[i.head] == '-' || i.str[i.head] == '+' {
			i.head++
		}
		for s = i.head; i.head < len(i.str); i.head++ {
			if i.isHeadDigit() {
				continue
			} else if i.isHeadNumEnd() {
				if i.head == s {
					// Expected at least one digit
					i.errc = ErrInvalNum
					i.expect = ExpectVal
					goto ERROR
				}
				// Number with (fraction and) exponent
				goto ON_NUM_VAL
			}
			break
		}
		// Unexpected rune
		i.errc = ErrInvalNum
		i.expect = ExpectVal
		goto ERROR

	ON_NUM_VAL:
		// Callback for argument
		i.token = TokenNum
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
	} else {
		// Invalid value
		i.errc = ErrInvalVal
		i.expect = ExpectVal
		goto ERROR
	}
	i.expect = ExpectAfterValue
	goto AFTER_VALUE_COMMENT

AFTER_VALUE_COMMENT:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	if t := i.stackTop(); t == TokenObj {
		if i.str[i.head] == '}' {
			i.tail = -1
			i.stackPop()

			// Callback for end of object
			i.token = TokenObjEnd
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			i.head++
			i.skipSTNRC()
			if i.stackLen() > 0 {
				i.expect = ExpectAfterValue
				goto AFTER_VALUE_COMMENT
			}
		} else {
			// Proceed to next field in the object
			i.expect = ExpectObjFieldName
			goto NAME
		}
	} else if t == TokenArr {
		if i.str[i.head] == ']' {
			i.tail = -1
			i.stackPop()

			// Callback for end of array
			i.token = TokenArrEnd
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			if i.stackLen() > 0 {
				i.expect = ExpectAfterValue
				goto AFTER_VALUE_COMMENT
			}
		} else {
			// Proceed to next value in the array
			goto VALUE
		}
	}
	if i.str[i.head] == ')' {
		// End of argument list
		i.head++
		i.expect = ExpectAfterArgList
		goto AFTER_ARG_LIST
	}
	// Proceed to the next argument
	i.expect = ExpectArgName
	goto NAME

AFTER_ARG_LIST:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '}' {
		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION
	}
	i.expect = ExpectSel
	goto SELECTION

SELECTION:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		i.expect = ExpectSel
		goto ERROR
	} else if i.str[i.head] == '#' {
		i.expect = ExpectSel
		goto COMMENT
	} else if i.str[i.head] != '.' {
		// Field selection
		i.expect = ExpectFieldName
		goto NAME
	}

	i.expect = ExpectFrag
	if i.head+2 >= len(i.str) {
		i.errc = ErrUnexpEOF
		if i.head+1 >= len(i.str) {
			i.head++
		} else {
			i.head += 2
		}
		goto ERROR
	} else if i.str[i.head+2] != '.' ||
		i.str[i.head+1] != '.' {
		i.errc = ErrUnexpToken
		if i.str[i.head+1] != '.' {
			i.head += 1
		} else if i.str[i.head+2] != '.' {
			i.head += 2
		}
		goto ERROR
	}

	i.head += len("...")
	goto FRAGMENT

FRAGMENT:
	i.skipSTNRC()
	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head+1] == 'n' &&
		i.str[i.head] == 'o' {
		if i.head+2 >= len(i.str) {
			i.head = len(i.str)
			i.errc = ErrUnexpEOF
			goto ERROR
		} else if i.str[i.head+2] == ' ' ||
			i.str[i.head+2] == '\n' ||
			i.str[i.head+2] == '\r' ||
			i.str[i.head+2] == '\t' ||
			i.str[i.head+2] == ',' ||
			i.str[i.head+2] == '#' {
			// ... on Type {
			i.head += len("on ")
			i.skipSTNRC()
			i.expect = ExpectFragInlined
			goto NAME
		}
	}
	// ...fragmentName
	i.expect = ExpectFragRef
	goto NAME

AFTER_DECL_VAR_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.expect = ExpectVarType
	goto VAR_TYPE

VAR_TYPE:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '[' {
		i.tail = -1
		i.token = TokenVarTypeArr
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		typeArrLvl++
		goto VAR_TYPE
	}
	i.expect = ExpectVarType
	goto NAME

NAME:
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	i.tail = i.head
	if i.isHeadNotNameStart() {
		i.errc = ErrUnexpToken
		goto ERROR
	}

	for ; i.head < len(i.str); i.head++ {
		if i.isHeadNameBody() {
			continue
		} else if i.isHeadSNTRC() {
			goto AFTER_NAME
		} else if i.isHeadCtrl() {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		break
	}

	goto AFTER_NAME

COLUMN_AFTER_ARG_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.stackReset()
	i.expect = ExpectVal
	goto VALUE

ARG_LIST:
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	goto NAME

AFTER_VAR_TYPE_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == ']' {
		if typeArrLvl < 1 {
			i.errc = ErrUnexpToken
			i.expect = ExpectVarName
			goto ERROR
		}
		i.tail = -1
		i.token = TokenVarTypeArrEnd
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		typeArrLvl--

		if typeArrLvl > 0 {
			goto AFTER_VAR_TYPE_NAME
		}
	}
	i.expect = ExpectAfterVarType
	goto AFTER_VAR_TYPE

AFTER_NAME:
	switch i.expect {
	case ExpectFieldName:
		// Callback for field name
		i.token = TokenField
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		// Lookahead
		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectFieldName
			goto ERROR
		} else if i.str[i.head] == '(' {
			// Argument list
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		} else if i.str[i.head] == '{' {
			// Field selector expands without arguments
			i.expect = ExpectSelSet
			goto SELECTION_SET
		}
		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION

	case ExpectArgName:
		// Callback for argument name
		i.token = TokenArg
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.skipSTNRC()
		i.expect = ExpectColumnAfterArg
		goto COLUMN_AFTER_ARG_NAME

	case ExpectObjFieldName:
		// Callback for object field
		i.token = TokenObjField
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectColObjFieldName
			goto ERROR
		} else if i.str[i.head] != ':' {
			i.errc = ErrUnexpToken
			i.expect = ExpectColObjFieldName
			goto ERROR
		}
		i.head++
		i.skipSTNRC()
		i.expect = ExpectVal
		goto VALUE

	case ExpectVarRefName:
		i.token = TokenVarRef
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectAfterValue
		goto AFTER_VALUE_COMMENT

	case ExpectVarType:
		i.token = TokenVarTypeName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectAfterVarTypeName
		goto AFTER_VAR_TYPE_NAME

	case ExpectVarName, ExpectAfterVarType:
		i.token = TokenVarName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectColumnAfterVar
		goto AFTER_DECL_VAR_NAME

	case ExpectQryName:
		i.token = TokenQryName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.skipSTNRC()

		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectSelSet
			goto ERROR
		} else if i.str[i.head] == '{' {
			i.expect = ExpectSelSet
			goto SELECTION_SET
		} else if i.str[i.head] == '(' {
			// Query variable list
			i.head++
			i.expect = ExpectVarName
			goto QUERY_VAR
		}
		i.errc = ErrUnexpToken
		i.expect = ExpectSelSet
		goto ERROR

	case ExpectMutName:
		i.token = TokenMutName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.skipSTNRC()

		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectSelSet
			goto ERROR
		} else if i.str[i.head] == '{' {
			i.expect = ExpectSelSet
			goto SELECTION_SET
		} else if i.str[i.head] == '(' {
			// Mutation variable list
			i.head++
			i.expect = ExpectVarName
			goto QUERY_VAR
		}
		i.errc = ErrUnexpToken
		i.expect = ExpectSelSet
		goto ERROR

	case ExpectFragInlined:
		i.token = TokenFragInline
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectSelSet
		goto SELECTION_SET

	case ExpectFragRef:
		i.token = TokenFragRef
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION

	case ExpectFragName:
		i.token = TokenFragName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectFragKeywordOn
		goto FRAG_KEYWORD_ON

	case ExpectFragTypeCond:
		i.token = TokenFragTypeCond
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectSelSet
		goto SELECTION_SET
	default:
		// This line should never be executed!
		// The panic is triggered only if we forgot to handle an expectation.
		panic(fmt.Errorf("unhandled expectation: %q", i.expect))
	}

FRAG_KEYWORD_ON:
	i.skipSTNRC()
	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head+1] != 'n' ||
		i.str[i.head] != 'o' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head += len("on")
	i.expect = ExpectFragTypeCond
	goto FRAG_TYPE_COND

FRAG_TYPE_COND:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	goto NAME

COMMENT:
	i.head++
	for i.tail = i.head; i.str[i.head] != '\n'; i.head++ {
	}
	i.tail = -1
	i.skipSTNRC()
	switch i.expect {
	case ExpectDef:
		goto DEFINITION
	case ExpectSelSet:
		goto SELECTION_SET
	case ExpectSel:
		goto SELECTION
	case ExpectAfterSelection:
		goto AFTER_SELECTION
	case ExpectVarName:
		goto QUERY_VAR
	case ExpectArgName:
		goto ARG_LIST
	case ExpectColumnAfterArg:
		goto COLUMN_AFTER_ARG_NAME
	case ExpectVal:
		goto VALUE
	case ExpectAfterValue:
		goto AFTER_VALUE_COMMENT
	case ExpectAfterArgList:
		goto AFTER_ARG_LIST
	case ExpectAfterKeywordQuery:
		goto AFTER_KEYWORD_QUERY
	case ExpectAfterKeywordMutation:
		goto AFTER_KEYWORD_MUTATION
	case ExpectFragName:
		goto AFTER_KEYWORD_FRAGMENT
	case ExpectFragKeywordOn:
		goto FRAG_KEYWORD_ON
	case ExpectFragTypeCond:
		goto FRAG_TYPE_COND
	case ExpectFrag:
		goto FRAGMENT
	case ExpectColumnAfterVar:
		goto AFTER_DECL_VAR_NAME
	case ExpectVarType:
		goto VAR_TYPE
	case ExpectAfterVarType:
		goto AFTER_VAR_TYPE
	case ExpectAfterVarTypeName:
		goto AFTER_VAR_TYPE_NAME
	}

DEFINITION_END:
	i.levelSel, i.expect = 0, ExpectDef
	// Expect end of file
	i.skipSTNRC()
	if i.head < len(i.str) {
		goto DEFINITION
	}
	return Error{}

ERROR:
	{
		var atIndex rune
		if i.head < len(i.str) {
			atIndex, _ = utf8.DecodeRune(i.str[i.head:])
		}
		return Error{
			Index:       i.head,
			AtIndex:     atIndex,
			Code:        i.errc,
			Expectation: i.expect,
		}
	}
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
	ExpectEndOfString
	ExpectColumnAfterArg
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
	case ExpectEndOfString:
		return "end of string"
	case ExpectColumnAfterArg:
		return "column after argument name"
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
	TokenSel
	TokenSelEnd
	TokenFragTypeCond
	TokenFragName
	TokenFragInline
	TokenFragRef
	TokenField
	TokenArg
	TokenArr
	TokenArrEnd
	TokenStr
	TokenNum
	TokenTrue
	TokenFalse
	TokenNull
	TokenVarName
	TokenVarTypeName
	TokenVarTypeArr
	TokenVarTypeArrEnd
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
