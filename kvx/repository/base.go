package repository

import (
	"context"

	"github.com/DaiYuANg/arcgo/kvx/mapping"
)

type repositoryBase[T any] struct {
	keyBuilder *mapping.KeyBuilder
	tagParser  *mapping.TagParser
	indexer    *Indexer[T]
}

func (b repositoryBase[T]) metadata(entity *T) (*mapping.EntityMetadata, error) {
	return b.tagParser.ParseType(entity)
}

func (b repositoryBase[T]) metadataForType() (*mapping.EntityMetadata, error) {
	var zero T
	return b.tagParser.ParseType(&zero)
}

func (b repositoryBase[T]) keyFromID(id string) string {
	return b.keyBuilder.BuildWithID(id)
}

func (b repositoryBase[T]) idsByField(ctx context.Context, fieldName, fieldValue string) ([]string, error) {
	metadata, err := b.metadataForType()
	if err != nil {
		return nil, err
	}
	_, fieldTag, ok := metadata.ResolveField(fieldName)
	if !ok {
		return nil, ErrFieldNotFound
	}
	return b.indexer.GetEntityIDsByField(ctx, fieldTag.IndexNameOrDefault(), fieldValue)
}

func (b repositoryBase[T]) hydrateEntityID(entity *T, metadata *mapping.EntityMetadata, key string) error {
	return metadata.SetEntityID(entity, extractIDFromKey(key))
}
