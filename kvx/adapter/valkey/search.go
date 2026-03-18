package valkey

import (
	"context"
	"fmt"
	"github.com/DaiYuANg/arcgo/kvx"
)

// ============== Search Interface ==============

// CreateIndex creates a secondary index.
func (a *Adapter) CreateIndex(ctx context.Context, indexName string, prefix string, schema []kvx.SchemaField) error {
	args := []string{indexName, "ON", "HASH", "PREFIX", "1", prefix, "SCHEMA"}

	for _, f := range schema {
		args = append(args, f.Name, string(f.Type))
		if f.Sortable {
			args = append(args, "SORTABLE")
		}
	}

	return a.client.Do(ctx, a.client.B().Arbitrary("FT.CREATE").Args(args...).Build()).Error()
}

// DropIndex drops a secondary index.
func (a *Adapter) DropIndex(ctx context.Context, indexName string) error {
	return a.client.Do(ctx, a.client.B().Arbitrary("FT.DROPINDEX").Args(indexName).Build()).Error()
}

// Search performs a search query.
func (a *Adapter) Search(ctx context.Context, indexName string, query string, limit int) ([]string, error) {
	resp := a.client.Do(ctx, a.client.B().Arbitrary("FT.SEARCH").Args(indexName, query, "LIMIT", "0", fmt.Sprintf("%d", limit)).Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	// Use AsFtSearch to parse the response
	total, docs, err := resp.AsFtSearch()
	if err != nil {
		return nil, err
	}
	_ = total // We don't need the total count here

	keys := make([]string, len(docs))
	for i, doc := range docs {
		keys[i] = doc.Key
	}
	return keys, nil
}
