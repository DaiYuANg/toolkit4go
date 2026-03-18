package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/DaiYuANg/arcgo/kvx/mapping"
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
		return err
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

// UpdateFieldIndex updates the index when a field value changes.
func (i *Indexer[T]) UpdateFieldIndex(ctx context.Context, entity *T, metadata *mapping.EntityMetadata, fieldName string, entityKey string) error {
	resolvedField, fieldTag, exists := metadata.ResolveField(fieldName)
	if !exists || !fieldTag.Index {
		return nil
	}
	v := reflect.Indirect(reflect.ValueOf(entity))
	fieldVal := v.FieldByName(resolvedField)
	if !fieldVal.IsValid() {
		return nil
	}
	fieldValue := formatIndexValue(fieldVal)
	if fieldValue == "" {
		return nil
	}
	indexKey := i.buildIndexKey(fieldTag.IndexNameOrDefault(), fieldValue)
	return i.addToIndex(ctx, indexKey, extractIDFromKey(entityKey))
}

// GetEntityIDsByField returns entity IDs that have the specified field value.
func (i *Indexer[T]) GetEntityIDsByField(ctx context.Context, fieldName string, fieldValue string) ([]string, error) {
	return i.getIndexMembers(ctx, i.buildIndexKey(fieldName, fieldValue))
}

func (i *Indexer[T]) buildIndexKey(fieldName, fieldValue string) string {
	prefix := strings.TrimSuffix(i.keyBuilder.BuildWithID(""), ":")
	if prefix == "" {
		return fmt.Sprintf("idx:%s:%s", fieldName, fieldValue)
	}
	return fmt.Sprintf("%s:idx:%s:%s", prefix, fieldName, fieldValue)
}

func (i *Indexer[T]) addToIndex(ctx context.Context, indexKey string, entityID string) error {
	return i.kv.Set(ctx, fmt.Sprintf("%s:%s", indexKey, entityID), []byte("1"), 0)
}

func (i *Indexer[T]) removeFromIndex(ctx context.Context, indexKey string, entityID string) error {
	return i.kv.Delete(ctx, fmt.Sprintf("%s:%s", indexKey, entityID))
}

func (i *Indexer[T]) getIndexMembers(ctx context.Context, indexKey string) ([]string, error) {
	keys, err := i.kv.Keys(ctx, indexKey+":*")
	if err != nil {
		return nil, err
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
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func extractIDFromKey(key string) string {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			return key[i+1:]
		}
	}
	return key
}
