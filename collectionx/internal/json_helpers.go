package internal

import "encoding/json"

// MarshalJSONValue serializes a value to JSON bytes.
func MarshalJSONValue(value any) ([]byte, error) {
	return json.Marshal(value)
}

// ForwardToJSON delegates json.Marshaler implementation to ToJSON-style methods.
func ForwardToJSON(toJSON func() ([]byte, error)) ([]byte, error) {
	if toJSON == nil {
		return json.Marshal(nil)
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
