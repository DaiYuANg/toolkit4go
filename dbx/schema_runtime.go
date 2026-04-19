package dbx

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"
)

type schemaSelectItem struct {
	meta schemax.ColumnMeta
}

func (s schemaSelectItem) QueryExpression() {}
func (s schemaSelectItem) QuerySelectItem() {}

func (s schemaSelectItem) columnRef() schemax.ColumnMeta {
	return s.meta
}

func (s schemaSelectItem) ColumnRef() schemax.ColumnMeta {
	return s.columnRef()
}

func (s schemaSelectItem) RenderOperand(state *querydsl.State) (string, error) {
	return renderColumnOperand(state, s.columnRef())
}

func (s Schema[E]) AllColumns() collectionx.List[querydsl.SelectItem] {
	return collectionx.MapList(s.def.columns, func(_ int, column schemax.ColumnMeta) querydsl.SelectItem {
		return schemaSelectItem{meta: cloneColumnMeta(column)}
	})
}

func (s Schema[E]) PrimaryColumn() (schemax.ColumnMeta, bool) {
	column, ok := collectionx.FindList(s.def.columns, func(_ int, column schemax.ColumnMeta) bool {
		return column.PrimaryKey
	})
	if !ok {
		return schemax.ColumnMeta{}, false
	}
	return cloneColumnMeta(column), true
}

func (s Schema[E]) ColumnByName(name string) (schemax.ColumnMeta, bool) {
	column, ok := s.def.columnByName(name)
	if !ok {
		return schemax.ColumnMeta{}, false
	}
	return cloneColumnMeta(column), true
}

type metadataAssignment struct {
	meta  schemax.ColumnMeta
	value any
}

type metadataComparisonPredicate struct {
	left  schemax.ColumnMeta
	op    querydsl.ComparisonOperator
	right any
}

func (metadataAssignment) QueryAssignment() {}

func (metadataComparisonPredicate) QueryExpression() {}
func (metadataComparisonPredicate) QueryPredicate()  {}

func (a metadataAssignment) assignmentColumn() schemax.ColumnMeta {
	return a.meta
}

func (a metadataAssignment) AssignmentColumn() schemax.ColumnMeta {
	return a.assignmentColumn()
}

func (a metadataAssignment) RenderAssignment(state *querydsl.State) error {
	state.WriteQuotedIdent(a.meta.Name)
	state.WriteString(" = ")
	operand, err := querydsl.RenderOperandValue(state, a.value)
	if err != nil {
		return fmt.Errorf("dbx: render metadata assignment operand: %w", err)
	}
	state.WriteString(operand)
	return nil
}

func (a metadataAssignment) RenderAssignmentValue(state *querydsl.State) error {
	operand, err := querydsl.RenderOperandValue(state, a.value)
	if err != nil {
		return fmt.Errorf("dbx: render metadata assignment value: %w", err)
	}
	state.WriteString(operand)
	return nil
}

func (p metadataComparisonPredicate) RenderPredicate(state *querydsl.State) error {
	state.RenderColumn(p.left)
	if p.op == querydsl.OpIs || p.op == querydsl.OpIsNot {
		state.WriteRawByte(' ')
		state.WriteString(string(p.op))
		state.WriteString(" NULL")
		return nil
	}

	operand, err := querydsl.RenderOperandValue(state, p.right)
	if err != nil {
		return fmt.Errorf("dbx: render metadata predicate operand: %w", err)
	}
	state.WriteRawByte(' ')
	state.WriteString(string(p.op))
	state.WriteRawByte(' ')
	state.WriteString(operand)
	return nil
}
