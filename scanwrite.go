package gqlscan

import (
	"fmt"
	"unicode/utf8"
)

// ScanWrite writes every token it scans in str to buffer.
// Returns the number of tokens written to the buffer.
// Returns ErrBufferOverflow if the buffer is too small.
func ScanWrite(str []byte, buffer []TokenRef) (int, Error) {
	i := acquireIterator(str)
	defer iteratorPool.Put(i)

	var typeArrLvl int
	bufferIndex := 0

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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenDefQry,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.isHeadKeywordQuery() {
		// Query
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenDefQry,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.head += len("query")
		i.expect = ExpectAfterKeywordQuery
		goto AFTER_KEYWORD_QUERY
	} else if i.isHeadKeywordMutation() {
		// Mutation
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenDefMut,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.head += len("mutation")
		i.expect = ExpectAfterKeywordMutation
		goto AFTER_KEYWORD_MUTATION
	} else if i.isHeadKeywordFragment() {
		// Fragment
		i.tail = -1
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenDefFrag,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		i.tail = -1
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarList,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		i.tail = -1
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarList,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
	i.tail = -1
	if bufferIndex >= len(buffer) {
		goto ERROR_BUFFER_OVERFLOW
	}
	buffer[bufferIndex] = TokenRef{
		Type:        TokenVarListEnd,
		IndexTail:   i.tail,
		IndexHead:   i.head,
		LevelSelect: i.levelSel,
	}
	bufferIndex++
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
	if bufferIndex >= len(buffer) {
		goto ERROR_BUFFER_OVERFLOW
	}
	buffer[bufferIndex] = TokenRef{
		Type:        TokenSel,
		IndexTail:   i.tail,
		IndexHead:   i.head,
		LevelSelect: i.levelSel,
	}
	bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenSelEnd,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenObj,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.stackPush(TokenObj)
		i.head++
		i.skipSTNRC()

		i.expect = ExpectObjFieldName
		goto NAME
	} else if i.str[i.head] == '[' {
		i.tail = -1
		// Callback for argument
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenArr,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.head++
		i.skipSTNRC()

		// Lookahead
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectVal
			goto ERROR
		} else if i.str[i.head] == ']' {
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenArrEnd,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenStr,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenNull,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenTrue,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenFalse,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenNum,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenObjEnd,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++

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
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenArrEnd,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenArgListEnd,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarTypeArr,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarTypeArrEnd,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenField,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++

		// Lookahead
		i.skipSTNRC()
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			i.expect = ExpectFieldName
			goto ERROR
		} else if i.str[i.head] == '(' {
			// Argument list
			i.tail = -1
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenArgList,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++
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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenArg,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.skipSTNRC()
		i.expect = ExpectColumnAfterArg
		goto COLUMN_AFTER_ARG_NAME

	case ExpectObjFieldName:
		// Callback for object field
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenObjField,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++

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
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarRef,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectAfterValue
		goto AFTER_VALUE_COMMENT

	case ExpectVarType:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarTypeName,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectAfterVarTypeName
		goto AFTER_VAR_TYPE_NAME

	case ExpectVarName, ExpectAfterVarType:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenVarName,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectColumnAfterVar
		goto AFTER_DECL_VAR_NAME

	case ExpectQryName:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenQryName,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
			i.tail = -1
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenVarList,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++
			i.head++
			i.expect = ExpectVarName
			goto QUERY_VAR
		}
		i.errc = ErrUnexpToken
		i.expect = ExpectSelSet
		goto ERROR

	case ExpectMutName:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenMutName,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
			i.tail = -1
			if bufferIndex >= len(buffer) {
				goto ERROR_BUFFER_OVERFLOW
			}
			buffer[bufferIndex] = TokenRef{
				Type:        TokenVarList,
				IndexTail:   i.tail,
				IndexHead:   i.head,
				LevelSelect: i.levelSel,
			}
			bufferIndex++
			i.head++
			i.expect = ExpectVarName
			goto QUERY_VAR
		}
		i.errc = ErrUnexpToken
		i.expect = ExpectSelSet
		goto ERROR

	case ExpectFragInlined:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenFragInline,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectSelSet
		goto SELECTION_SET

	case ExpectFragRef:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenFragRef,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectAfterSelection
		goto AFTER_SELECTION

	case ExpectFragName:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenFragName,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
		i.expect = ExpectFragKeywordOn
		goto FRAG_KEYWORD_ON

	case ExpectFragTypeCond:
		if bufferIndex >= len(buffer) {
			goto ERROR_BUFFER_OVERFLOW
		}
		buffer[bufferIndex] = TokenRef{
			Type:        TokenFragTypeCond,
			IndexTail:   i.tail,
			IndexHead:   i.head,
			LevelSelect: i.levelSel,
		}
		bufferIndex++
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
	return bufferIndex, Error{}

ERROR:
	{
		var atIndex rune
		if i.head < len(i.str) {
			atIndex, _ = utf8.DecodeRune(i.str[i.head:])
		}
		return bufferIndex, Error{
			Index:       i.head,
			AtIndex:     atIndex,
			Code:        i.errc,
			Expectation: i.expect,
		}
	}

ERROR_BUFFER_OVERFLOW:
	{
		var atIndex rune
		if i.head < len(i.str) {
			atIndex, _ = utf8.DecodeRune(i.str[i.head:])
		}
		return bufferIndex, Error{
			Index:       i.head,
			AtIndex:     atIndex,
			Code:        ErrBufferOverflow,
			Expectation: 0,
		}
	}
}