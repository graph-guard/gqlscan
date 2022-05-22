package gqlscan

import (
	"fmt"
	"unicode/utf8"
)

// ScanAll calls fn for every token it scans in str.
//
// WARNING: *Iterator passed to fn should never be aliased and
// used after ScanAll returns!
func ScanAll(str []byte, fn func(*Iterator)) Error {
	i := acquireIterator(str)
	defer iteratorPool.Put(i)

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
		fn(i)
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.isHeadKeywordQuery() {
		// Query
		i.token = TokenDefQry
		fn(i)
		i.head += len("query")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordMutation() {
		// Mutation
		i.token = TokenDefMut
		fn(i)
		i.head += len("mutation")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordSubscription() {
		// Subscription
		i.token = TokenDefSub
		fn(i)
		i.head += len("subscription")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordFragment() {
		// Fragment
		i.tail = -1
		i.token = TokenDefFrag
		fn(i)
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
		fn(i)
		i.head++
		i.expect = ExpectVarName
		goto OPR_VAR
	case '@':
		i.head++
		i.skipSTNRC()
		dirOn, i.expect = dirOpr, ExpectDirName
		goto NAME
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
			fn(i)
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
			fn(i)
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
			fn(i)
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
		case ')':
			dirOn = 0
			goto VAR_LIST_END
		default:
			i.expect, dirOn = ExpectVarName, 0
			goto NAME
		}
	case dirFragRef:
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
			fn(i)
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
			fn(i)
			i.head++
			i.skipSTNRC()
			i.expect = ExpectArgName
			goto ARG_LIST
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
		default:
			i.expect = ExpectSelSet
			goto SELECTION_SET
		}
	default:
		// The panic is triggered only if we forgot to handle a directive type.
		panic(fmt.Errorf("unhandled directive type: %q", dirOn))
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
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
		case ')':
			goto VAR_LIST_END
		default:
			i.expect, dirOn = ExpectAfterVarType, 0
			goto OPR_VAR
		}
	case dirFragRef:
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '@':
			i.head++
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
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
			i.skipSTNRC()
			i.expect = ExpectDirName
			goto NAME
		default:
			i.expect = ExpectSelSet
			goto SELECTION_SET
		}
	default:
		panic(fmt.Errorf("no: %#v", dirOn))
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
	} else if i.str[i.head] == '@' {
		i.head++
		i.skipSTNRC()
		dirOn, i.expect = dirVar, ExpectDirName
		goto NAME
	} else if i.str[i.head] == ')' {
		goto VAR_LIST_END
	}
	i.expect = ExpectAfterVarType
	goto OPR_VAR

VAR_LIST_END:
	i.tail = -1
	i.token = TokenVarListEnd
	fn(i)
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
		i.skipSTNRC()
		dirOn, i.expect = dirOpr, ExpectDirName
		goto NAME
	}
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
	fn(i)
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
	i.token = TokenSelEnd
	fn(i)
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
		fn(i)
		i.stackPush(TokenObj)
		i.head++
		i.skipSTNRC()

		i.expect = ExpectObjFieldName
		goto NAME
	case '[':
		i.tail = -1
		// Callback for argument
		i.token = TokenArr
		fn(i)
		i.head++
		i.skipSTNRC()

		// Lookahead
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] == ']' {
			i.token = TokenArrEnd
			fn(i)
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
		fn(i)
		// Advance head index to include the closing double-quotes
		i.head++
	case '$':
		// Variable reference
		i.head++

		// Variable name
		i.expect = ExpectVarRefName
		goto NAME

	case 'n':
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
		fn(i)
	case 't':
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
		fn(i)
	case 'f':
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
		fn(i)
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
		fn(i)
	default:
		// Invalid value
		i.errc = ErrInvalVal
		i.expect = ExpectVal
		goto ERROR
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
			fn(i)
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
			fn(i)

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
			fn(i)
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
		i.tail = -1
		i.token = TokenArgListEnd
		fn(i)
		// End of argument list
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
		i.skipSTNRC()
		dirOn, i.expect = dirField, ExpectDirName
		goto NAME
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
			i.head += len("on")
			i.expect = ExpectFragInlined
			goto FRAG_TYPE_COND
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
		fn(i)
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
		fn(i)
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
			i.errc, i.expect = ErrUnexpToken, ExpectVarName
			goto ERROR
		}
		i.tail = -1
		i.token = TokenVarTypeArrEnd
		fn(i)
		i.head++
		typeArrLvl--

		i.skipSTNRC()
		if i.head < len(i.str) && i.str[i.head] == '!' {
			i.tail = -1
			i.token = TokenVarTypeNotNull
			fn(i)
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
		fn(i)
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
		i.skipSTNRC()
		dirOn, i.expect = dirField, ExpectDirName
		goto NAME
	}
	i.expect = ExpectAfterSelection
	goto AFTER_SELECTION

AFTER_NAME:
	switch i.expect {
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
			fn(i)

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
		fn(i)
		goto AFTER_FIELD_NAME

	case ExpectDirName:
		// Callback directive name
		i.token = TokenDirName
		fn(i)
		goto AFTER_DIR_NAME

	case ExpectArgName:
		// Callback for argument name
		i.token = TokenArg
		fn(i)
		i.skipSTNRC()
		i.expect = ExpectColumnAfterArg
		goto COLUMN_AFTER_ARG_NAME

	case ExpectObjFieldName:
		// Callback for object field
		i.token = TokenObjField
		fn(i)

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
		fn(i)
		i.expect = ExpectAfterValue
		goto AFTER_VALUE_COMMENT

	case ExpectVarType:
		i.token = TokenVarTypeName
		fn(i)
		i.expect = ExpectAfterVarTypeName
		goto AFTER_VAR_TYPE_NAME

	case ExpectVarName, ExpectAfterVarType:
		i.token = TokenVarName
		fn(i)
		i.expect = ExpectColumnAfterVar
		goto AFTER_DECL_VAR_NAME

	case ExpectOprName:
		i.token = TokenOprName
		fn(i)
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
			fn(i)
			i.head++
			i.expect = ExpectVarName
			goto OPR_VAR
		case '@':
			i.head++
			i.skipSTNRC()
			dirOn, i.expect = dirOpr, ExpectDirName
			goto NAME
		}
		i.errc = ErrUnexpToken
		i.expect = ExpectSelSet
		goto ERROR

	case ExpectFragInlined:
		i.token = TokenFragInline
		fn(i)

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

	case ExpectFragRef:
		i.token = TokenFragRef
		fn(i)

		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		} else if i.str[i.head] == '@' {
			dirOn = dirFragRef
			goto AFTER_DIR_NAME
		}

		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION

	case ExpectFragName:
		i.token = TokenFragName
		fn(i)
		i.expect = ExpectFragKeywordOn
		goto FRAG_KEYWORD_ON

	case ExpectFragTypeCond:
		i.token = TokenFragTypeCond
		fn(i)

		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
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
	case ExpectDef:
		goto DEFINITION
	case ExpectDirName:
		goto AFTER_DIR_NAME
	case ExpectSelSet:
		goto SELECTION_SET
	case ExpectSel:
		goto SELECTION
	case ExpectAfterSelection:
		goto AFTER_SELECTION
	case ExpectVarName:
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
