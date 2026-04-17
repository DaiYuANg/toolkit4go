package dbx

import (
	"errors"
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
)

type relationSchemaSource interface {
	querydsl.TableSource
	schemaRef() schemaDefinition
}

func (q *SelectQuery) JoinRelation(source relationSchemaSource, relation relationAccessor, target querydsl.TableSource) (*SelectQuery, error) {
	return q.joinRelation(querydsl.InnerJoin, source, relation, target)
}

func (q *SelectQuery) LeftJoinRelation(source relationSchemaSource, relation relationAccessor, target querydsl.TableSource) (*SelectQuery, error) {
	return q.joinRelation(querydsl.LeftJoin, source, relation, target)
}

func (q *SelectQuery) RightJoinRelation(source relationSchemaSource, relation relationAccessor, target querydsl.TableSource) (*SelectQuery, error) {
	return q.joinRelation(querydsl.RightJoin, source, relation, target)
}

func (q *SelectQuery) joinRelation(joinType querydsl.JoinType, source relationSchemaSource, relation relationAccessor, target querydsl.TableSource) (*SelectQuery, error) {
	if q == nil {
		return nil, errors.New("dbx: select query is nil")
	}
	if source == nil {
		return nil, errors.New("dbx: relation join requires source schema")
	}
	if relation == nil {
		return nil, errors.New("dbx: relation join requires relation")
	}
	if target == nil {
		return nil, errors.New("dbx: relation join requires target table")
	}

	sourceTable := tableRef(source)
	if !q.containsTable(sourceTable) {
		return nil, fmt.Errorf("dbx: source table %s is not part of the query", sourceTable.Ref())
	}

	meta := relation.relationRef()
	targetTable := tableRef(target)
	if meta.TargetTable != "" && targetTable.Name() != meta.TargetTable {
		return nil, fmt.Errorf("dbx: relation %s targets table %s, got %s", meta.Name, meta.TargetTable, targetTable.Name())
	}

	joins, err := buildRelationJoins(joinType, source, meta, targetTable)
	if err != nil {
		return nil, err
	}
	q.Joins = mergeList(q.Joins, joins)
	return q, nil
}

func (q *SelectQuery) containsTable(table querydsl.Table) bool {
	if sameTable(q.FromItem, table) {
		return true
	}
	_, ok := collectionx.FindList(q.Joins, func(_ int, join Join) bool {
		return sameTable(join.Table, table)
	})
	return ok
}

func sameTable(left, right querydsl.Table) bool {
	return left.Name() == right.Name() && left.Alias() == right.Alias()
}

func buildDirectRelationPredicate(source relationSchemaSource, meta RelationMeta, target querydsl.Table) (Predicate, error) {
	localColumn, err := relationSourceColumn(source, meta)
	if err != nil {
		return nil, err
	}
	targetColumn, err := relationTargetColumn(target, meta)
	if err != nil {
		return nil, err
	}
	return metadataComparisonPredicate{
		left:  localColumn,
		op:    querydsl.OpEq,
		right: metadataColumnOperand{meta: targetColumn},
	}, nil
}

func buildRelationJoins(joinType querydsl.JoinType, source relationSchemaSource, meta RelationMeta, target querydsl.Table) (collectionx.List[Join], error) {
	joins := collectionx.NewListWithCapacity[Join](2)

	switch meta.Kind {
	case RelationBelongsTo, RelationHasOne, RelationHasMany:
		predicate, err := buildDirectRelationPredicate(source, meta, target)
		if err != nil {
			return nil, err
		}
		joins.Add(Join{Type: joinType, Table: target, Predicate: predicate})
		return joins, nil
	case RelationManyToMany:
		through, first, second, err := buildManyToManyJoins(source, meta, target)
		if err != nil {
			return nil, err
		}
		joins.Add(Join{Type: joinType, Table: through, Predicate: first})
		joins.Add(Join{Type: joinType, Table: target, Predicate: second})
		return joins, nil
	default:
		return nil, fmt.Errorf("dbx: unsupported relation kind %d", meta.Kind)
	}
}

func buildManyToManyJoins(source relationSchemaSource, meta RelationMeta, target querydsl.Table) (querydsl.Table, Predicate, Predicate, error) {
	if meta.ThroughTable == "" {
		return querydsl.Table{}, nil, nil, fmt.Errorf("dbx: many-to-many relation %s requires join table", meta.Name)
	}
	if meta.ThroughLocalColumn == "" || meta.ThroughTargetColumn == "" {
		return querydsl.Table{}, nil, nil, fmt.Errorf("dbx: many-to-many relation %s requires join_local and join_target", meta.Name)
	}

	sourceColumn, err := relationSourceColumn(source, meta)
	if err != nil {
		return querydsl.Table{}, nil, nil, err
	}
	targetColumn, err := relationTargetColumn(target, meta)
	if err != nil {
		return querydsl.Table{}, nil, nil, err
	}

	through := querydsl.NamedTable(meta.ThroughTable)
	throughSourceColumn := ColumnMeta{Name: meta.ThroughLocalColumn, Table: through.Name(), Alias: through.Alias()}
	throughTargetColumn := ColumnMeta{Name: meta.ThroughTargetColumn, Table: through.Name(), Alias: through.Alias()}

	first := metadataComparisonPredicate{
		left:  sourceColumn,
		op:    querydsl.OpEq,
		right: metadataColumnOperand{meta: throughSourceColumn},
	}
	second := metadataComparisonPredicate{
		left:  throughTargetColumn,
		op:    querydsl.OpEq,
		right: metadataColumnOperand{meta: targetColumn},
	}
	return through, first, second, nil
}

func relationSourceColumn(source relationSchemaSource, meta RelationMeta) (ColumnMeta, error) {
	name := meta.LocalColumn
	if name == "" {
		primaryKey := derivePrimaryKey(source.schemaRef())
		if primaryKey == nil || primaryKey.Columns.Len() != 1 {
			return ColumnMeta{}, fmt.Errorf("dbx: relation %s requires local column or single-column primary key", meta.Name)
		}
		name, _ = primaryKey.Columns.GetFirst()
	}

	column, ok := sourceColumnByName(source.schemaRef(), name)
	if !ok {
		return ColumnMeta{}, fmt.Errorf("dbx: relation %s source column %s not found", meta.Name, name)
	}
	return column, nil
}

func relationTargetColumn(target querydsl.Table, meta RelationMeta) (ColumnMeta, error) {
	if meta.TargetColumn == "" {
		return ColumnMeta{}, fmt.Errorf("dbx: relation %s requires target column", meta.Name)
	}
	return ColumnMeta{
		Name:  meta.TargetColumn,
		Table: target.Name(),
		Alias: target.Alias(),
	}, nil
}

func sourceColumnByName(def schemaDefinition, name string) (ColumnMeta, bool) {
	return def.columnByName(name)
}
