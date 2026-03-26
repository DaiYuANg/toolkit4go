package codec

import (
	"encoding/json"
	"fmt"
)

type jsonCodec struct{}

func (c jsonCodec) Name() string {
	return "json"
}

func (c jsonCodec) Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	return data, nil
}

func (c jsonCodec) Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}

// JSON is the built-in JSON codec.
var JSON Codec = jsonCodec{}
