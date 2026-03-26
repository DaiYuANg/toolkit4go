package codec

import (
	"encoding"
	"fmt"
)

type textCodec struct{}

func (c textCodec) Name() string {
	return "text"
}

func (c textCodec) Marshal(v any) ([]byte, error) {
	switch value := v.(type) {
	case string:
		return []byte(value), nil
	case []byte:
		return append([]byte(nil), value...), nil
	case encoding.TextMarshaler:
		data, err := value.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("marshal text: %w", err)
		}
		return data, nil
	case fmt.Stringer:
		return []byte(value.String()), nil
	default:
		return nil, fmt.Errorf("%w: codec=text marshal %T", ErrUnsupportedValue, v)
	}
}

func (c textCodec) Unmarshal(data []byte, v any) error {
	switch target := v.(type) {
	case *string:
		*target = string(data)
		return nil
	case *[]byte:
		*target = append((*target)[:0], data...)
		return nil
	case encoding.TextUnmarshaler:
		if err := target.UnmarshalText(data); err != nil {
			return fmt.Errorf("unmarshal text: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("%w: codec=text unmarshal %T", ErrUnsupportedValue, v)
	}
}

// Text is the built-in text codec.
var Text Codec = textCodec{}
