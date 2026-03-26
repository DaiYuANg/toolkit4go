package internal

import (
	"encoding/json"
	"fmt"
)

// MarshalJSONValue serializes a value to JSON bytes.
func MarshalJSONValue(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal json value: %w", err)
	}
	return data, nil
}

// ForwardToJSON delegates json.Marshaler implementation to ToJSON-style methods.
func ForwardToJSON(toJSON func() ([]byte, error)) ([]byte, error) {
	if toJSON == nil {
		data, err := json.Marshal(nil)
		if err != nil {
			return nil, fmt.Errorf("marshal nil json value: %w", err)
		}
		return data, nil
	}
	return toJSON()
}

// StringFromToJSON converts ToJSON-style methods into fmt.Stringer output.
func StringFromToJSON(toJSON func() ([]byte, error), fallback string) string {
	if toJSON == nil {
		return fallback
	}
	data, err := toJSON()
	return JSONResultString(data, err, fallback)
}
