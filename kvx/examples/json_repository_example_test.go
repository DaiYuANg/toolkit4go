package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DaiYuANg/archgo/kvx"
	"github.com/DaiYuANg/archgo/kvx/repository"
)

type exampleJSONBackend struct {
	data map[string][]byte
}

func newExampleJSONBackend() *exampleJSONBackend {
	return &exampleJSONBackend{data: map[string][]byte{}}
}
func (m *exampleJSONBackend) Get(_ context.Context, key string) ([]byte, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, kvx.ErrNil
	}
	return v, nil
}
func (m *exampleJSONBackend) MGet(context.Context, []string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}
func (m *exampleJSONBackend) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	m.data[key] = value
	return nil
}
func (m *exampleJSONBackend) MSet(context.Context, map[string][]byte, time.Duration) error {
	return nil
}
func (m *exampleJSONBackend) Delete(_ context.Context, key string) error {
	delete(m.data, key)
	return nil
}
func (m *exampleJSONBackend) DeleteMulti(context.Context, []string) error { return nil }
func (m *exampleJSONBackend) Exists(_ context.Context, key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}
func (m *exampleJSONBackend) ExistsMulti(context.Context, []string) (map[string]bool, error) {
	return map[string]bool{}, nil
}
func (m *exampleJSONBackend) Expire(context.Context, string, time.Duration) error { return nil }
func (m *exampleJSONBackend) TTL(context.Context, string) (time.Duration, error)  { return 0, nil }
func (m *exampleJSONBackend) Scan(_ context.Context, pattern string, _ uint64, _ int64) ([]string, uint64, error) {
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys, 0, nil
}
func (m *exampleJSONBackend) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, _, err := m.Scan(ctx, pattern, 0, 0)
	return keys, err
}
func (m *exampleJSONBackend) JSONSet(_ context.Context, key string, _ string, value []byte, _ time.Duration) error {
	m.data[key] = value
	return nil
}
func (m *exampleJSONBackend) JSONGet(_ context.Context, key string, _ string) ([]byte, error) {
	return m.Get(context.Background(), key)
}
func (m *exampleJSONBackend) JSONSetField(_ context.Context, key string, path string, value []byte) error {
	var doc map[string]any
	_ = json.Unmarshal(m.data[key], &doc)
	doc[path[2:]] = string(value)
	encoded, _ := json.Marshal(doc)
	m.data[key] = encoded
	return nil
}
func (m *exampleJSONBackend) JSONGetField(context.Context, string, string) ([]byte, error) {
	return nil, nil
}
func (m *exampleJSONBackend) JSONDelete(_ context.Context, key string, _ string) error {
	delete(m.data, key)
	return nil
}

func ExampleJSONRepository() {
	backend := newExampleJSONBackend()
	repo := repository.NewJSONRepository[ExampleUser](backend, backend, "json:user")
	_ = repo.Save(context.Background(), &ExampleUser{ID: "u-2", Name: "Bob", Email: "bob@example.com"})
	entity, _ := repo.FindByID(context.Background(), "u-2")
	fmt.Println(entity.ID, entity.Email)
	// Output: u-2 bob@example.com
}
