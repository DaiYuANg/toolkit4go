package set

import (
	"fmt"

	common "github.com/DaiYuANg/arcgo/collectionx/internal"
)

// ToJSON serializes set values to JSON.
func (s *Set[T]) ToJSON() ([]byte, error) {
	return marshalSetJSON("set", s.Values())
}

// MarshalJSON implements json.Marshaler.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return forwardSetJSON("set", s.ToJSON)
}

// String implements fmt.Stringer.
func (s *Set[T]) String() string {
	return common.StringFromToJSON(s.ToJSON, "[]")
}

// ToJSON serializes concurrent set values to JSON.
func (s *ConcurrentSet[T]) ToJSON() ([]byte, error) {
	return marshalSetJSON("concurrent set", s.Values())
}

// MarshalJSON implements json.Marshaler.
func (s *ConcurrentSet[T]) MarshalJSON() ([]byte, error) {
	return forwardSetJSON("concurrent set", s.ToJSON)
}

// String implements fmt.Stringer.
func (s *ConcurrentSet[T]) String() string {
	return common.StringFromToJSON(s.ToJSON, "[]")
}

// ToJSON serializes multiset counts to JSON.
func (s *MultiSet[T]) ToJSON() ([]byte, error) {
	return marshalSetJSON("multiset", s.AllCounts())
}

// MarshalJSON implements json.Marshaler.
func (s *MultiSet[T]) MarshalJSON() ([]byte, error) {
	return forwardSetJSON("multiset", s.ToJSON)
}

// String implements fmt.Stringer.
func (s *MultiSet[T]) String() string {
	return common.StringFromToJSON(s.ToJSON, "{}")
}

// ToJSON serializes ordered set values to JSON.
func (s *OrderedSet[T]) ToJSON() ([]byte, error) {
	return marshalSetJSON("ordered set", s.Values())
}

// MarshalJSON implements json.Marshaler.
func (s *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	return forwardSetJSON("ordered set", s.ToJSON)
}

// String implements fmt.Stringer.
func (s *OrderedSet[T]) String() string {
	return common.StringFromToJSON(s.ToJSON, "[]")
}

func marshalSetJSON[T any](kind string, value T) ([]byte, error) {
	data, err := common.MarshalJSONValue(value)
	if err != nil {
		return nil, fmt.Errorf("marshal %s JSON: %w", kind, err)
	}

	return data, nil
}

func forwardSetJSON(kind string, fn func() ([]byte, error)) ([]byte, error) {
	data, err := common.ForwardToJSON(fn)
	if err != nil {
		return nil, fmt.Errorf("marshal %s JSON: %w", kind, err)
	}

	return data, nil
}
