package gqlscan

import (
	"fmt"
	"unicode/utf8"
)

// Scan calls fn for every token it scans in str.
// If fn returns true then an error with code ErrCallbackFn is returned.
// If the returned error code == 0 then there was no error during the scan,
// this can also be checked using err.IsErr().
//
// WARNING: *Iterator passed to fn should never be aliased and
// used after Scan returns because it's returned to the pool
// and may be acquired by another call to Scan!
func Scan(str []byte, fn func(*Iterator) (err bool)) Error {
	i := acquireIterator(str)
	defer iteratorPool.Put(i)

	// inDefVal triggers different expectations after values
	// when the iterator is in a variable default value definition.
	var inDefVal bool
	var typeArrLvl int
	var dirOn dirTarget

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
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordMutation() {
		// Mutation
		i.token = TokenDefMut
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head += len("mutation")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordSubscription() {
		// Subscription
		i.token = TokenDefSub
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head += len("subscription")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
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

AFTER_DEF_KEYWORD:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	switch i.str[i.head] {
	case '#':
		goto COMMENT
	case '{':
		i.expect = ExpectSelSet
		goto SELECTION_SET
	case '(':
		// Variable list
		i.tail = -1
		i.token = TokenVarList
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		i.expect = ExpectVar
		goto OPR_VAR
	case '@':
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	i.expect = ExpectOprName
	goto NAME

AFTER_DIR_NAME:
	i.skipSTNRC()
	switch dirOn {
	case dirField:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		case '{':
			// Field selector expands without arguments
			i.expect = ExpectSelSet
			goto SELECTION_SET
		default:
			i.expect, dirOn = ExpectAfterSelection, 0
			goto AFTER_SELECTION
		}
	case dirOpr:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectSelSet, 0
			goto SELECTION_SET
		}
	case dirVar:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		case ')':
			dirOn = 0
			goto VAR_LIST_END
		default:
			i.expect, dirOn = ExpectVar, 0
			goto OPR_VAR
		}
	case dirFragRef:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectAfterSelection, 0
			goto AFTER_SELECTION
		}
	case dirFragInlineOrDef:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectSelSet, 0
			goto SELECTION_SET
		}
	default:
		// This line is only executed if we forgot to handle a dirOn case.
		panic(fmt.Errorf("unhandled dirOn case: %#v", dirOn))
	}

AFTER_DIR_ARGS:
	i.skipSTNRC()
	switch dirOn {
	case dirField:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		case '{':
			i.expect = ExpectSelSet
			goto SELECTION_SET
		default:
			i.expect, dirOn = ExpectAfterSelection, 0
			goto AFTER_SELECTION
		}
	case dirOpr:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectSelSet, 0
			goto SELECTION_SET
		}
	case dirVar:
		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectAfterVarType, 0
			goto OPR_VAR
		}
	case dirFragRef:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect, dirOn = ExpectAfterSelection, 0
			goto AFTER_SELECTION
		}
	case dirFragInlineOrDef:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.expect = ExpectDir
			goto DIR_NAME
		default:
			i.expect = ExpectSelSet
			goto SELECTION_SET
		}
	default:
		// This line is only executed if we forgot to handle a dirOn case.
		panic(fmt.Errorf("unhandled dirOn case: %#v", dirOn))
	}

AFTER_KEYWORD_FRAGMENT:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	goto NAME

OPR_VAR:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}

	// Variable name
	if i.str[i.head] != '$' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.expect = ExpectVarName
	goto VAR_NAME

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
	} else if i.str[i.head] == '@' {
		i.head++
		dirOn, i.expect = dirVar, ExpectDir
		goto DIR_NAME
	} else if i.str[i.head] == '=' {
		i.head++
		i.skipSTNRC()
		i.expect, inDefVal = ExpectVal, true
		goto VALUE
	} else if i.str[i.head] == ')' {
		goto VAR_LIST_END
	}
	i.expect = ExpectAfterVarType
	goto OPR_VAR

VAR_LIST_END:
	i.tail = -1
	i.token = TokenVarListEnd
	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}
	i.head++
	i.skipSTNRC()
	i.expect = ExpectSelSet
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		dirOn, i.expect = dirOpr, ExpectDirName
		goto AFTER_DIR_NAME
	} else if i.str[i.head] == '@' {
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	goto SELECTION_SET

SELECTION_SET:
	i.skipSTNRC()
	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != '{' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.tail = -1
	i.token = TokenSet
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
		goto SEL_END
	}
	i.expect = ExpectSel
	goto SELECTION

SEL_END:
	i.tail = -1
	i.token = TokenSetEnd
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

VALUE:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	switch i.str[i.head] {
	case '#':
		goto COMMENT
	case '{':
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
	case '[':
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

	case '"':
		i.head++
		i.tail = i.head

		if i.head+1 < len(i.str) &&
			i.str[i.head] == '"' &&
			i.str[i.head+1] == '"' {
			i.head += 2
			i.tail = i.head
			goto BLOCK_STRING
		}

		// String value
		escaped := false
		if i.head < len(i.str) && i.str[i.head] == '"' {
			goto AFTER_STR_VAL
		}
		for {
			for !escaped && i.head+7 < len(i.str) {
				// Fast path
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
				if i.str[i.head] == '"' ||
					i.str[i.head] == '\\' ||
					i.str[i.head] < 0x20 {
					break
				}
				i.head++
			}
			if i.head >= len(i.str) {
				break
			}
			if i.isHeadCtrl() {
				i.errc = ErrUnexpToken
				i.expect = ExpectEndOfString
				goto ERROR
			}
			if escaped {
				switch i.str[i.head] {
				case '\\':
					// Backslash
					i.head++
				case '/':
					// Solidus
					i.head++
				case '"':
					// Double-quotes
					i.head++
				case 'b':
					// Backspace
					i.head++
				case 'f':
					// Form-feed
					i.head++
				case 'r':
					// Carriage-return
					i.head++
				case 'n':
					// Line-break
					i.head++
				case 't':
					// Tab
					i.head++
				case 'u':
					// Unicode sequence
					i.head++
					if i.head >= len(i.str) {
						i.errc = ErrUnexpEOF
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if !i.isHeadHexDigit() {
						i.errc = ErrUnexpToken
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					i.head++
					if i.head >= len(i.str) {
						i.errc = ErrUnexpEOF
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if !i.isHeadHexDigit() {
						i.errc = ErrUnexpToken
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					i.head++
					if i.head >= len(i.str) {
						i.errc = ErrUnexpEOF
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if !i.isHeadHexDigit() {
						i.errc = ErrUnexpToken
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					i.head++
					if i.head >= len(i.str) {
						i.errc = ErrUnexpEOF
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if !i.isHeadHexDigit() {
						i.errc = ErrUnexpToken
						i.expect = ExpectEscapedUnicodeSequence
						goto ERROR
					}
				default:
					i.errc = ErrUnexpToken
					i.expect = ExpectEscapedSequence
					goto ERROR
				}
				escaped = false
				continue
			} else if i.str[i.head] == '"' {
				goto AFTER_STR_VAL
			} else if i.str[i.head] == '\\' {
				escaped = true
			}
			i.head++
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
	case '$':
		if inDefVal {
			i.errc, i.expect = ErrUnexpToken, ExpectDefaultVarVal
			goto ERROR
		}

		// Variable reference
		i.head++

		// Variable name
		i.expect = ExpectVarRefName
		goto VAR_NAME

	case 'n':
		// Null
		if i.head+4 < len(i.str) &&
			i.str[i.head+3] == 'l' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'u' &&
			i.str[i.head] == 'n' &&
			(i.str[i.head+4] == ' ' ||
				i.str[i.head+4] == '\t' ||
				i.str[i.head+4] == '\r' ||
				i.str[i.head+4] == '\n' ||
				i.str[i.head+4] == ',' ||
				i.str[i.head+4] == ')' ||
				i.str[i.head+4] == '}' ||
				i.str[i.head+4] == '{' ||
				i.str[i.head+4] == ']' ||
				i.str[i.head+4] == '[' ||
				i.str[i.head+4] == '#') {
			i.tail = -1
			i.head += len("null")

			// Callback for null value
			i.token = TokenNull
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
		} else {
			i.expect = ExpectValEnum
			goto NAME
		}
	case 't':
		if i.head+4 < len(i.str) &&
			i.str[i.head+3] == 'e' &&
			i.str[i.head+2] == 'u' &&
			i.str[i.head+1] == 'r' &&
			i.str[i.head] == 't' &&
			(i.str[i.head+4] == ' ' ||
				i.str[i.head+4] == '\t' ||
				i.str[i.head+4] == '\r' ||
				i.str[i.head+4] == '\n' ||
				i.str[i.head+4] == ',' ||
				i.str[i.head+4] == ')' ||
				i.str[i.head+4] == '}' ||
				i.str[i.head+4] == '{' ||
				i.str[i.head+4] == ']' ||
				i.str[i.head+4] == '[' ||
				i.str[i.head+4] == '#') {
			i.tail = -1
			i.head += len("true")

			// Callback for true value
			i.token = TokenTrue
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
		} else {
			i.expect = ExpectValEnum
			goto NAME
		}
	case 'f':
		// False
		if i.head+5 < len(i.str) &&
			i.str[i.head+4] == 'e' &&
			i.str[i.head+3] == 's' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'a' &&
			i.str[i.head] == 'f' &&
			(i.str[i.head+5] == ' ' ||
				i.str[i.head+5] == '\t' ||
				i.str[i.head+5] == '\r' ||
				i.str[i.head+5] == '\n' ||
				i.str[i.head+5] == ',' ||
				i.str[i.head+5] == ')' ||
				i.str[i.head+5] == '}' ||
				i.str[i.head+5] == '{' ||
				i.str[i.head+5] == ']' ||
				i.str[i.head+5] == '[' ||
				i.str[i.head+5] == '#') {
			i.tail = -1
			i.head += len("false")

			// Callback for false value
			i.token = TokenFalse
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
		} else {
			i.expect = ExpectValEnum
			goto NAME
		}
	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
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
					i.token = TokenInt
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
				i.token = TokenInt
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
			i.token = TokenInt
			goto ON_NUM_VAL
		}
		// Continue to fraction

	FRACTION:
		_ = 0 // Make code coverage count the label above
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
				i.token = TokenFloat
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
			i.token = TokenFloat
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
				i.token = TokenFloat
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
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
	default:
		// Invalid value
		i.expect = ExpectValEnum
		goto NAME
	}
	i.expect = ExpectAfterValue
	goto AFTER_VALUE_COMMENT

BLOCK_STRING:
	i.expect = ExpectEndOfBlockString
	for ; i.head < len(i.str); i.head++ {
		if i.str[i.head] == '\\' &&
			i.str[i.head+3] == '"' &&
			i.str[i.head+2] == '"' &&
			i.str[i.head+1] == '"' {
			i.head += 3
			continue
		}
		if i.str[i.head] == '"' &&
			i.str[i.head+1] == '"' &&
			i.str[i.head+2] == '"' {
			i.token = TokenStrBlock
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head += 3
			goto AFTER_VALUE_COMMENT
		}
	}

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

	if inDefVal {
		inDefVal = false
		if i.str[i.head] == ')' {
			goto VAR_LIST_END
		} else if i.str[i.head] == '@' {
			i.head++
			dirOn, i.expect = dirVar, ExpectDir
			goto DIR_NAME
		}
		i.expect = ExpectVar
		goto OPR_VAR
	}

	if i.str[i.head] == ')' {
		// End of argument list
		i.tail = -1
		i.token = TokenArgListEnd
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		i.expect = ExpectAfterArgList
		goto AFTER_ARG_LIST
	}

	// Proceed to the next argument
	i.expect = ExpectArgName
	goto NAME

AFTER_ARG_LIST:
	if dirOn != 0 {
		goto AFTER_DIR_ARGS
	}

	i.skipSTNRC()

	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}

	if i.str[i.head] == '{' {
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '}' {
		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION
	} else if i.str[i.head] == '@' {
		i.head++
		dirOn, i.expect = dirField, ExpectDir
		goto DIR_NAME
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
		i.expect = ExpectFieldNameOrAlias
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
	goto SPREAD

SPREAD:
	i.skipSTNRC()
	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token, i.tail = TokenFragInline, -1
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '@' {
		i.token, i.tail = TokenFragInline, -1
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect, dirOn = ExpectDirName, dirFragInlineOrDef
		goto AFTER_DIR_NAME
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
			i.head += len("on")
			i.expect = ExpectFragInlined
			goto FRAG_TYPE_COND
		}
	}
	// ...fragmentName
	i.expect = ExpectSpreadName
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

VAR_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	goto NAME

DIR_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	}
	i.expect = ExpectDirName
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
	i.head++
	for {
		if i.head+7 >= len(i.str) {
			goto NAME_LOOP
		}
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
		if !i.isHeadNameBody() {
			break
		}
		i.head++
	}
	if i.isHeadSNTRC() {
		goto AFTER_NAME
	} else if i.isHeadCtrl() {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	goto AFTER_NAME

NAME_LOOP:
	_ = 0 // Make code coverage count the label above
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
	if i.head < len(i.str) && i.str[i.head] == '!' {
		i.tail = -1
		i.token = TokenVarTypeNotNull
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
	}
	goto AFTER_VAR_TYPE_NOT_NULL

AFTER_VAR_TYPE_NOT_NULL:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == ']' {
		if typeArrLvl < 1 {
			i.errc, i.expect = ErrUnexpToken, ExpectVar
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

		i.skipSTNRC()
		if i.head < len(i.str) && i.str[i.head] == '!' {
			i.tail = -1
			i.token = TokenVarTypeNotNull
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
		}

		if typeArrLvl > 0 {
			goto AFTER_VAR_TYPE_NAME
		}
	}
	i.expect = ExpectAfterVarType
	goto AFTER_VAR_TYPE

AFTER_FIELD_NAME:
	i.skipSTNRC()
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	// Lookahead
	switch i.str[i.head] {
	case '(':
		// Argument list
		i.tail = -1
		i.token = TokenArgList
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.head++
		i.skipSTNRC()
		i.expect = ExpectArgName
		goto ARG_LIST
	case '{':
		// Field selector expands without arguments
		i.expect = ExpectSelSet
		goto SELECTION_SET
	case '#':
		i.expect = ExpectAfterFieldName
		goto COMMENT
	case '@':
		i.head++
		dirOn, i.expect = dirField, ExpectDir
		goto DIR_NAME
	}
	i.expect = ExpectAfterSelection
	goto AFTER_SELECTION

AFTER_NAME:
	_ = 0 // Make code coverage count the label above
	switch i.expect {
	case ExpectValEnum:
		// Enum value
		i.token = TokenEnumVal
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectAfterValue
		goto AFTER_VALUE_COMMENT

	case ExpectFieldNameOrAlias:
		head := i.head
		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		} else if i.str[i.head] == ':' {
			h2 := i.head
			i.head = head
			// Callback for field alias name
			i.token = TokenFieldAlias
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			i.head = h2 + 1
			i.skipSTNRC()
			i.expect = ExpectFieldName
			goto NAME
		}
		i.head = head
		fallthrough

	case ExpectFieldName:
		// Callback for field name
		i.token = TokenField
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		goto AFTER_FIELD_NAME

	case ExpectDirName:
		// Callback directive name
		i.token = TokenDirName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		goto AFTER_DIR_NAME

	case ExpectArgName:
		// Callback for argument name
		i.token = TokenArgName
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

	case ExpectVarName:
		i.token = TokenVarName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.expect = ExpectColumnAfterVar
		goto AFTER_DECL_VAR_NAME

	case ExpectOprName:
		i.token = TokenOprName
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}
		i.skipSTNRC()

		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectSelSet
			goto ERROR
		}
		switch i.str[i.head] {
		case '{':
			i.expect = ExpectSelSet
			goto SELECTION_SET
		case '(':
			// Variable list
			i.tail = -1
			i.token = TokenVarList
			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}
			i.head++
			i.expect = ExpectVar
			goto OPR_VAR
		case '@':
			i.head++
			dirOn, i.expect = dirOpr, ExpectDir
			goto DIR_NAME
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

		i.expect, dirOn = ExpectDirName, dirFragInlineOrDef
		goto AFTER_DIR_NAME

	case ExpectSpreadName:
		i.token = TokenNamedSpread
		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		i.expect, dirOn = ExpectDirName, dirFragRef
		goto AFTER_DIR_NAME

	case ExpectFragName:
		if i.head-i.tail == 2 &&
			i.str[i.tail+1] == 'n' &&
			i.str[i.tail] == 'o' {
			i.errc, i.head = ErrIllegalFragName, i.tail
			goto ERROR
		}

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

		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		} else if i.str[i.head] == '@' {
			dirOn = dirFragInlineOrDef
			goto AFTER_DIR_NAME
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
	for {
		if i.head+7 >= len(i.str) {
			for ; i.head < len(i.str) && i.str[i.head] != '\n'; i.head++ {
			}
			break
		}
		if i.str[i.head] != '\n' &&
			i.str[i.head+1] != '\n' &&
			i.str[i.head+2] != '\n' &&
			i.str[i.head+3] != '\n' &&
			i.str[i.head+4] != '\n' &&
			i.str[i.head+5] != '\n' &&
			i.str[i.head+6] != '\n' &&
			i.str[i.head+7] != '\n' {
			i.head += 8
			continue
		}
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
		i.head++
		if i.str[i.head] == '\n' {
			break
		}
	}
	i.tail = -1
	i.skipSTNRC()
	switch i.expect {
	case ExpectVarRefName, ExpectVarName:
		goto VAR_NAME
	case ExpectDef:
		goto DEFINITION
	case ExpectDir:
		goto DIR_NAME
	case ExpectDirName:
		goto AFTER_DIR_NAME
	case ExpectSelSet:
		goto SELECTION_SET
	case ExpectSel:
		goto SELECTION
	case ExpectAfterSelection:
		goto AFTER_SELECTION
	case ExpectVar:
		goto OPR_VAR
	case ExpectArgName:
		goto ARG_LIST
	case ExpectColumnAfterArg:
		goto COLUMN_AFTER_ARG_NAME
	case ExpectVal:
		goto VALUE
	case ExpectAfterFieldName:
		goto AFTER_FIELD_NAME
	case ExpectAfterValue:
		goto AFTER_VALUE_COMMENT
	case ExpectAfterArgList:
		goto AFTER_ARG_LIST
	case ExpectAfterDefKeyword:
		goto AFTER_DEF_KEYWORD
	case ExpectFragName:
		goto AFTER_KEYWORD_FRAGMENT
	case ExpectFragKeywordOn:
		goto FRAG_KEYWORD_ON
	case ExpectFragInlined:
		goto FRAG_TYPE_COND
	case ExpectFragTypeCond:
		goto FRAG_TYPE_COND
	case ExpectFrag:
		goto SPREAD
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
