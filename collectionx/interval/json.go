package interval

import (
	"fmt"
	"slices"

	common "github.com/DaiYuANg/arcgo/collectionx/internal"
)

// ToJSON serializes normalized ranges to JSON.
func (s *RangeSet[T]) ToJSON() ([]byte, error) {
	if s != nil && !s.jsonDirty && s.jsonCache != nil {
		return slices.Clone(s.jsonCache), nil
	}
	data, err := common.MarshalJSONValue(s.Ranges())
	if err != nil {
		return nil, fmt.Errorf("marshal range set json: %w", err)
	}
	if s != nil {
		s.jsonCache = data
		s.stringCache = string(data)
		s.jsonDirty = false
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (s *RangeSet[T]) MarshalJSON() ([]byte, error) {
	data, err := common.ForwardToJSON(s.ToJSON)
	if err != nil {
		return nil, fmt.Errorf("marshal range set: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (s *RangeSet[T]) String() string {
	if s != nil && !s.jsonDirty && s.stringCache != "" {
		return s.stringCache
	}
	data, err := s.ToJSON()
	return common.JSONResultString(data, err, "[]")
}

// ToJSON serializes range-map entries to JSON.
func (m *RangeMap[T, V]) ToJSON() ([]byte, error) {
	if m != nil && !m.jsonDirty && m.jsonCache != nil {
		return slices.Clone(m.jsonCache), nil
	}
	data, err := common.MarshalJSONValue(m.Entries())
	if err != nil {
		return nil, fmt.Errorf("marshal range map json: %w", err)
	}
	if m != nil {
		m.jsonCache = data
		m.stringCache = string(data)
		m.jsonDirty = false
	}
	return slices.Clone(data), nil
}

// MarshalJSON implements json.Marshaler.
func (m *RangeMap[T, V]) MarshalJSON() ([]byte, error) {
	data, err := common.ForwardToJSON(m.ToJSON)
	if err != nil {
		return nil, fmt.Errorf("marshal range map: %w", err)
	}
	return data, nil
}

// String implements fmt.Stringer.
func (m *RangeMap[T, V]) String() string {
	if m != nil && !m.jsonDirty && m.stringCache != "" {
		return m.stringCache
	}
	data, err := m.ToJSON()
	return common.JSONResultString(data, err, "[]")
}
