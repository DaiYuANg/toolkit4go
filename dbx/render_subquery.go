package dbx

import (
	"errors"
	"strings"
)

type subqueryOperand struct {
	Query *SelectQuery
}

func (subqueryOperand) expressionNode() {}

func (s subqueryOperand) renderOperand(state *renderState) (string, error) {
	if s.Query == nil {
		return "", errors.New("dbx: subquery is nil")
	}
	original := state.buf
	var builder strings.Builder
	state.buf = builder
	if err := renderSelectStatement(state, s.Query); err != nil {
		state.buf = original
		return "", err
	}
	rendered := state.buf.String()
	state.buf = original
	return "(" + rendered + ")", nil
}
