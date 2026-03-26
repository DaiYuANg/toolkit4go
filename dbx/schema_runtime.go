package dbx

import (
	"strings"

	"github.com/samber/lo"
)

type schemaSelectItem struct {
	meta ColumnMeta
}

func (s schemaSelectItem) selectItemNode() {}

func (s schemaSelectItem) columnRef() ColumnMeta {
	return s.meta
}

func (s Schema[E]) AllColumns() []SelectItem {
	return lo.Map(s.def.columns, func(column ColumnMeta, _ int) SelectItem {
		return schemaSelectItem{meta: column}
	})
}

func (s Schema[E]) PrimaryColumn() (ColumnMeta, bool) {
	return lo.Find(s.def.columns, func(column ColumnMeta) bool {
		return column.PrimaryKey
	})
}

func (s Schema[E]) ColumnByName(name string) (ColumnMeta, bool) {
	return lo.Find(s.def.columns, func(column ColumnMeta) bool {
		return column.Name == name
	})
}

type metadataAssignment struct {
	meta  ColumnMeta
	value any
}

type metadataComparisonPredicate struct {
	left  ColumnMeta
	op    ComparisonOperator
	right any
}

type metadataColumnOperand struct {
	meta ColumnMeta
}

func (metadataAssignment) assignmentNode() {}

func (metadataComparisonPredicate) expressionNode() {}
func (metadataComparisonPredicate) predicateNode()  {}
func (metadataColumnOperand) expressionNode()       {}

func (a metadataAssignment) assignmentColumn() ColumnMeta {
	return a.meta
}

func (a metadataAssignment) renderAssignment(state *renderState) error {
	state.writeQuotedIdent(a.meta.Name)
	state.buf.WriteString(" = ")
	operand, err := renderOperandValue(state, a.value)
	if err != nil {
		return err
	}
	state.buf.WriteString(operand)
	return nil
}

func (a metadataAssignment) renderAssignmentValue(state *renderState) error {
	operand, err := renderOperandValue(state, a.value)
	if err != nil {
		return err
	}
	state.buf.WriteString(operand)
	return nil
}

func (p metadataComparisonPredicate) renderPredicate(state *renderState) error {
	state.renderColumn(p.left)
	if p.op == OpIs || p.op == OpIsNot {
		state.buf.WriteByte(' ')
		state.buf.WriteString(string(p.op))
		state.buf.WriteString(" NULL")
		return nil
	}

	operand, err := renderOperandValue(state, p.right)
	if err != nil {
		return err
	}
	state.buf.WriteByte(' ')
	state.buf.WriteString(string(p.op))
	state.buf.WriteByte(' ')
	state.buf.WriteString(operand)
	return nil
}

func (o metadataColumnOperand) renderOperand(state *renderState) (string, error) {
	var builder strings.Builder
	table := o.meta.Table
	if o.meta.Alias != "" {
		table = o.meta.Alias
	}
	builder.WriteString(state.dialect.QuoteIdent(table))
	builder.WriteByte('.')
	builder.WriteString(state.dialect.QuoteIdent(o.meta.Name))
	return builder.String(), nil
}
