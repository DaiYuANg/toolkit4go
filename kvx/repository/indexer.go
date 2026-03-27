package repository

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/DaiYuANg/arcgo/collectionx/set"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/DaiYuANg/arcgo/kvx/mapping"
	"github.com/samber/lo"
)

// Indexer manages secondary indexes for entities.
type Indexer[T any] struct {
	kv         kvx.KV
	keyBuilder *mapping.KeyBuilder
}

// NewIndexer creates a new Indexer.
func NewIndexer[T any](kv kvx.KV, keyPrefix string) *Indexer[T] {
	return &Indexer[T]{
		kv:         kv,
		keyBuilder: mapping.NewKeyBuilder(keyPrefix),
	}
}

// IndexEntity adds an entity to secondary indexes.
func (i *Indexer[T]) IndexEntity(ctx context.Context, entity *T, metadata *mapping.EntityMetadata, entityKey string) error {
	v := reflect.Indirect(reflect.ValueOf(entity))
	for _, fieldName := range metadata.IndexFields {
		fieldTag, ok := metadata.Fields[fieldName]
		if !ok {
			continue
		}
		fieldVal := v.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			continue
		}
		fieldValue := formatIndexValue(fieldVal)
		if fieldValue == "" {
			continue
		}
		indexKey := i.buildIndexKey(fieldTag.IndexNameOrDefault(), fieldValue)
		if err := i.addToIndex(ctx, indexKey, extractIDFromKey(entityKey)); err != nil {
			return fmt.Errorf("failed to index field %s: %w", fieldName, err)
		}
	}
	return nil
}

// RemoveEntityFromIndexes removes an entity from all secondary indexes.
func (i *Indexer[T]) RemoveEntityFromIndexes(ctx context.Context, entity *T, metadata *mapping.EntityMetadata) error {
	v := reflect.Indirect(reflect.ValueOf(entity))
	entityKey, err := i.keyBuilder.Build(entity, metadata)
	if err != nil {
		return wrapRepositoryError(err, "build entity key for index removal")
	}
	entityID := extractIDFromKey(entityKey)

	for _, fieldName := range metadata.IndexFields {
		fieldTag, ok := metadata.Fields[fieldName]
		if !ok {
			continue
		}
		fieldVal := v.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			continue
		}
		fieldValue := formatIndexValue(fieldVal)
		if fieldValue == "" {
			continue
		}
		indexKey := i.buildIndexKey(fieldTag.IndexNameOrDefault(), fieldValue)
		if err := i.removeFromIndex(ctx, indexKey, entityID); err != nil {
			return fmt.Errorf("failed to remove index for field %s: %w", fieldName, err)
		}
	}
	return nil
}

// EntityIndexEntries returns all concrete index entry keys for an entity.
func (i *Indexer[T]) EntityIndexEntries(entity *T, metadata *mapping.EntityMetadata, entityKey string) ([]string, error) {
	if entity == nil {
		return nil, nil
	}
	v := reflect.Indirect(reflect.ValueOf(entity))
	entityID := extractIDFromKey(entityKey)
	results := make([]string, 0, len(metadata.IndexFields))
	for _, fieldName := range metadata.IndexFields {
		fieldTag, ok := metadata.Fields[fieldName]
		if !ok {
			continue
		}
		fieldVal := v.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			continue
		}
		if entry := i.indexEntryKey(fieldTag.IndexNameOrDefault(), formatIndexValue(fieldVal), entityID); entry != "" {
			results = append(results, entry)
		}
	}
	return results, nil
}

// ReplaceEntityIndexEntries calculates the index diff for replacing one entity with another.
func (i *Indexer[T]) ReplaceEntityIndexEntries(_ context.Context, oldEntity, newEntity *T, metadata *mapping.EntityMetadata, entityKey string) ([]string, []string, error) {
	oldEntries, err := i.EntityIndexEntries(oldEntity, metadata, entityKey)
	if err != nil {
		return nil, nil, err
	}
	newEntries, err := i.EntityIndexEntries(newEntity, metadata, entityKey)
	if err != nil {
		return nil, nil, err
	}
	return diffIndexEntries(oldEntries, newEntries), diffIndexEntries(newEntries, oldEntries), nil
}

// ReplaceFieldIndexEntries calculates the index diff for updating one indexed field.
func (i *Indexer[T]) ReplaceFieldIndexEntries(metadata *mapping.EntityMetadata, fieldName, entityKey string, entity *T, newValue any) ([]string, []string, error) {
	resolvedField, fieldTag, exists := metadata.ResolveField(fieldName)
	if !exists || !fieldTag.Index || entity == nil {
		return nil, nil, nil
	}
	v := reflect.Indirect(reflect.ValueOf(entity))
	fieldVal := v.FieldByName(resolvedField)
	if !fieldVal.IsValid() {
		return nil, nil, nil
	}

	entityID := extractIDFromKey(entityKey)
	oldEntry := i.indexEntryKey(fieldTag.IndexNameOrDefault(), formatIndexValue(fieldVal), entityID)
	newEntry := i.indexEntryKey(fieldTag.IndexNameOrDefault(), formatIndexValue(reflect.ValueOf(newValue)), entityID)
	if oldEntry == newEntry {
		return nil, nil, nil
	}

	oldEntries := make([]string, 0, 1)
	newEntries := make([]string, 0, 1)
	if oldEntry != "" {
		oldEntries = append(oldEntries, oldEntry)
	}
	if newEntry != "" {
		newEntries = append(newEntries, newEntry)
	}
	return oldEntries, newEntries, nil
}

// ApplyIndexDiff removes stale index entries and writes the new ones.
func (i *Indexer[T]) ApplyIndexDiff(ctx context.Context, removeEntries, addEntries []string) error {
	for _, entry := range removeEntries {
		if err := i.kv.Delete(ctx, entry); err != nil {
			return wrapRepositoryError(err, "remove stale index entry")
		}
	}
	for _, entry := range addEntries {
		if err := i.kv.Set(ctx, entry, []byte("1"), 0); err != nil {
			return wrapRepositoryError(err, "write index entry")
		}
	}
	return nil
}

// GetEntityIDsByField returns entity IDs that have the specified field value.
func (i *Indexer[T]) GetEntityIDsByField(ctx context.Context, fieldName, fieldValue string) ([]string, error) {
	return i.getIndexMembers(ctx, i.buildIndexKey(fieldName, fieldValue))
}

func (i *Indexer[T]) buildIndexKey(fieldName, fieldValue string) string {
	prefix := strings.TrimSuffix(i.keyBuilder.BuildWithID(""), ":")
	if prefix == "" {
		return "idx:" + fieldName + ":" + fieldValue
	}
	return prefix + ":idx:" + fieldName + ":" + fieldValue
}

func (i *Indexer[T]) indexEntryKey(fieldName, fieldValue, entityID string) string {
	if fieldValue == "" || entityID == "" {
		return ""
	}
	return i.buildIndexKey(fieldName, fieldValue) + ":" + entityID
}

func (i *Indexer[T]) addToIndex(ctx context.Context, indexKey, entityID string) error {
	return wrapRepositoryError(i.kv.Set(ctx, indexKey+":"+entityID, []byte("1"), 0), "write index entry")
}

func (i *Indexer[T]) removeFromIndex(ctx context.Context, indexKey, entityID string) error {
	return wrapRepositoryError(i.kv.Delete(ctx, indexKey+":"+entityID), "delete index entry")
}

func (i *Indexer[T]) getIndexMembers(ctx context.Context, indexKey string) ([]string, error) {
	keys, err := i.kv.Keys(ctx, indexKey+":*")
	if err != nil {
		return nil, wrapRepositoryError(err, "list index members")
	}
	results := make([]string, 0, len(keys))
	prefixLen := len(indexKey) + 1
	for _, key := range keys {
		if len(key) > prefixLen {
			results = append(results, key[prefixLen:])
		}
	}
	return results, nil
}

func formatIndexValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}
	if v.Kind() == reflect.String {
		return v.String()
	}
	if isSignedIndexKind(v.Kind()) {
		return strconv.FormatInt(v.Int(), 10)
	}
	if isUnsignedIndexKind(v.Kind()) {
		return strconv.FormatUint(v.Uint(), 10)
	}
	if v.Kind() == reflect.Bool {
		return strconv.FormatBool(v.Bool())
	}
	return fmt.Sprint(v.Interface())
}

func diffIndexEntries(left, right []string) []string {
	if len(left) == 0 {
		return nil
	}
	rightSet := set.NewSet[string](right...)
	return lo.Filter(left, func(entry string, _ int) bool {
		return !rightSet.Contains(entry)
	})
}

func extractIDFromKey(key string) string {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			return key[i+1:]
		}
	}
	return key
}

func isSignedIndexKind(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Int64
}

func isUnsignedIndexKind(kind reflect.Kind) bool {
	return kind >= reflect.Uint && kind <= reflect.Uintptr
}
