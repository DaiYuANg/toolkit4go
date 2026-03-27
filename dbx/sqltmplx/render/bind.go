package render

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var (
	errSpreadParamEmpty       = errors.New("sqltmplx: spread parameter is empty")
	errSpreadParamType        = errors.New("sqltmplx: spread parameter must be slice or array")
	errUnterminatedSQLComment = errors.New("sqltmplx: unterminated sql comment")
	errInvalidCastSuffix      = errors.New("invalid cast suffix")
	errUnterminatedBalanced   = errors.New("unterminated balanced literal")
	errUnterminatedQuoted     = errors.New("unterminated quoted literal")
	errEmptyScalarLiteral     = errors.New("empty scalar literal")
)

func bindText(input string, st *state) (string, error) {
	var out strings.Builder
	for i := 0; i < len(input); {
		if i+1 < len(input) && input[i] == '/' && input[i+1] == '*' {
			text, next, handled, err := bindCommentPlaceholder(input, i, st)
			if err != nil {
				return "", err
			}
			if handled {
				writeBuilderString(&out, text)
				i = next
				continue
			}
		}
		writeBuilderByte(&out, input[i])
		i++
	}
	return out.String(), nil
}

func bindCommentPlaceholder(input string, i int, st *state) (string, int, bool, error) {
	raw, commentEnd, handled, err := parseCommentPlaceholder(input, i)
	if err != nil || !handled {
		return "", commentEnd, handled, err
	}

	sampleStart, err := placeholderSampleStart(input, commentEnd, raw)
	if err != nil {
		return "", 0, false, err
	}

	spread := input[sampleStart] == '(' || looksLikeCollectionSample(input, sampleStart)
	text, err := bindParam(raw, spread, st)
	if err != nil {
		return "", 0, false, err
	}

	k, err := skipPlaceholderSample(input, sampleStart)
	if err != nil {
		return "", 0, false, fmt.Errorf("sqltmplx: placeholder %q invalid test literal: %w", raw, err)
	}
	return text, k, true, nil
}

func bindParam(name string, spread bool, st *state) (string, error) {
	val, err := lookupParam(st.params, name)
	if err != nil {
		return "", err
	}
	if !spread {
		st.args = append(st.args, val)
		return st.nextBind(), nil
	}

	return bindSpreadParam(val, st)
}

func parseCommentPlaceholder(input string, start int) (string, int, bool, error) {
	endComment := strings.Index(input[start+2:], "*/")
	if endComment < 0 {
		return "", 0, false, errUnterminatedSQLComment
	}

	commentEnd := start + 2 + endComment + 2
	raw := strings.TrimSpace(input[start+2 : start+2+endComment])
	if raw == "" || strings.HasPrefix(raw, "%") || !isParamPath(raw) {
		return "", commentEnd, false, nil
	}
	return raw, commentEnd, true, nil
}

func placeholderSampleStart(input string, commentEnd int, raw string) (int, error) {
	start := skipSpaces(input, commentEnd)
	if start >= len(input) {
		return 0, fmt.Errorf("sqltmplx: placeholder %q missing test literal", raw)
	}
	return start, nil
}

func lookupParam(params any, name string) (any, error) {
	valOpt := lookup(params, name)
	if valOpt.IsAbsent() {
		return nil, fmt.Errorf("sqltmplx: parameter %q not found", name)
	}
	return valOpt.MustGet(), nil
}

func bindSpreadParam(val any, st *state) (string, error) {
	rv, err := spreadValue(val)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	length := rv.Len()
	for j := range length {
		appendSpreadBind(&out, j, st.nextBind())
		st.args = append(st.args, rv.Index(j).Interface())
	}
	return out.String(), nil
}

func spreadValue(val any) (reflect.Value, error) {
	rv := reflect.ValueOf(val)
	for rv.IsValid() && rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return reflect.Value{}, errSpreadParamEmpty
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return reflect.Value{}, errSpreadParamEmpty
	}
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return reflect.Value{}, errSpreadParamType
	}
	if rv.Len() == 0 {
		return reflect.Value{}, errSpreadParamEmpty
	}
	return rv, nil
}

func appendSpreadBind(out *strings.Builder, index int, bind string) {
	if index > 0 {
		writeBuilderString(out, ", ")
	}
	writeBuilderString(out, bind)
}

func isParamPath(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case r == '.':
			if i == 0 || i == len(s)-1 {
				return false
			}
		case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
		default:
			return false
		}
	}
	return true
}

func looksLikeCollectionSample(input string, start int) bool {
	i := start
	if i < len(input) && isIdentifierStart(rune(input[i])) {
		j := i + 1
		for j < len(input) && isIdentifierPart(rune(input[j])) {
			j++
		}
		return j < len(input) && input[j] == '['
	}
	return false
}

func skipPlaceholderSample(input string, start int) (int, error) {
	i := start
	switch {
	case input[i] == '(':
		var err error
		i, err = skipBalanced(input, i, '(', ')')
		if err != nil {
			return 0, err
		}
	case input[i] == '\'' || input[i] == '"':
		var err error
		i, err = skipQuoted(input, i)
		if err != nil {
			return 0, err
		}
	case isIdentifierStart(rune(input[i])):
		var err error
		i, err = skipIdentifierExpr(input, i)
		if err != nil {
			return 0, err
		}
	default:
		var err error
		i, err = skipScalarToken(input, i)
		if err != nil {
			return 0, err
		}
	}
	return skipExprSuffixes(input, i)
}

func skipIdentifierExpr(input string, start int) (int, error) {
	i := start + 1
	for i < len(input) && isIdentifierPart(rune(input[i])) {
		i++
	}
	for {
		next, handled, err := skipIdentifierCall(input, i)
		if err != nil {
			return 0, err
		}
		if !handled {
			return i, nil
		}
		i = next
	}
}

func skipExprSuffixes(input string, start int) (int, error) {
	i := start
	for {
		next, handled, err := nextExprSuffix(input, i)
		if err != nil {
			return 0, err
		}
		if !handled {
			return i, nil
		}
		i = next
	}
}

func isIdentifierStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentifierPart(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func skipBalanced(input string, start int, open, closeByte byte) (int, error) {
	depth := 0
	for i := start; i < len(input); i++ {
		switch input[i] {
		case '\'', '"':
			j, err := skipQuoted(input, i)
			if err != nil {
				return 0, err
			}
			i = j - 1
		case open:
			depth++
		case closeByte:
			depth--
			if depth == 0 {
				return i + 1, nil
			}
		}
	}
	return 0, errUnterminatedBalanced
}

func skipQuoted(input string, start int) (int, error) {
	quote := input[start]
	for i := start + 1; i < len(input); i++ {
		if input[i] != quote {
			continue
		}
		if i+1 < len(input) && input[i+1] == quote {
			i++
			continue
		}
		return i + 1, nil
	}
	return 0, errUnterminatedQuoted
}

func skipScalarToken(input string, start int) (int, error) {
	i := start
	for i < len(input) {
		r := rune(input[i])
		if unicode.IsSpace(r) || r == ',' || r == ')' || r == '(' || r == ']' {
			break
		}
		i++
	}
	if i == start {
		return 0, errEmptyScalarLiteral
	}
	return i, nil
}

func skipIdentifierCall(input string, start int) (int, bool, error) {
	next := skipSpaces(input, start)
	if next >= len(input) || input[next] != '(' {
		return 0, false, nil
	}
	index, err := skipBalanced(input, next, '(', ')')
	if err != nil {
		return 0, false, err
	}
	return index, true, nil
}

func skipCastSuffix(input string, start int) (int, bool, error) {
	if !hasCastPrefix(input, start) {
		return 0, false, nil
	}

	index := start + 2
	if index >= len(input) || !isIdentifierStart(rune(input[index])) {
		return 0, false, errInvalidCastSuffix
	}
	return scanCastTypeSuffix(input, index+1), true, nil
}

func skipIndexSuffix(input string, start int) (int, bool, error) {
	if start >= len(input) || input[start] != '[' {
		return 0, false, nil
	}
	index, err := skipBalanced(input, start, '[', ']')
	if err != nil {
		return 0, false, err
	}
	return index, true, nil
}

func skipSpaces(input string, start int) int {
	index := start
	for index < len(input) && unicode.IsSpace(rune(input[index])) {
		index++
	}
	return index
}

func writeBuilderByte(builder *strings.Builder, value byte) {
	if err := builder.WriteByte(value); err != nil {
		panic(err)
	}
}

func nextExprSuffix(input string, start int) (int, bool, error) {
	next := skipSpaces(input, start)
	if castNext, handled, err := skipCastSuffix(input, next); handled || err != nil {
		return castNext, handled, err
	}
	return skipIndexSuffix(input, next)
}

func hasCastPrefix(input string, start int) bool {
	return start+1 < len(input) && input[start] == ':' && input[start+1] == ':'
}

func scanCastTypeSuffix(input string, start int) int {
	index := start
	for index < len(input) {
		if isCastIdentifierChar(rune(input[index])) {
			index++
			continue
		}
		if hasArraySuffix(input, index) {
			index += 2
			continue
		}
		break
	}
	return index
}

func isCastIdentifierChar(r rune) bool {
	return isIdentifierPart(r) || r == '.'
}

func hasArraySuffix(input string, start int) bool {
	return start+1 < len(input) && input[start] == '[' && input[start+1] == ']'
}
