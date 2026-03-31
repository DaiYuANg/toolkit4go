package dbx

import (
	"context"
	"database/sql"
)

func AutoMigrate(ctx context.Context, session Session, schemas ...SchemaResource) (ValidationReport, error) {
	logRuntimeNode(session, "schema.auto_migrate.start", "schemas", len(schemas))
	plan, err := PlanSchemaChanges(ctx, session, schemas...)
	if err != nil {
		logRuntimeNode(session, "schema.auto_migrate.error", "stage", "plan", "error", err)
		return ValidationReport{}, err
	}
	if plan.HasManualActions() {
		logRuntimeNode(session, "schema.auto_migrate.manual_required", "actions", len(plan.Actions))
		return plan.Report, SchemaDriftError{Report: plan.Report}
	}

	report, err := applyMigrationPlan(ctx, session, plan, schemas...)
	if err != nil {
		return ValidationReport{}, err
	}
	logRuntimeNode(session, "schema.auto_migrate.done", "actions", len(plan.Actions))
	return report, nil
}

func applyMigrationPlan(ctx context.Context, session Session, plan MigrationPlan, schemas ...SchemaResource) (ValidationReport, error) {
	executableActions := plan.ExecutableActions()
	execSession, finalize, rollback, transactional, err := autoMigrateExecutionSession(ctx, session, len(executableActions) > 0)
	if err != nil {
		logRuntimeNode(session, "schema.auto_migrate.error", "stage", "begin_tx", "error", err)
		return ValidationReport{}, err
	}

	committed := false
	if rollback != nil {
		defer rollbackPendingMigration(session, rollback, &committed)
	}
	execErr := executeMigrationActions(ctx, session, execSession, plan.Actions)
	if execErr != nil {
		return ValidationReport{}, execErr
	}

	report, err := validateAppliedMigration(ctx, session, execSession, schemas...)
	if err != nil {
		return ValidationReport{}, err
	}
	report = appendNonTransactionalWarning(report, transactional, len(executableActions))
	if err := ensureMigrationReportValid(session, report); err != nil {
		return report, err
	}
	if err := finalizeMigration(finalize, &committed, session); err != nil {
		return ValidationReport{}, err
	}
	return report, nil
}

func rollbackPendingMigration(session Session, rollback func() error, committed *bool) {
	if *committed {
		return
	}
	if rollbackErr := rollback(); rollbackErr != nil {
		logRuntimeNode(session, "schema.auto_migrate.error", "stage", "rollback", "error", rollbackErr)
	}
}

func executeMigrationActions(ctx context.Context, session, execSession Session, actions []MigrationAction) error {
	for _, action := range actions {
		if !action.Executable {
			continue
		}
		logRuntimeNode(execSession, "schema.auto_migrate.exec_action", "kind", action.Kind, "table", action.Table, "summary", action.Summary)
		if _, err := execSession.ExecBoundContext(ctx, action.Statement); err != nil {
			logRuntimeNode(session, "schema.auto_migrate.error", "stage", "exec", "kind", action.Kind, "table", action.Table, "error", err)
			return wrapDBError("apply schema migration action", err)
		}
	}
	return nil
}

func validateAppliedMigration(ctx context.Context, session, execSession Session, schemas ...SchemaResource) (ValidationReport, error) {
	report, err := ValidateSchemas(ctx, execSession, schemas...)
	if err != nil {
		logRuntimeNode(session, "schema.auto_migrate.error", "stage", "validate", "error", err)
		return ValidationReport{}, err
	}
	return report, nil
}

func appendNonTransactionalWarning(report ValidationReport, transactional bool, actionCount int) ValidationReport {
	if transactional || actionCount <= 1 {
		return report
	}
	return report.withWarning("dbx: auto migrate executed without transaction; partial application is possible on failure")
}

func ensureMigrationReportValid(session Session, report ValidationReport) error {
	if report.Valid() {
		return nil
	}
	logRuntimeNode(session, "schema.auto_migrate.invalid_after_apply", "tables", len(report.Tables))
	return SchemaDriftError{Report: report}
}

func finalizeMigration(finalize func() error, committed *bool, session Session) error {
	if finalize == nil {
		return nil
	}
	if err := finalize(); err != nil {
		logRuntimeNode(session, "schema.auto_migrate.error", "stage", "commit", "error", err)
		return err
	}
	*committed = true
	return nil
}

type txStarter interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
}

func autoMigrateExecutionSession(ctx context.Context, session Session, needExec bool) (Session, func() error, func() error, bool, error) {
	if !needExec {
		logRuntimeNode(session, "schema.auto_migrate.execution_session", "need_exec", false, "transactional", false)
		return session, nil, nil, false, nil
	}

	starter, ok := session.(txStarter)
	if !ok {
		logRuntimeNode(session, "schema.auto_migrate.execution_session", "need_exec", true, "transactional", false, "reason", "session_has_no_begin_tx")
		return session, nil, nil, false, nil
	}

	tx, err := starter.BeginTx(ctx, nil)
	if err != nil {
		logRuntimeNode(session, "schema.auto_migrate.execution_session.error", "error", err)
		return nil, nil, nil, false, wrapDBError("begin schema migration transaction", err)
	}
	logRuntimeNode(session, "schema.auto_migrate.execution_session", "need_exec", true, "transactional", true)
	return tx, tx.Commit, tx.Rollback, true, nil
}
