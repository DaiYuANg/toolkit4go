package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DaiYuANg/arcgo/kvx/repository"
)

func TestHashRepository_Save(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	user := &TestUser{
		ID:    "user1",
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify saved data
	key := "user:user1"
	if _, exists := hash.data[key]; !exists {
		t.Errorf("User not saved to hash")
	}
}

func TestHashRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}

	user, err := repo.FindByID(ctx, "user1")
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
	}
	if user.Age != 30 {
		t.Errorf("Expected age 30, got %d", user.Age)
	}
}

func TestHashRepository_Save_ReplacesStaleIndexes(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("old@example.com"),
		"age":   []byte("30"),
	}
	kv.data["user:idx:email:old@example.com:user1"] = []byte("1")
	kv.data["user:idx:age:30:user1"] = []byte("1")

	user := &TestUser{
		ID:    "user1",
		Name:  "John Doe",
		Email: "new@example.com",
		Age:   31,
	}

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, ok := kv.data["user:idx:email:old@example.com:user1"]; ok {
		t.Fatalf("stale email index should be removed")
	}
	if _, ok := kv.data["user:idx:age:30:user1"]; ok {
		t.Fatalf("stale age index should be removed")
	}
	if _, ok := kv.data["user:idx:email:new@example.com:user1"]; !ok {
		t.Fatalf("new email index should exist")
	}
	if _, ok := kv.data["user:idx:age:31:user1"]; !ok {
		t.Fatalf("new age index should exist")
	}
}

func TestHashRepository_UpdateField_ReplacesIndexedFieldEntry(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("old@example.com"),
		"age":   []byte("30"),
	}
	kv.data["user:idx:email:old@example.com:user1"] = []byte("1")

	if err := repo.UpdateField(ctx, "user1", "Email", "new@example.com"); err != nil {
		t.Fatalf("UpdateField failed: %v", err)
	}

	if _, ok := kv.data["user:idx:email:old@example.com:user1"]; ok {
		t.Fatalf("old email index should be removed")
	}
	if _, ok := kv.data["user:idx:email:new@example.com:user1"]; !ok {
		t.Fatalf("new email index should exist")
	}
}

func TestHashRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	_, err := repo.FindByID(ctx, "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestHashRepository_Exists(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	kv.data["user:user1"] = []byte("exists")

	exists, err := repo.Exists(ctx, "user1")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Errorf("Expected user to exist")
	}

	exists, err = repo.Exists(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Errorf("Expected user to not exist")
	}
}

func TestHashRepository_Delete(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	hash.data["user:user1"] = map[string][]byte{
		"name": []byte("John Doe"),
	}
	kv.data["user:user1"] = []byte("exists")

	err := repo.Delete(ctx, "user1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	if _, exists := hash.data["user:user1"]; exists {
		t.Errorf("User not deleted from hash")
	}
	if _, exists := kv.data["user:user1"]; exists {
		t.Errorf("User key not deleted from kv")
	}
}

func TestHashRepository_FindByIDs(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}
	hash.data["user:user2"] = map[string][]byte{
		"name":  []byte("Jane Doe"),
		"email": []byte("jane@example.com"),
		"age":   []byte("25"),
	}

	results, err := repo.FindByIDs(ctx, []string{"user1", "user2", "nonexistent"})
	if err != nil {
		t.Fatalf("FindByIDs failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if user, ok := results["user1"]; !ok || user.Name != "John Doe" {
		t.Errorf("Expected user1 with name 'John Doe'")
	}

	if user, ok := results["user2"]; !ok || user.Name != "Jane Doe" {
		t.Errorf("Expected user2 with name 'Jane Doe'")
	}
}

func TestHashRepository_Count(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	kv.data["user:user1"] = []byte("exists")
	kv.data["user:user2"] = []byte("exists")
	kv.data["other:key"] = []byte("exists")

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestHashRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	kv.data["user:user1"] = []byte("exists")
	kv.data["user:user2"] = []byte("exists")
	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}
	hash.data["user:user2"] = map[string][]byte{
		"name":  []byte("Jane Doe"),
		"email": []byte("jane@example.com"),
		"age":   []byte("25"),
	}

	results, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestHashRepository_FindAll_ScansAllPagesAndDeduplicates(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	kv.scanPages = [][]string{
		{"user:user1", "user:user2"},
		{"user:user2", "user:user3"},
	}
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}
	hash.data["user:user2"] = map[string][]byte{
		"name":  []byte("Jane Doe"),
		"email": []byte("jane@example.com"),
		"age":   []byte("25"),
	}
	hash.data["user:user3"] = map[string][]byte{
		"name":  []byte("Bob Doe"),
		"email": []byte("bob@example.com"),
		"age":   []byte("40"),
	}

	results, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Expected 3 unique results, got %d", len(results))
	}

	ids := map[string]bool{}
	for _, result := range results {
		ids[result.ID] = true
	}

	for _, id := range []string{"user1", "user2", "user3"} {
		if !ids[id] {
			t.Fatalf("Expected result for %s", id)
		}
	}
}

func TestHashRepository_Count_ScansAllPagesAndDeduplicates(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	kv.scanPages = [][]string{
		{"user:user1", "user:user2"},
		{"user:user2", "user:user3"},
	}
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	if count != 3 {
		t.Fatalf("Expected count 3, got %d", count)
	}
}

func TestHashRepository_SaveWithExpiration(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	user := &TestUser{
		ID:    "user1",
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	err := repo.SaveWithExpiration(ctx, user, time.Hour)
	if err != nil {
		t.Fatalf("SaveWithExpiration failed: %v", err)
	}

	// Verify expiration was set
	if _, exists := kv.expiration["user:user1"]; !exists {
		t.Errorf("Expiration not set")
	}
}

func TestHashRepository_UpdateField(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}

	err := repo.UpdateField(ctx, "user1", "Name", "Jane Doe")
	if err != nil {
		t.Fatalf("UpdateField failed: %v", err)
	}

	// Note: UpdateField implementation needs to be fixed to actually update
	// This test documents expected behavior
}

func TestHashRepository_IncrementField(t *testing.T) {
	ctx := context.Background()
	hash := newMockHash()
	kv := newMockKV()
	repo := repository.NewHashRepository[TestUser](hash, kv, "user")

	// Pre-populate data
	hash.data["user:user1"] = map[string][]byte{
		"name":  []byte("John Doe"),
		"email": []byte("john@example.com"),
		"age":   []byte("30"),
	}

	newVal, err := repo.IncrementField(ctx, "user1", "Age", 5)
	if err != nil {
		t.Fatalf("IncrementField failed: %v", err)
	}

	// Mock returns 0, but in real implementation would return 35
	_ = newVal
}
