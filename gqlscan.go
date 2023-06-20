package gqlscan

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
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

	/*<scan_body>*/
	i := iteratorPool.Get().(*Iterator)
	i.stackReset()
	i.expect = ExpectDef
	i.tail, i.head = -1, 0
	i.str = str
	i.levelSel = 0
	i.errc = 0
	defer iteratorPool.Put(i)

	// inDefVal triggers different expectations after values
	// when the iterator is in a variable default value definition.
	var inDefVal bool
	var typeArrLvl int
	var dirOn dirTarget

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectDef
		goto ERROR
	}
	/*</check_eof>*/

	/*<l_definition>*/
DEFINITION:
	if i.head >= len(i.str) {
		goto DEFINITION_END
	} else if i.str[i.head] == '#' {
		i.expect = ExpectDef
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token = TokenDefQry
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.isHeadKeywordQuery() {
		// Query
		i.token = TokenDefQry
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head += len("query")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordMutation() {
		// Mutation
		i.token = TokenDefMut
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head += len("mutation")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordSubscription() {
		// Subscription
		i.token = TokenDefSub
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head += len("subscription")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordFragment() {
		// Fragment
		i.tail = -1
		i.token = TokenDefFrag
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head += len("fragment")
		i.expect = ExpectFragName
		goto AFTER_KEYWORD_FRAGMENT
	}

	i.errc = ErrUnexpToken
	i.expect = ExpectDef
	goto ERROR
	/*</l_definition>*/

	/*<l_after_def_keyword>*/
AFTER_DEF_KEYWORD:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

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
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++
		i.expect = ExpectVar
		goto OPR_VAR
	case '@':
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	i.expect = ExpectOprName

	/*<name>*/
	// Followed by oprname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectOprName after name>
	i.token = TokenOprName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	goto AFTER_OPR_NAME
	// </ExpectOprName after name>

	/*</name>*/

	/*</l_after_def_keyword>*/

	/*<l_after_dir_name>*/
AFTER_DIR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch dirOn {
	case dirField:

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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
	/*</l_after_dir_name>*/

	/*<l_after_dir_args>*/
AFTER_DIR_ARGS:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch dirOn {
	case dirField:

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		/*</check_eof>*/

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
	/*</l_after_dir_args>*/

	/*<l_after_keyword_fragment>*/
AFTER_KEYWORD_FRAGMENT:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fragname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragName after name>
	if i.head-i.tail == 2 &&
		i.str[i.tail+1] == 'n' &&
		i.str[i.tail] == 'o' {
		i.errc, i.head = ErrIllegalFragName, i.tail
		goto ERROR
	}
	i.token = TokenFragName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect = ExpectFragKeywordOn
	goto FRAG_KEYWORD_ON
	// </ExpectFragName after name>

	/*</name>*/

	/*</l_after_keyword_fragment>*/

	/*<l_opr_var>*/
OPR_VAR:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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
	/*</l_opr_var>*/

	/*<l_after_var_type>*/
AFTER_VAR_TYPE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect, inDefVal = ExpectVal, true
		goto VALUE
	} else if i.str[i.head] == ')' {
		goto VAR_LIST_END
	}
	i.expect = ExpectAfterVarType
	goto OPR_VAR
	/*</l_after_var_type>*/

	/*<l_var_list_end>*/
VAR_LIST_END:
	i.tail = -1
	i.token = TokenVarListEnd
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.head++

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectSelSet

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		dirOn, i.expect = dirOpr, ExpectDirName
		goto AFTER_DIR_NAME
	} else if i.str[i.head] == '@' {
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	goto SELECTION_SET
	/*</l_var_list_end>*/

	/*<l_selection_set>*/
SELECTION_SET:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != '{' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.tail = -1
	i.token = TokenSet
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.levelSel++
	i.head++
	i.expect = ExpectSel
	goto SELECTION
	/*</l_selection_set>*/

	/*<l_after_selection>*/
AFTER_SELECTION:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '}' {
		goto SEL_END
	}
	i.expect = ExpectSel
	goto SELECTION
	/*</l_after_selection>*/

	/*<l_sel_end>*/
SEL_END:
	i.tail = -1
	i.token = TokenSetEnd
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.levelSel--
	i.head++

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.levelSel < 1 {
		goto DEFINITION_END
	}
	goto AFTER_SELECTION
	/*</l_sel_end>*/

	/*<l_value>*/
VALUE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	switch i.str[i.head] {
	case '#':
		goto COMMENT

	case '{':
		// Object begin
		i.tail = -1
		// Callback for argument
		i.token = TokenObj
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.stackPush(TokenObj)
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect = ExpectObjFieldName

		/*<name>*/
		// Followed by objfieldname>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectObjFieldName after name>
		i.token = TokenObjField
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectColObjFieldName
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] != ':' {
			i.errc = ErrUnexpToken
			i.expect = ExpectColObjFieldName
			goto ERROR
		}
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect = ExpectVal
		goto VALUE
	// </ExpectObjFieldName after name>

	/*</name>*/

	case '[':
		i.tail = -1
		// Callback for argument
		i.token = TokenArr
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		// Lookahead

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectVal
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == ']' {
			i.token = TokenArrEnd
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
		}
		i.stackPush(TokenArr)
		i.expect = ExpectAfterValueInner
		goto AFTER_VALUE_INNER

	case '"':

		/*<str>*/
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
		if i.head < len(i.str) && i.str[i.head] == '"' {
			goto AFTER_STR_VAL
		}
		for {
			if i.head+15 < len(i.str) {
				// Fast path
				if lutStr[i.str[i.head]] == 1 {
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+1]] == 1 {
					i.head += 1
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+2]] == 1 {
					i.head += 2
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+3]] == 1 {
					i.head += 3
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+4]] == 1 {
					i.head += 4
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+5]] == 1 {
					i.head += 5
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+6]] == 1 {
					i.head += 6
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+7]] == 1 {
					i.head += 7
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+8]] == 1 {
					i.head += 8
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+9]] == 1 {
					i.head += 9
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+10]] == 1 {
					i.head += 10
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+11]] == 1 {
					i.head += 11
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+12]] == 1 {
					i.head += 12
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+13]] == 1 {
					i.head += 13
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+14]] == 1 {
					i.head += 14
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+15]] == 1 {
					i.head += 15
					goto CHECK_STR_CHARACTER
				}
				i.head += 16
			}
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectEndOfString
				goto ERROR
			}
		CHECK_STR_CHARACTER:
			switch i.str[i.head] {
			case '"':
				// i.head++
				goto AFTER_STR_VAL
			case '\\':
				i.head++
				if i.head >= len(i.str) {
					i.errc, i.expect = ErrUnexpEOF, ExpectEscapedSequence
					goto ERROR
				}
				switch i.str[i.head] {
				case '\\', '/', '"', 'b', 'f', 'r', 'n', 't':
					i.head++
				case 'u':
					if i.head+4 >= len(i.str) {
						switch {
						case i.head+1 < len(i.str) && lutHex[i.str[i.head+1]] != 2:
							i.head++
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+2 < len(i.str) && lutHex[i.str[i.head+2]] != 2:
							i.head += 2
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+3 < len(i.str) && lutHex[i.str[i.head+3]] != 2:
							i.head += 3
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+4 < len(i.str) && lutHex[i.str[i.head+4]] != 2:
							i.head += 4
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						}
						i.head = len(i.str)
						i.errc, i.expect = ErrUnexpEOF, ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if lutHex[i.str[i.head+4]] != 2 ||
						lutHex[i.str[i.head+3]] != 2 ||
						lutHex[i.str[i.head+2]] != 2 ||
						lutHex[i.str[i.head+1]] != 2 {
						switch {
						case lutHex[i.str[i.head+1]] != 2:
							i.head++
						case lutHex[i.str[i.head+2]] != 2:
							i.head += 2
						case lutHex[i.str[i.head+3]] != 2:
							i.head += 3
						case lutHex[i.str[i.head+4]] != 2:
							i.head += 4
						}
						i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
						goto ERROR
					}
					i.head += 4
				default:
					i.errc, i.expect = ErrUnexpToken, ExpectEscapedSequence
					goto ERROR
				}
				continue
			}
			if i.str[i.head] < 0x20 {
				i.errc, i.expect = ErrUnexpToken, ExpectEndOfString
				goto ERROR
			}
			i.head++
		}

	AFTER_STR_VAL:
		// Callback for argument
		i.token = TokenStr
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		// Advance head index to include the closing double-quotes
		i.head++
	/*</str>*/

	case '$':
		if inDefVal {
			i.errc, i.expect = ErrUnexpToken, ExpectDefaultVarVal
			goto ERROR
		}

		// Variable reference
		i.head++

		// Variable name
		i.expect = ExpectVarRefName
		goto VAR_REF_NAME

	case 'n':

		/*<null>*/
		if i.head+4 < len(i.str) &&
			lutEndVal[i.str[i.head+4]] == 1 &&
			i.str[i.head+3] == 'l' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'u' &&
			i.str[i.head] == 'n' {
			i.tail = -1
			i.head += len("null")

			// Callback for null value
			i.token = TokenNull
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</null>*/

	case 't':

		/*<true>*/
		if i.head+4 < len(i.str) &&
			lutEndVal[i.str[i.head+4]] == 1 &&
			i.str[i.head+3] == 'e' &&
			i.str[i.head+2] == 'u' &&
			i.str[i.head+1] == 'r' &&
			i.str[i.head] == 't' {
			i.tail = -1
			i.head += len("true")

			// Callback for true value
			i.token = TokenTrue
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</true>*/

	case 'f':

		/*<false>*/
		if i.head+5 < len(i.str) &&
			lutEndVal[i.str[i.head+5]] == 1 &&
			i.str[i.head+4] == 'e' &&
			i.str[i.head+3] == 's' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'a' &&
			i.str[i.head] == 'f' {
			i.tail = -1
			i.head += len("false")

			// Callback for false value
			i.token = TokenFalse
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</false>*/

	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

		/*<num>*/
		// Number
		i.tail = i.head

		var s int

		switch i.str[i.head] {
		case '-':
			// Signed
			i.head++

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectVal
				goto ERROR
			}
		/*</check_eof>*/

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
				} else if lutEndNum[i.str[i.head]] == 1 {
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
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if i.str[i.head] == '.' {
				i.head++
				goto FRACTION
			} else if lutEndNum[i.str[i.head]] == 1 {
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
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if lutEndNum[i.str[i.head]] == 1 {
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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectVal
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == '-' || i.str[i.head] == '+' {
			i.head++
		}
		for s = i.head; i.head < len(i.str); i.head++ {
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if lutEndNum[i.str[i.head]] == 1 {
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
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

	/*</callback>*/
	/*</num>*/

	default:
		// Invalid value
		i.expect = ExpectValEnum

		/*<name>*/
		// Followed by valenum>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectValEnum after name>
		i.token = TokenEnumVal
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.expect = ExpectAfterValueInner
		goto AFTER_VALUE_INNER
		// </ExpectValEnum after name>

		/*</name>*/

	}
	i.expect = ExpectAfterValueInner
	goto AFTER_VALUE_INNER
	/*</l_value>*/

	/*<l_block_string>*/
BLOCK_STRING:
	i.expect = ExpectEndOfBlockString
	for {
		for i.head+15 < len(i.str) {
			if lutStr[i.str[i.head]] == 1 {
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+1]] == 1 {
				i.head += 1
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+2]] == 1 {
				i.head += 2
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+3]] == 1 {
				i.head += 3
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+4]] == 1 {
				i.head += 4
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+5]] == 1 {
				i.head += 5
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+6]] == 1 {
				i.head += 6
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+7]] == 1 {
				i.head += 7
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+8]] == 1 {
				i.head += 8
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+9]] == 1 {
				i.head += 9
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+10]] == 1 {
				i.head += 10
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+11]] == 1 {
				i.head += 11
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+12]] == 1 {
				i.head += 12
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+13]] == 1 {
				i.head += 13
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+14]] == 1 {
				i.head += 14
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+15]] == 1 {
				i.head += 15
				goto CHECK_BLOCKSTR_CHARACTER
			}
			i.head += 16
		}

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

	CHECK_BLOCKSTR_CHARACTER:
		if i.str[i.head] == '\\' &&
			i.head+3 < len(i.str) &&
			i.str[i.head+3] == '"' &&
			i.str[i.head+2] == '"' &&
			i.str[i.head+1] == '"' {
			i.head += len(`\"""`)
			continue
		} else if i.str[i.head] == '"' &&
			i.head+2 < len(i.str) &&
			i.str[i.head+2] == '"' &&
			i.str[i.head+1] == '"' {
			i.token = TokenStrBlock
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head += len(`"""`)
			goto AFTER_VALUE_INNER
		} else if i.str[i.head] < 0x20 &&
			i.str[i.head] != '\t' &&
			i.str[i.head] != '\n' &&
			i.str[i.head] != '\r' {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
	}
	/*</l_block_string>*/

	/*<l_after_value_inner>*/
AFTER_VALUE_INNER:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}
	if t := i.stackTop(); t == TokenObj {
		if i.str[i.head] == '}' {
			i.tail = -1
			i.stackPop()

			// Callback for end of object
			i.token = TokenObjEnd
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			if i.stackLen() > 0 {
				i.expect = ExpectAfterValueInner
				goto AFTER_VALUE_INNER
			}
		} else {
			// Proceed to next field in the object
			i.expect = ExpectObjFieldName

			/*<name>*/
			// Followed by objfieldname>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectObjFieldName after name>
			i.token = TokenObjField
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectColObjFieldName
				goto ERROR
			}
			/*</check_eof>*/

			if i.str[i.head] != ':' {
				i.errc = ErrUnexpToken
				i.expect = ExpectColObjFieldName
				goto ERROR
			}
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			i.expect = ExpectVal
			goto VALUE
			// </ExpectObjFieldName after name>

			/*</name>*/

		}
	} else if t == TokenArr {
		if i.str[i.head] == ']' {
			i.tail = -1
			i.stackPop()

			// Callback for end of array
			i.token = TokenArrEnd
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			if i.stackLen() > 0 {
				i.expect = ExpectAfterValueInner
				goto AFTER_VALUE_INNER
			}
		} else {
			// Proceed to next value in the array
			goto VALUE
		}
	}
	goto AFTER_VALUE_OUTER
	/*</l_after_value_inner>*/

	/*<l_after_value_outer>*/
AFTER_VALUE_OUTER:

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if inDefVal {
		switch i.str[i.head] {
		case ')':
			inDefVal = false
			goto VAR_LIST_END
		case '@':
			inDefVal = false
			i.head++
			dirOn, i.expect = dirVar, ExpectDir
			goto DIR_NAME
		case '#':
			goto COMMENT
		}
		inDefVal = false
		i.expect = ExpectVar
		goto OPR_VAR
	}

	if i.str[i.head] == ')' {
		// End of argument list
		i.tail = -1
		i.token = TokenArgListEnd
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++
		i.expect = ExpectAfterArgList
		goto AFTER_ARG_LIST
	}

	// Proceed to the next argument
	i.expect = ExpectArgName

	/*<name>*/
	// Followed by argname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectArgName after name>
	i.token = TokenArgName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectColumnAfterArg
	goto COLUMN_AFTER_ARG_NAME
	// </ExpectArgName after name>

	/*</name>*/

	/*</l_after_value_outer>*/

	/*<l_after_arg_list>*/
AFTER_ARG_LIST:
	if dirOn != 0 {
		goto AFTER_DIR_ARGS
	}

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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
	/*</l_after_arg_list>*/

	/*<l_selection>*/
SELECTION:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSel
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		i.expect = ExpectSel
		goto COMMENT
	} else if i.str[i.head] != '.' {
		// Field selection
		i.expect = ExpectFieldNameOrAlias

		/*<name>*/
		// Followed by fieldnameoralias>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectFieldNameOrAlias after name>
		head := i.head

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == ':' {
			h2 := i.head
			i.head = head
			i.token = TokenFieldAlias
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head = h2 + 1

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			i.expect = ExpectFieldName

			/*<name>*/
			// Followed by fieldname>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectFieldName after name>
			i.token = TokenField
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			goto AFTER_FIELD_NAME
			// </ExpectFieldName after name>

			/*</name>*/

		}
		i.head = head
		i.token = TokenField
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		goto AFTER_FIELD_NAME
		// </ExpectFieldNameOrAlias after name>

		/*</name>*/

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
	/*</l_selection>*/

	/*<l_spread>*/
SPREAD:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token, i.tail = TokenFragInline, -1
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '@' {
		i.token, i.tail = TokenFragInline, -1
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
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
			goto FRAG_INLINED
		}
	}
	// ...fragmentName
	i.expect = ExpectSpreadName

	/*<name>*/
	// Followed by spreadname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectSpreadName after name>
	i.token = TokenNamedSpread
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect, dirOn = ExpectDirName, dirFragRef
	goto AFTER_DIR_NAME
	// </ExpectSpreadName after name>

	/*</name>*/

	/*</l_spread>*/

	/*<l_after_decl_varname>*/
AFTER_DECL_VAR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.expect = ExpectVarType
	goto VAR_TYPE
	/*</l_after_decl_varname>*/

	/*<l_var_type>*/
VAR_TYPE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '[' {
		i.tail = -1
		i.token = TokenVarTypeArr
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++
		typeArrLvl++
		goto VAR_TYPE
	}
	i.expect = ExpectVarType

	/*<name>*/
	// Followed by vartype>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarType after name>
	i.token = TokenVarTypeName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect = ExpectAfterVarTypeName
	goto AFTER_VAR_TYPE_NAME
	// </ExpectVarType after name>

	/*</name>*/

	/*</l_var_type>*/

	/*<l_var_name>*/
VAR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by varname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarName after name>
	i.token = TokenVarName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect = ExpectColumnAfterVar
	goto AFTER_DECL_VAR_NAME
	// </ExpectVarName after name>

	/*</name>*/

	/*</l_var_name>*/

	/*<l_var_ref>*/
VAR_REF_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by varrefname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarRefName after name>
	i.token = TokenVarRef
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect = ExpectAfterValueInner
	goto AFTER_VALUE_INNER
	// </ExpectVarRefName after name>

	/*</name>*/

	/*</l_var_ref>*/

	/*<l_dir_name>*/
DIR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}
	i.expect = ExpectDirName

	/*<name>*/
	// Followed by dirname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectDirName after name>
	i.token = TokenDirName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	goto AFTER_DIR_NAME
	// </ExpectDirName after name>

	/*</name>*/

	/*</l_dir_name>*/

	/*<l_column_after_arg_name>*/
COLUMN_AFTER_ARG_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.stackReset()
	i.expect = ExpectVal
	goto VALUE
	/*</l_column_after_arg_name>*/

	/*<l_arg_list>*/
ARG_LIST:

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by argname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectArgName after name>
	i.token = TokenArgName
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectColumnAfterArg
	goto COLUMN_AFTER_ARG_NAME
	// </ExpectArgName after name>

	/*</name>*/

	/*</l_arg_list>*/

	/*<l_after_var_type_name>*/
AFTER_VAR_TYPE_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head < len(i.str) && i.str[i.head] == '!' {
		i.tail = -1
		i.token = TokenVarTypeNotNull
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++
	}
	goto AFTER_VAR_TYPE_NOT_NULL
	/*</l_after_var_type_name>*/

	/*<l_after_var_type_not_null>*/
AFTER_VAR_TYPE_NOT_NULL:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == ']' {
		if typeArrLvl < 1 {
			i.errc, i.expect = ErrUnexpToken, ExpectVar
			goto ERROR
		}
		i.tail = -1
		i.token = TokenVarTypeArrEnd
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++
		typeArrLvl--

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		if i.head < len(i.str) && i.str[i.head] == '!' {
			i.tail = -1
			i.token = TokenVarTypeNotNull
			/*<callback>*/

			if fn(i) {
				i.errc = ErrCallbackFn
				goto ERROR
			}

			/*</callback>*/
			i.head++
		}

		if typeArrLvl > 0 {
			goto AFTER_VAR_TYPE_NAME
		}
	}
	i.expect = ExpectAfterVarType
	goto AFTER_VAR_TYPE
	/*</l_after_var_type_not_null>*/

	/*<l_after_field_name>*/
AFTER_FIELD_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	// Lookahead
	switch i.str[i.head] {
	case '(':
		// Argument list
		i.tail = -1
		i.token = TokenArgList
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

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
	/*</l_after_field_name>*/

	/*<l_after_opr_name>*/
AFTER_OPR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
		goto ERROR
	}
	/*</check_eof>*/

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
		/*<callback>*/

		if fn(i) {
			i.errc = ErrCallbackFn
			goto ERROR
		}

		/*</callback>*/
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
	/*</l_after_opr_name>*/

	/*<l_frag_keyword_on>*/
FRAG_KEYWORD_ON:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

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

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fragtypecond>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragTypeCond after name>
	i.token = TokenFragTypeCond
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '@' {
		dirOn = dirFragInlineOrDef
		goto AFTER_DIR_NAME
	}
	i.expect = ExpectSelSet
	goto SELECTION_SET
	// </ExpectFragTypeCond after name>

	/*</name>*/

	/*</l_frag_keyword_on>*/

	/*<l_frag_inlined>*/
FRAG_INLINED:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fraginlined>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragInlined after name>
	i.token = TokenFragInline
	/*<callback>*/

	if fn(i) {
		i.errc = ErrCallbackFn
		goto ERROR
	}

	/*</callback>*/
	i.expect, dirOn = ExpectDirName, dirFragInlineOrDef
	goto AFTER_DIR_NAME
	// </ExpectFragInlined after name>

	/*</name>*/

	/*</l_frag_inlined>*/

	/*<l_comment>*/
COMMENT:
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str) && i.str[i.head] != '\n'; i.head++ {
			}
			break
		}
		if i.str[i.head] == '\n' {
			break
		}
		if i.str[i.head+1] == '\n' {
			i.head += 1
			break
		}
		if i.str[i.head+2] == '\n' {
			i.head += 2
			break
		}
		if i.str[i.head+3] == '\n' {
			i.head += 3
			break
		}
		if i.str[i.head+4] == '\n' {
			i.head += 4
			break
		}
		if i.str[i.head+5] == '\n' {
			i.head += 5
			break
		}
		if i.str[i.head+6] == '\n' {
			i.head += 6
			break
		}
		if i.str[i.head+7] == '\n' {
			i.head += 7
			break
		}
		if i.str[i.head+8] == '\n' {
			i.head += 8
			break
		}
		if i.str[i.head+9] == '\n' {
			i.head += 9
			break
		}
		if i.str[i.head+10] == '\n' {
			i.head += 10
			break
		}
		if i.str[i.head+11] == '\n' {
			i.head += 11
			break
		}
		if i.str[i.head+12] == '\n' {
			i.head += 12
			break
		}
		if i.str[i.head+13] == '\n' {
			i.head += 13
			break
		}
		if i.str[i.head+14] == '\n' {
			i.head += 14
			break
		}
		if i.str[i.head+15] == '\n' {
			i.head += 15
			break
		}
		i.head += 16
	}
	i.tail = -1

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch i.expect {
	case ExpectOprName:
		goto AFTER_OPR_NAME
	case ExpectVarRefName:
		goto VAR_REF_NAME
	case ExpectVarName:
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
	case ExpectAfterValueInner:
		goto AFTER_VALUE_INNER
	case ExpectAfterValueOuter:
		goto AFTER_VALUE_OUTER
	case ExpectAfterArgList:
		goto AFTER_ARG_LIST
	case ExpectAfterDefKeyword:
		goto AFTER_DEF_KEYWORD
	case ExpectFragName:
		goto AFTER_KEYWORD_FRAGMENT
	case ExpectFragKeywordOn:
		goto FRAG_KEYWORD_ON
	case ExpectFragInlined:
		goto FRAG_INLINED
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
	/*</l_comment>*/

	/*<l_definition_end>*/
DEFINITION_END:
	i.levelSel, i.expect = 0, ExpectDef
	// Expect end of file

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head < len(i.str) {
		goto DEFINITION
	}
	return Error{}
	/*</l_definition_end>*/

	/*<l_error>*/
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
	/*</l_error>*/

	/*</scan_body>*/

}

// ScanAll calls fn for every token it scans in str.
// If the returned error code == 0 then there was no error during the scan,
// this can also be checked using err.IsErr().
//
// WARNING: *Iterator passed to fn should never be aliased and
// used after ScanAll returns because it's returned to the pool
// and may be acquired by another call to ScanAll!
func ScanAll(str []byte, fn func(*Iterator)) Error {

	/*<scan_body>*/
	i := iteratorPool.Get().(*Iterator)
	i.stackReset()
	i.expect = ExpectDef
	i.tail, i.head = -1, 0
	i.str = str
	i.levelSel = 0
	i.errc = 0
	defer iteratorPool.Put(i)

	// inDefVal triggers different expectations after values
	// when the iterator is in a variable default value definition.
	var inDefVal bool
	var typeArrLvl int
	var dirOn dirTarget

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectDef
		goto ERROR
	}
	/*</check_eof>*/

	/*<l_definition>*/
DEFINITION:
	if i.head >= len(i.str) {
		goto DEFINITION_END
	} else if i.str[i.head] == '#' {
		i.expect = ExpectDef
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token = TokenDefQry
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.isHeadKeywordQuery() {
		// Query
		i.token = TokenDefQry
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head += len("query")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordMutation() {
		// Mutation
		i.token = TokenDefMut
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head += len("mutation")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordSubscription() {
		// Subscription
		i.token = TokenDefSub
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head += len("subscription")
		i.expect = ExpectAfterDefKeyword
		goto AFTER_DEF_KEYWORD
	} else if i.isHeadKeywordFragment() {
		// Fragment
		i.tail = -1
		i.token = TokenDefFrag
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head += len("fragment")
		i.expect = ExpectFragName
		goto AFTER_KEYWORD_FRAGMENT
	}

	i.errc = ErrUnexpToken
	i.expect = ExpectDef
	goto ERROR
	/*</l_definition>*/

	/*<l_after_def_keyword>*/
AFTER_DEF_KEYWORD:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

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
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++
		i.expect = ExpectVar
		goto OPR_VAR
	case '@':
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	i.expect = ExpectOprName

	/*<name>*/
	// Followed by oprname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectOprName after name>
	i.token = TokenOprName
	/*<callback>*/

	fn(i)

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	goto AFTER_OPR_NAME
	// </ExpectOprName after name>

	/*</name>*/

	/*</l_after_def_keyword>*/

	/*<l_after_dir_name>*/
AFTER_DIR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch dirOn {
	case dirField:

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		/*</check_eof>*/

		switch i.str[i.head] {
		case '#':
			goto COMMENT
		case '(':
			// Directive argument list
			i.tail = -1
			i.token = TokenArgList
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

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
	/*</l_after_dir_name>*/

	/*<l_after_dir_args>*/
AFTER_DIR_ARGS:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch dirOn {
	case dirField:

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterFieldName
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterDefKeyword
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterVarType
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectAfterSelection
			goto ERROR
		}
		/*</check_eof>*/

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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
			goto ERROR
		}
		/*</check_eof>*/

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
	/*</l_after_dir_args>*/

	/*<l_after_keyword_fragment>*/
AFTER_KEYWORD_FRAGMENT:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fragname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragName after name>
	if i.head-i.tail == 2 &&
		i.str[i.tail+1] == 'n' &&
		i.str[i.tail] == 'o' {
		i.errc, i.head = ErrIllegalFragName, i.tail
		goto ERROR
	}
	i.token = TokenFragName
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect = ExpectFragKeywordOn
	goto FRAG_KEYWORD_ON
	// </ExpectFragName after name>

	/*</name>*/

	/*</l_after_keyword_fragment>*/

	/*<l_opr_var>*/
OPR_VAR:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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
	/*</l_opr_var>*/

	/*<l_after_var_type>*/
AFTER_VAR_TYPE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect, inDefVal = ExpectVal, true
		goto VALUE
	} else if i.str[i.head] == ')' {
		goto VAR_LIST_END
	}
	i.expect = ExpectAfterVarType
	goto OPR_VAR
	/*</l_after_var_type>*/

	/*<l_var_list_end>*/
VAR_LIST_END:
	i.tail = -1
	i.token = TokenVarListEnd
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.head++

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectSelSet

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		dirOn, i.expect = dirOpr, ExpectDirName
		goto AFTER_DIR_NAME
	} else if i.str[i.head] == '@' {
		i.head++
		dirOn, i.expect = dirOpr, ExpectDir
		goto DIR_NAME
	}
	goto SELECTION_SET
	/*</l_var_list_end>*/

	/*<l_selection_set>*/
SELECTION_SET:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != '{' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.tail = -1
	i.token = TokenSet
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.levelSel++
	i.head++
	i.expect = ExpectSel
	goto SELECTION
	/*</l_selection_set>*/

	/*<l_after_selection>*/
AFTER_SELECTION:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '}' {
		goto SEL_END
	}
	i.expect = ExpectSel
	goto SELECTION
	/*</l_after_selection>*/

	/*<l_sel_end>*/
SEL_END:
	i.tail = -1
	i.token = TokenSetEnd
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.levelSel--
	i.head++

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.levelSel < 1 {
		goto DEFINITION_END
	}
	goto AFTER_SELECTION
	/*</l_sel_end>*/

	/*<l_value>*/
VALUE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	switch i.str[i.head] {
	case '#':
		goto COMMENT

	case '{':
		// Object begin
		i.tail = -1
		// Callback for argument
		i.token = TokenObj
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.stackPush(TokenObj)
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect = ExpectObjFieldName

		/*<name>*/
		// Followed by objfieldname>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectObjFieldName after name>
		i.token = TokenObjField
		/*<callback>*/

		fn(i)

		/*</callback>*/

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectColObjFieldName
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] != ':' {
			i.errc = ErrUnexpToken
			i.expect = ExpectColObjFieldName
			goto ERROR
		}
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		i.expect = ExpectVal
		goto VALUE
	// </ExpectObjFieldName after name>

	/*</name>*/

	case '[':
		i.tail = -1
		// Callback for argument
		i.token = TokenArr
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		// Lookahead

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectVal
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == ']' {
			i.token = TokenArrEnd
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
		}
		i.stackPush(TokenArr)
		i.expect = ExpectAfterValueInner
		goto AFTER_VALUE_INNER

	case '"':

		/*<str>*/
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
		if i.head < len(i.str) && i.str[i.head] == '"' {
			goto AFTER_STR_VAL
		}
		for {
			if i.head+15 < len(i.str) {
				// Fast path
				if lutStr[i.str[i.head]] == 1 {
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+1]] == 1 {
					i.head += 1
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+2]] == 1 {
					i.head += 2
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+3]] == 1 {
					i.head += 3
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+4]] == 1 {
					i.head += 4
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+5]] == 1 {
					i.head += 5
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+6]] == 1 {
					i.head += 6
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+7]] == 1 {
					i.head += 7
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+8]] == 1 {
					i.head += 8
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+9]] == 1 {
					i.head += 9
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+10]] == 1 {
					i.head += 10
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+11]] == 1 {
					i.head += 11
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+12]] == 1 {
					i.head += 12
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+13]] == 1 {
					i.head += 13
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+14]] == 1 {
					i.head += 14
					goto CHECK_STR_CHARACTER
				}
				if lutStr[i.str[i.head+15]] == 1 {
					i.head += 15
					goto CHECK_STR_CHARACTER
				}
				i.head += 16
			}
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectEndOfString
				goto ERROR
			}
		CHECK_STR_CHARACTER:
			switch i.str[i.head] {
			case '"':
				// i.head++
				goto AFTER_STR_VAL
			case '\\':
				i.head++
				if i.head >= len(i.str) {
					i.errc, i.expect = ErrUnexpEOF, ExpectEscapedSequence
					goto ERROR
				}
				switch i.str[i.head] {
				case '\\', '/', '"', 'b', 'f', 'r', 'n', 't':
					i.head++
				case 'u':
					if i.head+4 >= len(i.str) {
						switch {
						case i.head+1 < len(i.str) && lutHex[i.str[i.head+1]] != 2:
							i.head++
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+2 < len(i.str) && lutHex[i.str[i.head+2]] != 2:
							i.head += 2
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+3 < len(i.str) && lutHex[i.str[i.head+3]] != 2:
							i.head += 3
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						case i.head+4 < len(i.str) && lutHex[i.str[i.head+4]] != 2:
							i.head += 4
							i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
							goto ERROR
						}
						i.head = len(i.str)
						i.errc, i.expect = ErrUnexpEOF, ExpectEscapedUnicodeSequence
						goto ERROR
					}
					if lutHex[i.str[i.head+4]] != 2 ||
						lutHex[i.str[i.head+3]] != 2 ||
						lutHex[i.str[i.head+2]] != 2 ||
						lutHex[i.str[i.head+1]] != 2 {
						switch {
						case lutHex[i.str[i.head+1]] != 2:
							i.head++
						case lutHex[i.str[i.head+2]] != 2:
							i.head += 2
						case lutHex[i.str[i.head+3]] != 2:
							i.head += 3
						case lutHex[i.str[i.head+4]] != 2:
							i.head += 4
						}
						i.errc, i.expect = ErrUnexpToken, ExpectEscapedUnicodeSequence
						goto ERROR
					}
					i.head += 4
				default:
					i.errc, i.expect = ErrUnexpToken, ExpectEscapedSequence
					goto ERROR
				}
				continue
			}
			if i.str[i.head] < 0x20 {
				i.errc, i.expect = ErrUnexpToken, ExpectEndOfString
				goto ERROR
			}
			i.head++
		}

	AFTER_STR_VAL:
		// Callback for argument
		i.token = TokenStr
		/*<callback>*/

		fn(i)

		/*</callback>*/
		// Advance head index to include the closing double-quotes
		i.head++
	/*</str>*/

	case '$':
		if inDefVal {
			i.errc, i.expect = ErrUnexpToken, ExpectDefaultVarVal
			goto ERROR
		}

		// Variable reference
		i.head++

		// Variable name
		i.expect = ExpectVarRefName
		goto VAR_REF_NAME

	case 'n':

		/*<null>*/
		if i.head+4 < len(i.str) &&
			lutEndVal[i.str[i.head+4]] == 1 &&
			i.str[i.head+3] == 'l' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'u' &&
			i.str[i.head] == 'n' {
			i.tail = -1
			i.head += len("null")

			// Callback for null value
			i.token = TokenNull
			/*<callback>*/

			fn(i)

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</null>*/

	case 't':

		/*<true>*/
		if i.head+4 < len(i.str) &&
			lutEndVal[i.str[i.head+4]] == 1 &&
			i.str[i.head+3] == 'e' &&
			i.str[i.head+2] == 'u' &&
			i.str[i.head+1] == 'r' &&
			i.str[i.head] == 't' {
			i.tail = -1
			i.head += len("true")

			// Callback for true value
			i.token = TokenTrue
			/*<callback>*/

			fn(i)

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</true>*/

	case 'f':

		/*<false>*/
		if i.head+5 < len(i.str) &&
			lutEndVal[i.str[i.head+5]] == 1 &&
			i.str[i.head+4] == 'e' &&
			i.str[i.head+3] == 's' &&
			i.str[i.head+2] == 'l' &&
			i.str[i.head+1] == 'a' &&
			i.str[i.head] == 'f' {
			i.tail = -1
			i.head += len("false")

			// Callback for false value
			i.token = TokenFalse
			/*<callback>*/

			fn(i)

			/*</callback>*/
		} else {
			i.expect = ExpectValEnum

			/*<name>*/
			// Followed by valenum>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectValEnum after name>
			i.token = TokenEnumVal
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.expect = ExpectAfterValueInner
			goto AFTER_VALUE_INNER
			// </ExpectValEnum after name>

			/*</name>*/

		}
	/*</false>*/

	case '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':

		/*<num>*/
		// Number
		i.tail = i.head

		var s int

		switch i.str[i.head] {
		case '-':
			// Signed
			i.head++

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectVal
				goto ERROR
			}
		/*</check_eof>*/

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
				} else if lutEndNum[i.str[i.head]] == 1 {
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
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if i.str[i.head] == '.' {
				i.head++
				goto FRACTION
			} else if lutEndNum[i.str[i.head]] == 1 {
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
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if lutEndNum[i.str[i.head]] == 1 {
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

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc, i.expect = ErrUnexpEOF, ExpectVal
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == '-' || i.str[i.head] == '+' {
			i.head++
		}
		for s = i.head; i.head < len(i.str); i.head++ {
			if lutName[i.str[i.head]] == 2 {
				continue
			} else if lutEndNum[i.str[i.head]] == 1 {
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
		/*<callback>*/

		fn(i)

	/*</callback>*/
	/*</num>*/

	default:
		// Invalid value
		i.expect = ExpectValEnum

		/*<name>*/
		// Followed by valenum>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectValEnum after name>
		i.token = TokenEnumVal
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.expect = ExpectAfterValueInner
		goto AFTER_VALUE_INNER
		// </ExpectValEnum after name>

		/*</name>*/

	}
	i.expect = ExpectAfterValueInner
	goto AFTER_VALUE_INNER
	/*</l_value>*/

	/*<l_block_string>*/
BLOCK_STRING:
	i.expect = ExpectEndOfBlockString
	for {
		for i.head+15 < len(i.str) {
			if lutStr[i.str[i.head]] == 1 {
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+1]] == 1 {
				i.head += 1
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+2]] == 1 {
				i.head += 2
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+3]] == 1 {
				i.head += 3
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+4]] == 1 {
				i.head += 4
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+5]] == 1 {
				i.head += 5
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+6]] == 1 {
				i.head += 6
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+7]] == 1 {
				i.head += 7
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+8]] == 1 {
				i.head += 8
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+9]] == 1 {
				i.head += 9
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+10]] == 1 {
				i.head += 10
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+11]] == 1 {
				i.head += 11
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+12]] == 1 {
				i.head += 12
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+13]] == 1 {
				i.head += 13
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+14]] == 1 {
				i.head += 14
				goto CHECK_BLOCKSTR_CHARACTER
			}
			if lutStr[i.str[i.head+15]] == 1 {
				i.head += 15
				goto CHECK_BLOCKSTR_CHARACTER
			}
			i.head += 16
		}

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

	CHECK_BLOCKSTR_CHARACTER:
		if i.str[i.head] == '\\' &&
			i.head+3 < len(i.str) &&
			i.str[i.head+3] == '"' &&
			i.str[i.head+2] == '"' &&
			i.str[i.head+1] == '"' {
			i.head += len(`\"""`)
			continue
		} else if i.str[i.head] == '"' &&
			i.head+2 < len(i.str) &&
			i.str[i.head+2] == '"' &&
			i.str[i.head+1] == '"' {
			i.token = TokenStrBlock
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head += len(`"""`)
			goto AFTER_VALUE_INNER
		} else if i.str[i.head] < 0x20 &&
			i.str[i.head] != '\t' &&
			i.str[i.head] != '\n' &&
			i.str[i.head] != '\r' {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
	}
	/*</l_block_string>*/

	/*<l_after_value_inner>*/
AFTER_VALUE_INNER:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}
	if t := i.stackTop(); t == TokenObj {
		if i.str[i.head] == '}' {
			i.tail = -1
			i.stackPop()

			// Callback for end of object
			i.token = TokenObjEnd
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			if i.stackLen() > 0 {
				i.expect = ExpectAfterValueInner
				goto AFTER_VALUE_INNER
			}
		} else {
			// Proceed to next field in the object
			i.expect = ExpectObjFieldName

			/*<name>*/
			// Followed by objfieldname>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectObjFieldName after name>
			i.token = TokenObjField
			/*<callback>*/

			fn(i)

			/*</callback>*/

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc, i.expect = ErrUnexpEOF, ExpectColObjFieldName
				goto ERROR
			}
			/*</check_eof>*/

			if i.str[i.head] != ':' {
				i.errc = ErrUnexpToken
				i.expect = ExpectColObjFieldName
				goto ERROR
			}
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			i.expect = ExpectVal
			goto VALUE
			// </ExpectObjFieldName after name>

			/*</name>*/

		}
	} else if t == TokenArr {
		if i.str[i.head] == ']' {
			i.tail = -1
			i.stackPop()

			// Callback for end of array
			i.token = TokenArrEnd
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			if i.stackLen() > 0 {
				i.expect = ExpectAfterValueInner
				goto AFTER_VALUE_INNER
			}
		} else {
			// Proceed to next value in the array
			goto VALUE
		}
	}
	goto AFTER_VALUE_OUTER
	/*</l_after_value_inner>*/

	/*<l_after_value_outer>*/
AFTER_VALUE_OUTER:

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if inDefVal {
		switch i.str[i.head] {
		case ')':
			inDefVal = false
			goto VAR_LIST_END
		case '@':
			inDefVal = false
			i.head++
			dirOn, i.expect = dirVar, ExpectDir
			goto DIR_NAME
		case '#':
			goto COMMENT
		}
		inDefVal = false
		i.expect = ExpectVar
		goto OPR_VAR
	}

	if i.str[i.head] == ')' {
		// End of argument list
		i.tail = -1
		i.token = TokenArgListEnd
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++
		i.expect = ExpectAfterArgList
		goto AFTER_ARG_LIST
	}

	// Proceed to the next argument
	i.expect = ExpectArgName

	/*<name>*/
	// Followed by argname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectArgName after name>
	i.token = TokenArgName
	/*<callback>*/

	fn(i)

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectColumnAfterArg
	goto COLUMN_AFTER_ARG_NAME
	// </ExpectArgName after name>

	/*</name>*/

	/*</l_after_value_outer>*/

	/*<l_after_arg_list>*/
AFTER_ARG_LIST:
	if dirOn != 0 {
		goto AFTER_DIR_ARGS
	}

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
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
	/*</l_after_arg_list>*/

	/*<l_selection>*/
SELECTION:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSel
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		i.expect = ExpectSel
		goto COMMENT
	} else if i.str[i.head] != '.' {
		// Field selection
		i.expect = ExpectFieldNameOrAlias

		/*<name>*/
		// Followed by fieldnameoralias>

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		i.tail = i.head
		if lutName[i.str[i.head]] != 1 {
			i.errc = ErrUnexpToken
			goto ERROR
		}
		i.head++
		for {
			if i.head+15 >= len(i.str) {
				for ; i.head < len(i.str); i.head++ {
					if lutName[i.str[i.head]] > 2 {
						if lutName[i.str[i.head]] == 5 {
							i.errc = ErrUnexpToken
							goto ERROR
						}
						break
					}
				}
				break
			}
			if lutName[i.str[i.head]] > 2 {
				break
			}
			if lutName[i.str[i.head+1]] > 2 {
				i.head += 1
				break
			}
			if lutName[i.str[i.head+2]] > 2 {
				i.head += 2
				break
			}
			if lutName[i.str[i.head+3]] > 2 {
				i.head += 3
				break
			}
			if lutName[i.str[i.head+4]] > 2 {
				i.head += 4
				break
			}
			if lutName[i.str[i.head+5]] > 2 {
				i.head += 5
				break
			}
			if lutName[i.str[i.head+6]] > 2 {
				i.head += 6
				break
			}
			if lutName[i.str[i.head+7]] > 2 {
				i.head += 7
				break
			}
			if lutName[i.str[i.head+8]] > 2 {
				i.head += 8
				break
			}
			if lutName[i.str[i.head+9]] > 2 {
				i.head += 9
				break
			}
			if lutName[i.str[i.head+10]] > 2 {
				i.head += 10
				break
			}
			if lutName[i.str[i.head+11]] > 2 {
				i.head += 11
				break
			}
			if lutName[i.str[i.head+12]] > 2 {
				i.head += 12
				break
			}
			if lutName[i.str[i.head+13]] > 2 {
				i.head += 13
				break
			}
			if lutName[i.str[i.head+14]] > 2 {
				i.head += 14
				break
			}
			if lutName[i.str[i.head+15]] > 2 {
				i.head += 15
				break
			}
			i.head += 16
		}

		// <ExpectFieldNameOrAlias after name>
		head := i.head

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		/*<check_eof>*/
		if i.head >= len(i.str) {
			i.errc = ErrUnexpEOF
			goto ERROR
		}
		/*</check_eof>*/

		if i.str[i.head] == ':' {
			h2 := i.head
			i.head = head
			i.token = TokenFieldAlias
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head = h2 + 1

			/*<skip_irrelevant>*/
			for {
				if i.head+15 >= len(i.str) {
					for i.head < len(i.str) {
						if lutIrrel[i.str[i.head]] != 1 {
							break
						}
						i.head++
					}
					break
				}
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				if lutIrrel[i.str[i.head+1]] != 1 {
					i.head += 1
					break
				}
				if lutIrrel[i.str[i.head+2]] != 1 {
					i.head += 2
					break
				}
				if lutIrrel[i.str[i.head+3]] != 1 {
					i.head += 3
					break
				}
				if lutIrrel[i.str[i.head+4]] != 1 {
					i.head += 4
					break
				}
				if lutIrrel[i.str[i.head+5]] != 1 {
					i.head += 5
					break
				}
				if lutIrrel[i.str[i.head+6]] != 1 {
					i.head += 6
					break
				}
				if lutIrrel[i.str[i.head+7]] != 1 {
					i.head += 7
					break
				}
				if lutIrrel[i.str[i.head+8]] != 1 {
					i.head += 8
					break
				}
				if lutIrrel[i.str[i.head+9]] != 1 {
					i.head += 9
					break
				}
				if lutIrrel[i.str[i.head+10]] != 1 {
					i.head += 10
					break
				}
				if lutIrrel[i.str[i.head+11]] != 1 {
					i.head += 11
					break
				}
				if lutIrrel[i.str[i.head+12]] != 1 {
					i.head += 12
					break
				}
				if lutIrrel[i.str[i.head+13]] != 1 {
					i.head += 13
					break
				}
				if lutIrrel[i.str[i.head+14]] != 1 {
					i.head += 14
					break
				}
				if lutIrrel[i.str[i.head+15]] != 1 {
					i.head += 15
					break
				}
				i.head += 16
			}
			/*</skip_irrelevant>*/

			i.expect = ExpectFieldName

			/*<name>*/
			// Followed by fieldname>

			/*<check_eof>*/
			if i.head >= len(i.str) {
				i.errc = ErrUnexpEOF
				goto ERROR
			}
			/*</check_eof>*/

			i.tail = i.head
			if lutName[i.str[i.head]] != 1 {
				i.errc = ErrUnexpToken
				goto ERROR
			}
			i.head++
			for {
				if i.head+15 >= len(i.str) {
					for ; i.head < len(i.str); i.head++ {
						if lutName[i.str[i.head]] > 2 {
							if lutName[i.str[i.head]] == 5 {
								i.errc = ErrUnexpToken
								goto ERROR
							}
							break
						}
					}
					break
				}
				if lutName[i.str[i.head]] > 2 {
					break
				}
				if lutName[i.str[i.head+1]] > 2 {
					i.head += 1
					break
				}
				if lutName[i.str[i.head+2]] > 2 {
					i.head += 2
					break
				}
				if lutName[i.str[i.head+3]] > 2 {
					i.head += 3
					break
				}
				if lutName[i.str[i.head+4]] > 2 {
					i.head += 4
					break
				}
				if lutName[i.str[i.head+5]] > 2 {
					i.head += 5
					break
				}
				if lutName[i.str[i.head+6]] > 2 {
					i.head += 6
					break
				}
				if lutName[i.str[i.head+7]] > 2 {
					i.head += 7
					break
				}
				if lutName[i.str[i.head+8]] > 2 {
					i.head += 8
					break
				}
				if lutName[i.str[i.head+9]] > 2 {
					i.head += 9
					break
				}
				if lutName[i.str[i.head+10]] > 2 {
					i.head += 10
					break
				}
				if lutName[i.str[i.head+11]] > 2 {
					i.head += 11
					break
				}
				if lutName[i.str[i.head+12]] > 2 {
					i.head += 12
					break
				}
				if lutName[i.str[i.head+13]] > 2 {
					i.head += 13
					break
				}
				if lutName[i.str[i.head+14]] > 2 {
					i.head += 14
					break
				}
				if lutName[i.str[i.head+15]] > 2 {
					i.head += 15
					break
				}
				i.head += 16
			}

			// <ExpectFieldName after name>
			i.token = TokenField
			/*<callback>*/

			fn(i)

			/*</callback>*/
			goto AFTER_FIELD_NAME
			// </ExpectFieldName after name>

			/*</name>*/

		}
		i.head = head
		i.token = TokenField
		/*<callback>*/

		fn(i)

		/*</callback>*/
		goto AFTER_FIELD_NAME
		// </ExpectFieldNameOrAlias after name>

		/*</name>*/

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
	/*</l_selection>*/

	/*<l_spread>*/
SPREAD:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head+1 >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	} else if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '{' {
		i.token, i.tail = TokenFragInline, -1
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.expect = ExpectSelSet
		goto SELECTION_SET
	} else if i.str[i.head] == '@' {
		i.token, i.tail = TokenFragInline, -1
		/*<callback>*/

		fn(i)

		/*</callback>*/
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
			goto FRAG_INLINED
		}
	}
	// ...fragmentName
	i.expect = ExpectSpreadName

	/*<name>*/
	// Followed by spreadname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectSpreadName after name>
	i.token = TokenNamedSpread
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect, dirOn = ExpectDirName, dirFragRef
	goto AFTER_DIR_NAME
	// </ExpectSpreadName after name>

	/*</name>*/

	/*</l_spread>*/

	/*<l_after_decl_varname>*/
AFTER_DECL_VAR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.expect = ExpectVarType
	goto VAR_TYPE
	/*</l_after_decl_varname>*/

	/*<l_var_type>*/
VAR_TYPE:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == '[' {
		i.tail = -1
		i.token = TokenVarTypeArr
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++
		typeArrLvl++
		goto VAR_TYPE
	}
	i.expect = ExpectVarType

	/*<name>*/
	// Followed by vartype>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarType after name>
	i.token = TokenVarTypeName
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect = ExpectAfterVarTypeName
	goto AFTER_VAR_TYPE_NAME
	// </ExpectVarType after name>

	/*</name>*/

	/*</l_var_type>*/

	/*<l_var_name>*/
VAR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by varname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarName after name>
	i.token = TokenVarName
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect = ExpectColumnAfterVar
	goto AFTER_DECL_VAR_NAME
	// </ExpectVarName after name>

	/*</name>*/

	/*</l_var_name>*/

	/*<l_var_ref>*/
VAR_REF_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by varrefname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectVarRefName after name>
	i.token = TokenVarRef
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect = ExpectAfterValueInner
	goto AFTER_VALUE_INNER
	// </ExpectVarRefName after name>

	/*</name>*/

	/*</l_var_ref>*/

	/*<l_dir_name>*/
DIR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}
	i.expect = ExpectDirName

	/*<name>*/
	// Followed by dirname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectDirName after name>
	i.token = TokenDirName
	/*<callback>*/

	fn(i)

	/*</callback>*/
	goto AFTER_DIR_NAME
	// </ExpectDirName after name>

	/*</name>*/

	/*</l_dir_name>*/

	/*<l_column_after_arg_name>*/
COLUMN_AFTER_ARG_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] != ':' {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	i.stackReset()
	i.expect = ExpectVal
	goto VALUE
	/*</l_column_after_arg_name>*/

	/*<l_arg_list>*/
ARG_LIST:

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by argname>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectArgName after name>
	i.token = TokenArgName
	/*<callback>*/

	fn(i)

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	i.expect = ExpectColumnAfterArg
	goto COLUMN_AFTER_ARG_NAME
	// </ExpectArgName after name>

	/*</name>*/

	/*</l_arg_list>*/

	/*<l_after_var_type_name>*/
AFTER_VAR_TYPE_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head < len(i.str) && i.str[i.head] == '!' {
		i.tail = -1
		i.token = TokenVarTypeNotNull
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++
	}
	goto AFTER_VAR_TYPE_NOT_NULL
	/*</l_after_var_type_name>*/

	/*<l_after_var_type_not_null>*/
AFTER_VAR_TYPE_NOT_NULL:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	} else if i.str[i.head] == ']' {
		if typeArrLvl < 1 {
			i.errc, i.expect = ErrUnexpToken, ExpectVar
			goto ERROR
		}
		i.tail = -1
		i.token = TokenVarTypeArrEnd
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++
		typeArrLvl--

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

		if i.head < len(i.str) && i.str[i.head] == '!' {
			i.tail = -1
			i.token = TokenVarTypeNotNull
			/*<callback>*/

			fn(i)

			/*</callback>*/
			i.head++
		}

		if typeArrLvl > 0 {
			goto AFTER_VAR_TYPE_NAME
		}
	}
	i.expect = ExpectAfterVarType
	goto AFTER_VAR_TYPE
	/*</l_after_var_type_not_null>*/

	/*<l_after_field_name>*/
AFTER_FIELD_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	// Lookahead
	switch i.str[i.head] {
	case '(':
		// Argument list
		i.tail = -1
		i.token = TokenArgList
		/*<callback>*/

		fn(i)

		/*</callback>*/
		i.head++

		/*<skip_irrelevant>*/
		for {
			if i.head+15 >= len(i.str) {
				for i.head < len(i.str) {
					if lutIrrel[i.str[i.head]] != 1 {
						break
					}
					i.head++
				}
				break
			}
			if lutIrrel[i.str[i.head]] != 1 {
				break
			}
			if lutIrrel[i.str[i.head+1]] != 1 {
				i.head += 1
				break
			}
			if lutIrrel[i.str[i.head+2]] != 1 {
				i.head += 2
				break
			}
			if lutIrrel[i.str[i.head+3]] != 1 {
				i.head += 3
				break
			}
			if lutIrrel[i.str[i.head+4]] != 1 {
				i.head += 4
				break
			}
			if lutIrrel[i.str[i.head+5]] != 1 {
				i.head += 5
				break
			}
			if lutIrrel[i.str[i.head+6]] != 1 {
				i.head += 6
				break
			}
			if lutIrrel[i.str[i.head+7]] != 1 {
				i.head += 7
				break
			}
			if lutIrrel[i.str[i.head+8]] != 1 {
				i.head += 8
				break
			}
			if lutIrrel[i.str[i.head+9]] != 1 {
				i.head += 9
				break
			}
			if lutIrrel[i.str[i.head+10]] != 1 {
				i.head += 10
				break
			}
			if lutIrrel[i.str[i.head+11]] != 1 {
				i.head += 11
				break
			}
			if lutIrrel[i.str[i.head+12]] != 1 {
				i.head += 12
				break
			}
			if lutIrrel[i.str[i.head+13]] != 1 {
				i.head += 13
				break
			}
			if lutIrrel[i.str[i.head+14]] != 1 {
				i.head += 14
				break
			}
			if lutIrrel[i.str[i.head+15]] != 1 {
				i.head += 15
				break
			}
			i.head += 16
		}
		/*</skip_irrelevant>*/

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
	/*</l_after_field_name>*/

	/*<l_after_opr_name>*/
AFTER_OPR_NAME:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
		goto ERROR
	}
	/*</check_eof>*/

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
		/*<callback>*/

		fn(i)

		/*</callback>*/
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
	/*</l_after_opr_name>*/

	/*<l_frag_keyword_on>*/
FRAG_KEYWORD_ON:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

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

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fragtypecond>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragTypeCond after name>
	i.token = TokenFragTypeCond
	/*<callback>*/

	fn(i)

	/*</callback>*/

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc, i.expect = ErrUnexpEOF, ExpectSelSet
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '@' {
		dirOn = dirFragInlineOrDef
		goto AFTER_DIR_NAME
	}
	i.expect = ExpectSelSet
	goto SELECTION_SET
	// </ExpectFragTypeCond after name>

	/*</name>*/

	/*</l_frag_keyword_on>*/

	/*<l_frag_inlined>*/
FRAG_INLINED:

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	if i.str[i.head] == '#' {
		goto COMMENT
	}

	/*<name>*/
	// Followed by fraginlined>

	/*<check_eof>*/
	if i.head >= len(i.str) {
		i.errc = ErrUnexpEOF
		goto ERROR
	}
	/*</check_eof>*/

	i.tail = i.head
	if lutName[i.str[i.head]] != 1 {
		i.errc = ErrUnexpToken
		goto ERROR
	}
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str); i.head++ {
				if lutName[i.str[i.head]] > 2 {
					if lutName[i.str[i.head]] == 5 {
						i.errc = ErrUnexpToken
						goto ERROR
					}
					break
				}
			}
			break
		}
		if lutName[i.str[i.head]] > 2 {
			break
		}
		if lutName[i.str[i.head+1]] > 2 {
			i.head += 1
			break
		}
		if lutName[i.str[i.head+2]] > 2 {
			i.head += 2
			break
		}
		if lutName[i.str[i.head+3]] > 2 {
			i.head += 3
			break
		}
		if lutName[i.str[i.head+4]] > 2 {
			i.head += 4
			break
		}
		if lutName[i.str[i.head+5]] > 2 {
			i.head += 5
			break
		}
		if lutName[i.str[i.head+6]] > 2 {
			i.head += 6
			break
		}
		if lutName[i.str[i.head+7]] > 2 {
			i.head += 7
			break
		}
		if lutName[i.str[i.head+8]] > 2 {
			i.head += 8
			break
		}
		if lutName[i.str[i.head+9]] > 2 {
			i.head += 9
			break
		}
		if lutName[i.str[i.head+10]] > 2 {
			i.head += 10
			break
		}
		if lutName[i.str[i.head+11]] > 2 {
			i.head += 11
			break
		}
		if lutName[i.str[i.head+12]] > 2 {
			i.head += 12
			break
		}
		if lutName[i.str[i.head+13]] > 2 {
			i.head += 13
			break
		}
		if lutName[i.str[i.head+14]] > 2 {
			i.head += 14
			break
		}
		if lutName[i.str[i.head+15]] > 2 {
			i.head += 15
			break
		}
		i.head += 16
	}

	// <ExpectFragInlined after name>
	i.token = TokenFragInline
	/*<callback>*/

	fn(i)

	/*</callback>*/
	i.expect, dirOn = ExpectDirName, dirFragInlineOrDef
	goto AFTER_DIR_NAME
	// </ExpectFragInlined after name>

	/*</name>*/

	/*</l_frag_inlined>*/

	/*<l_comment>*/
COMMENT:
	i.head++
	for {
		if i.head+15 >= len(i.str) {
			for ; i.head < len(i.str) && i.str[i.head] != '\n'; i.head++ {
			}
			break
		}
		if i.str[i.head] == '\n' {
			break
		}
		if i.str[i.head+1] == '\n' {
			i.head += 1
			break
		}
		if i.str[i.head+2] == '\n' {
			i.head += 2
			break
		}
		if i.str[i.head+3] == '\n' {
			i.head += 3
			break
		}
		if i.str[i.head+4] == '\n' {
			i.head += 4
			break
		}
		if i.str[i.head+5] == '\n' {
			i.head += 5
			break
		}
		if i.str[i.head+6] == '\n' {
			i.head += 6
			break
		}
		if i.str[i.head+7] == '\n' {
			i.head += 7
			break
		}
		if i.str[i.head+8] == '\n' {
			i.head += 8
			break
		}
		if i.str[i.head+9] == '\n' {
			i.head += 9
			break
		}
		if i.str[i.head+10] == '\n' {
			i.head += 10
			break
		}
		if i.str[i.head+11] == '\n' {
			i.head += 11
			break
		}
		if i.str[i.head+12] == '\n' {
			i.head += 12
			break
		}
		if i.str[i.head+13] == '\n' {
			i.head += 13
			break
		}
		if i.str[i.head+14] == '\n' {
			i.head += 14
			break
		}
		if i.str[i.head+15] == '\n' {
			i.head += 15
			break
		}
		i.head += 16
	}
	i.tail = -1

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	switch i.expect {
	case ExpectOprName:
		goto AFTER_OPR_NAME
	case ExpectVarRefName:
		goto VAR_REF_NAME
	case ExpectVarName:
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
	case ExpectAfterValueInner:
		goto AFTER_VALUE_INNER
	case ExpectAfterValueOuter:
		goto AFTER_VALUE_OUTER
	case ExpectAfterArgList:
		goto AFTER_ARG_LIST
	case ExpectAfterDefKeyword:
		goto AFTER_DEF_KEYWORD
	case ExpectFragName:
		goto AFTER_KEYWORD_FRAGMENT
	case ExpectFragKeywordOn:
		goto FRAG_KEYWORD_ON
	case ExpectFragInlined:
		goto FRAG_INLINED
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
	/*</l_comment>*/

	/*<l_definition_end>*/
DEFINITION_END:
	i.levelSel, i.expect = 0, ExpectDef
	// Expect end of file

	/*<skip_irrelevant>*/
	for {
		if i.head+15 >= len(i.str) {
			for i.head < len(i.str) {
				if lutIrrel[i.str[i.head]] != 1 {
					break
				}
				i.head++
			}
			break
		}
		if lutIrrel[i.str[i.head]] != 1 {
			break
		}
		if lutIrrel[i.str[i.head+1]] != 1 {
			i.head += 1
			break
		}
		if lutIrrel[i.str[i.head+2]] != 1 {
			i.head += 2
			break
		}
		if lutIrrel[i.str[i.head+3]] != 1 {
			i.head += 3
			break
		}
		if lutIrrel[i.str[i.head+4]] != 1 {
			i.head += 4
			break
		}
		if lutIrrel[i.str[i.head+5]] != 1 {
			i.head += 5
			break
		}
		if lutIrrel[i.str[i.head+6]] != 1 {
			i.head += 6
			break
		}
		if lutIrrel[i.str[i.head+7]] != 1 {
			i.head += 7
			break
		}
		if lutIrrel[i.str[i.head+8]] != 1 {
			i.head += 8
			break
		}
		if lutIrrel[i.str[i.head+9]] != 1 {
			i.head += 9
			break
		}
		if lutIrrel[i.str[i.head+10]] != 1 {
			i.head += 10
			break
		}
		if lutIrrel[i.str[i.head+11]] != 1 {
			i.head += 11
			break
		}
		if lutIrrel[i.str[i.head+12]] != 1 {
			i.head += 12
			break
		}
		if lutIrrel[i.str[i.head+13]] != 1 {
			i.head += 13
			break
		}
		if lutIrrel[i.str[i.head+14]] != 1 {
			i.head += 14
			break
		}
		if lutIrrel[i.str[i.head+15]] != 1 {
			i.head += 15
			break
		}
		i.head += 16
	}
	/*</skip_irrelevant>*/

	if i.head < len(i.str) {
		goto DEFINITION
	}
	return Error{}
	/*</l_definition_end>*/

	/*<l_error>*/
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
	/*</l_error>*/

	/*</scan_body>*/

}

// Iterator is a GraphQL iterator for lexical analysis.
//
// WARNING: An iterator instance shall never be aliased and/or used
// after Scan or ScanAll returns because it's returned to a global pool!
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

// LevelSelect returns the current selector level.
func (i *Iterator) LevelSelect() int {
	return i.levelSel
}

// IndexHead returns the current head index.
func (i *Iterator) IndexHead() int {
	return i.head
}

// IndexTail returns the current tail index.
// Returns -1 if the current token doesn't reflect a dynamic value.
func (i *Iterator) IndexTail() int {
	return i.tail
}

// Token returns the current token type.
func (i *Iterator) Token() Token {
	return i.token
}

// Value returns the raw value of the current token.
// For TokenStrBlock it's the raw uninterpreted body of the string,
// use ScanInterpreted for the interpreted value of the block string.
//
// WARNING: The returned byte slice refers to the same underlying memory
// as the byte slice passed to Scan and ScanAll as str parameter,
// copy it or use with caution!
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

// isHeadKeywordSubscription returns true if
// the current head equals 'subscription'.
func (i *Iterator) isHeadKeywordSubscription() bool {
	return i.head+11 < len(i.str) &&
		i.str[i.head+11] == 'n' &&
		i.str[i.head+10] == 'o' &&
		i.str[i.head+9] == 'i' &&
		i.str[i.head+8] == 't' &&
		i.str[i.head+7] == 'p' &&
		i.str[i.head+6] == 'i' &&
		i.str[i.head+5] == 'r' &&
		i.str[i.head+4] == 'c' &&
		i.str[i.head+3] == 's' &&
		i.str[i.head+2] == 'b' &&
		i.str[i.head+1] == 'u' &&
		i.str[i.head] == 's'
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
	ExpectValEnum
	ExpectDefaultVarVal
	ExpectDef
	ExpectOprName
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
	ExpectDir
	ExpectDirName
	ExpectVar
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
	ExpectSpreadName
	ExpectFragInlined
	ExpectAfterFieldName
	ExpectAfterSelection
	ExpectAfterValueInner
	ExpectAfterValueOuter
	ExpectAfterArgList
	ExpectAfterDefKeyword
	ExpectAfterVarType
	ExpectAfterVarTypeName
)

func (e Expect) String() string {
	switch e {
	case ExpectDef:
		return "definition"
	case ExpectOprName:
		return "operation name"
	case ExpectVal:
		return "value"
	case ExpectValEnum:
		return "enum value"
	case ExpectDefaultVarVal:
		return "default variable value"
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
	case ExpectDir:
		return "directive name"
	case ExpectDirName:
		return "directive name"
	case ExpectVar:
		return "variable"
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
	case ExpectSpreadName:
		return "spread name"
	case ExpectFragInlined:
		return "inlined fragment"
	case ExpectAfterFieldName:
		return "selection, selection set or end of selection set"
	case ExpectAfterSelection:
		return "selection or end of selection set"
	case ExpectAfterValueInner:
		return "argument list closure or argument"
	case ExpectAfterValueOuter:
		return "argument list closure or argument"
	case ExpectAfterArgList:
		return "selection set or selection"
	case ExpectAfterDefKeyword:
		return "variable list or selection set"
	case ExpectAfterVarType:
		return "variable list closure or variable"
	case ExpectAfterVarTypeName:
		return "variable list closure or variable"
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
	TokenDefSub
	TokenDefFrag
	TokenOprName
	TokenDirName
	TokenVarList
	TokenVarListEnd
	TokenArgList
	TokenArgListEnd
	TokenSet
	TokenSetEnd
	TokenFragTypeCond
	TokenFragName
	TokenFragInline
	TokenNamedSpread
	TokenFieldAlias
	TokenField
	TokenArgName
	TokenEnumVal
	TokenArr
	TokenArrEnd
	TokenStr
	TokenStrBlock
	TokenInt
	TokenFloat
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
	case TokenDefSub:
		return "subscription definition"
	case TokenDefFrag:
		return "fragment definition"
	case TokenOprName:
		return "operation name"
	case TokenDirName:
		return "directive name"
	case TokenVarList:
		return "variable list"
	case TokenVarListEnd:
		return "variable list end"
	case TokenArgList:
		return "argument list"
	case TokenArgListEnd:
		return "argument list end"
	case TokenSet:
		return "selection set"
	case TokenSetEnd:
		return "selection set end"
	case TokenFragTypeCond:
		return "fragment type condition"
	case TokenFragName:
		return "fragment name"
	case TokenFragInline:
		return "fragment inline"
	case TokenNamedSpread:
		return "named spread"
	case TokenFieldAlias:
		return "field alias"
	case TokenField:
		return "field"
	case TokenArgName:
		return "argument name"
	case TokenEnumVal:
		return "enum value"
	case TokenArr:
		return "array"
	case TokenArrEnd:
		return "array end"
	case TokenStr:
		return "string"
	case TokenStrBlock:
		return "block string"
	case TokenInt:
		return "integer"
	case TokenFloat:
		return "float"
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
	ErrIllegalFragName
	ErrInvalNum
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
	case ErrIllegalFragName:
		b.WriteString(": illegal fragment name")
	case ErrInvalNum:
		b.WriteString(": invalid number value")
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

type dirTarget int

const (
	_ dirTarget = iota
	dirOpr
	dirVar
	dirField
	dirFragRef
	dirFragInlineOrDef
)

// lutStr maps 0 to all bytes that don't require checking during string traversal.
// 1 is mapped to control, quotation mark (") and reverse solidus ("\").
// All other characters are mapped to 0.
var lutStr = [256]byte{
	0x00: 1, 0x01: 1, 0x02: 1, 0x03: 1, 0x04: 1, 0x05: 1, 0x06: 1, 0x07: 1,
	0x08: 1, 0x0b: 1, 0x0c: 1, 0x0e: 1, 0x0f: 1, 0x10: 1, 0x11: 1, 0x12: 1,
	0x13: 1, 0x14: 1, 0x15: 1, 0x16: 1, 0x17: 1, 0x18: 1, 0x19: 1, 0x1a: 1,
	0x1b: 1, 0x1c: 1, 0x1d: 1, 0x1e: 1, 0x1f: 1,

	'\\': 1, '"': 1,
}

// lutHex maps valid hex digits to 2. All other characters are mapped to 0.
var lutHex = [256]byte{
	'0': 2, '1': 2, '2': 2, '3': 2, '4': 2, '5': 2, '6': 2, '7': 2, '8': 2, '9': 2,
	'a': 2, 'b': 2, 'c': 2, 'd': 2, 'e': 2, 'f': 2,
	'A': 2, 'B': 2, 'C': 2, 'D': 2, 'E': 2, 'F': 2,
}

// lutIrrel maps irrelevant characters such as comma, whitespace,
// line-break, tab and carriage-return to 1.
// All other characters are mapped to 0.
var lutIrrel = [256]byte{
	',': 1, ' ': 1, '\n': 1, '\t': 1, '\r': 1,
}

// lutEndVal maps characters that can appear at the
// end of a value to 1. All other characters are mapped to 0.
var lutEndVal = [256]byte{
	' ': 1, '\t': 1, '\r': 1, '\n': 1,
	',': 1, ')': 1, '}': 1, '{': 1, ']': 1, '[': 1, '#': 1,
}

// lutEndNum maps characters that can appear at the
// end of a number value to 1. All other characters are mapped to 0.
var lutEndNum = [256]byte{
	' ': 1, '\t': 1, '\r': 1, '\n': 1, ',': 1, ')': 1, '}': 1, ']': 1, '#': 1,
}

// lutName maps:
//
//	characters names can start with to 1,
//	characters names can contain to <3,
//	digit characters to 2.
//	characters names can end at to 3,
//	illegal control characters to 5,
//	all other characters to 4.
var lutName = [256]byte{
	'_': 1,

	'A': 1, 'B': 1, 'C': 1, 'D': 1, 'E': 1, 'F': 1, 'G': 1,
	'H': 1, 'I': 1, 'J': 1, 'K': 1, 'L': 1, 'M': 1, 'N': 1,
	'O': 1, 'P': 1, 'Q': 1, 'R': 1, 'S': 1, 'T': 1, 'U': 1,
	'V': 1, 'W': 1, 'X': 1, 'Y': 1, 'Z': 1,

	'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1, 'g': 1,
	'h': 1, 'i': 1, 'j': 1, 'k': 1, 'l': 1, 'm': 1, 'n': 1,
	'o': 1, 'p': 1, 'q': 1, 'r': 1, 's': 1, 't': 1, 'u': 1,
	'v': 1, 'w': 1, 'x': 1, 'y': 1, 'z': 1,

	'0': 2, '1': 2, '2': 2, '3': 2, '4': 2, '5': 2,
	'6': 2, '7': 2, '8': 2, '9': 2,

	',': 3, ' ': 3, '\n': 3, '\t': 3, '\r': 3,

	0x21: 4, 0x22: 4, 0x23: 4, 0x24: 4, 0x25: 4, 0x26: 4, 0x27: 4, 0x28: 4,
	0x29: 4, 0x2a: 4, 0x2b: 4, 0x2d: 4, 0x2e: 4, 0x2f: 4, 0x3a: 4, 0x3b: 4,
	0x3c: 4, 0x3d: 4, 0x3e: 4, 0x3f: 4, 0x40: 4, 0x5b: 4, 0x5c: 4, 0x5d: 4,
	0x5e: 4, 0x60: 4, 0x7b: 4, 0x7c: 4, 0x7d: 4, 0x7e: 4, 0x7f: 4, 0x80: 4,
	0x81: 4, 0x82: 4, 0x83: 4, 0x84: 4, 0x85: 4, 0x86: 4, 0x87: 4, 0x88: 4,
	0x89: 4, 0x8a: 4, 0x8b: 4, 0x8c: 4, 0x8d: 4, 0x8e: 4, 0x8f: 4, 0x90: 4,
	0x91: 4, 0x92: 4, 0x93: 4, 0x94: 4, 0x95: 4, 0x96: 4, 0x97: 4, 0x98: 4,
	0x99: 4, 0x9a: 4, 0x9b: 4, 0x9c: 4, 0x9d: 4, 0x9e: 4, 0x9f: 4, 0xa0: 4,
	0xa1: 4, 0xa2: 4, 0xa3: 4, 0xa4: 4, 0xa5: 4, 0xa6: 4, 0xa7: 4, 0xa8: 4,
	0xa9: 4, 0xaa: 4, 0xab: 4, 0xac: 4, 0xad: 4, 0xae: 4, 0xaf: 4, 0xb0: 4,
	0xb1: 4, 0xb2: 4, 0xb3: 4, 0xb4: 4, 0xb5: 4, 0xb6: 4, 0xb7: 4, 0xb8: 4,
	0xb9: 4, 0xba: 4, 0xbb: 4, 0xbc: 4, 0xbd: 4, 0xbe: 4, 0xbf: 4, 0xc0: 4,
	0xc1: 4, 0xc2: 4, 0xc3: 4, 0xc4: 4, 0xc5: 4, 0xc6: 4, 0xc7: 4, 0xc8: 4,
	0xc9: 4, 0xca: 4, 0xcb: 4, 0xcc: 4, 0xcd: 4, 0xce: 4, 0xcf: 4, 0xd0: 4,
	0xd1: 4, 0xd2: 4, 0xd3: 4, 0xd4: 4, 0xd5: 4, 0xd6: 4, 0xd7: 4, 0xd8: 4,
	0xd9: 4, 0xda: 4, 0xdb: 4, 0xdc: 4, 0xdd: 4, 0xde: 4, 0xdf: 4, 0xe0: 4,
	0xe1: 4, 0xe2: 4, 0xe3: 4, 0xe4: 4, 0xe5: 4, 0xe6: 4, 0xe7: 4, 0xe8: 4,
	0xe9: 4, 0xea: 4, 0xeb: 4, 0xec: 4, 0xed: 4, 0xee: 4, 0xef: 4, 0xf0: 4,
	0xf1: 4, 0xf2: 4, 0xf3: 4, 0xf4: 4, 0xf5: 4, 0xf6: 4, 0xf7: 4, 0xf8: 4,
	0xf9: 4, 0xfa: 4, 0xfb: 4, 0xfc: 4, 0xfd: 4, 0xfe: 4, 0xff: 4,

	0x00: 5, 0x01: 5, 0x02: 5, 0x03: 5, 0x04: 5, 0x05: 5, 0x06: 5, 0x07: 5,
	0x08: 5, 0x0b: 5, 0x0c: 5, 0x0e: 5, 0x0f: 5, 0x10: 5, 0x11: 5, 0x12: 5,
	0x13: 5, 0x14: 5, 0x15: 5, 0x16: 5, 0x17: 5, 0x18: 5, 0x19: 5, 0x1a: 5,
	0x1b: 5, 0x1c: 5, 0x1d: 5, 0x1e: 5, 0x1f: 5,
}
