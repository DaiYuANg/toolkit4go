package redis

import (
	"fmt"
)

// ============== Helper Functions ==============

func convertInterfaceMapToBytes(m map[string]interface{}) map[string][]byte {
	result := make(map[string][]byte, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case []byte:
			result[k] = val
		case string:
			result[k] = []byte(val)
		default:
			result[k] = []byte(fmt.Sprintf("%v", val))
		}
	}
	return result
}

func valueToBytes(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case nil:
		return nil, nil
	default:
		return []byte(fmt.Sprintf("%v", v)), nil
	}
}

func parseFTSearchResponse(val interface{}) ([]string, error) {
	arr, ok := val.([]interface{})
	if !ok {
		return nil, nil
	}

	if len(arr) < 1 {
		return nil, nil
	}

	// Extract keys from the response
	var keys []string
	for i := 1; i < len(arr); i += 2 {
		if key, ok := arr[i].(string); ok {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func parseFTAggregateResponse(val interface{}) ([]map[string]interface{}, error) {
	arr, ok := val.([]interface{})
	if !ok {
		return nil, nil
	}

	if len(arr) < 1 {
		return nil, nil
	}

	// Parse aggregation results
	var results []map[string]interface{}
	for i := 1; i < len(arr); i++ {
		if row, ok := arr[i].([]interface{}); ok {
			result := make(map[string]interface{})
			for j := 0; j < len(row)-1; j += 2 {
				if key, ok := row[j].(string); ok {
					result[key] = row[j+1]
				}
			}
			results = append(results, result)
		}
	}
	return results, nil
}
