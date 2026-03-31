package dbx

import (
	"errors"
	"fmt"
)

func renderUpsert(state *renderState, q *InsertQuery) error {
	if q.Upsert == nil {
		return nil
	}
	switch dialectFeatures(state.dialect).UpsertVariant {
	case "on_conflict":
		return renderUpsertOnConflict(state, q)
	case "on_duplicate_key":
		return renderUpsertOnDuplicateKey(state, q)
	default:
		return fmt.Errorf("dbx: upsert is not supported for dialect %s", state.dialect.Name())
	}
}

func renderUpsertOnConflict(state *renderState, q *InsertQuery) error {
	state.writeString(" ON CONFLICT")
	if err := renderUpsertTargets(state, q.Upsert.Targets); err != nil {
		return err
	}
	if q.Upsert.DoNothing {
		state.writeString(" DO NOTHING")
		return nil
	}
	if err := validateUpsertAssignments(q.Upsert); err != nil {
		return err
	}
	state.writeString(" DO UPDATE SET ")
	return renderUpsertAssignments(state, q.Upsert.Assignments)
}

func renderUpsertTargets(state *renderState, targets []Expression) error {
	if len(targets) == 0 {
		return nil
	}
	state.writeString(" (")
	for i, target := range targets {
		if i > 0 {
			state.writeString(", ")
		}
		if column, ok := target.(columnAccessor); ok {
			state.writeQuotedIdent(column.columnRef().Name)
			continue
		}
		operand, err := renderOperandValue(state, target)
		if err != nil {
			return err
		}
		state.writeString(operand)
	}
	state.writeByte(')')
	return nil
}

func validateUpsertAssignments(upsert *UpsertClause) error {
	switch {
	case len(upsert.Assignments) == 0:
		return errors.New("dbx: upsert update requires assignments")
	case len(upsert.Targets) == 0:
		return errors.New("dbx: upsert update requires conflict targets")
	default:
		return nil
	}
}

func renderUpsertOnDuplicateKey(state *renderState, q *InsertQuery) error {
	if q.Upsert.DoNothing {
		return nil
	}
	if len(q.Upsert.Assignments) == 0 {
		return errors.New("dbx: upsert update requires assignments")
	}
	state.writeString(" ON DUPLICATE KEY UPDATE ")
	return renderUpsertAssignments(state, q.Upsert.Assignments)
}

func renderUpsertAssignments(state *renderState, assignments []Assignment) error {
	for i, assignment := range assignments {
		if i > 0 {
			state.writeString(", ")
		}
		if err := renderAssignment(state, assignment); err != nil {
			return err
		}
	}
	return nil
}

func renderReturning(state *renderState, items []SelectItem) error {
	if len(items) == 0 {
		return nil
	}
	if !dialectFeatures(state.dialect).SupportsReturning {
		return fmt.Errorf("dbx: RETURNING is not supported for dialect %s", state.dialect.Name())
	}
	state.writeString(" RETURNING ")
	for i, item := range items {
		if i > 0 {
			state.writeString(", ")
		}
		if err := renderSelectItem(state, item); err != nil {
			return err
		}
	}
	return nil
}
