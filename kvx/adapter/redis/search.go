package redis

import (
	"context"
	"github.com/DaiYuANg/archgo/kvx"
)

// ============== Search Interface ==============

// CreateIndex creates a secondary index.
func (a *Adapter) CreateIndex(ctx context.Context, indexName string, prefix string, schema []kvx.SchemaField) error {
	args := make([]interface{}, 0)
	args = append(args, indexName, "ON", "HASH", "PREFIX", 1, prefix, "SCHEMA")

	for _, f := range schema {
		args = append(args, f.Name, string(f.Type))
		if f.Sortable {
			args = append(args, "SORTABLE")
		}
	}

	return a.client.Do(ctx, args...).Err()
}

// DropIndex drops a secondary index.
func (a *Adapter) DropIndex(ctx context.Context, indexName string) error {
	return a.client.Do(ctx, "FT.DROPINDEX", indexName).Err()
}

// Search performs a search query.
func (a *Adapter) Search(ctx context.Context, indexName string, query string, limit int) ([]string, error) {
	val, err := a.client.Do(ctx, "FT.SEARCH", indexName, query, "LIMIT", 0, limit).Result()
	if err != nil {
		return nil, err
	}

	// Parse FT.SEARCH response
	// Response format: [total, key1, [field1, value1, ...], key2, ...]
	return parseFTSearchResponse(val)
}

// SearchWithSort performs a search query with sorting.
func (a *Adapter) SearchWithSort(ctx context.Context, indexName string, query string, sortBy string, ascending bool, limit int) ([]string, error) {
	args := []interface{}{"FT.SEARCH", indexName, query, "SORTBY", sortBy}
	if !ascending {
		args = append(args, "DESC")
	}
	args = append(args, "LIMIT", 0, limit)

	val, err := a.client.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}

	return parseFTSearchResponse(val)
}

// SearchAggregate performs an aggregation query.
func (a *Adapter) SearchAggregate(ctx context.Context, indexName string, query string, limit int) ([]map[string]interface{}, error) {
	val, err := a.client.Do(ctx, "FT.AGGREGATE", indexName, query, "LIMIT", 0, limit).Result()
	if err != nil {
		return nil, err
	}

	// Parse FT.AGGREGATE response
	return parseFTAggregateResponse(val)
}
