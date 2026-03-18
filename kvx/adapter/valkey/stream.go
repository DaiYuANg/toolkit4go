package valkey

import (
	"context"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/valkey-io/valkey-go"
	"strconv"
)

// ============== Stream Interface ==============

// XAdd adds an entry to a stream.
func (a *Adapter) XAdd(ctx context.Context, key string, id string, values map[string][]byte) (string, error) {
	// Build the command with FieldValue chain
	cmd := a.client.B().Xadd().Key(key).Id(id).FieldValue()
	for k, v := range values {
		cmd = cmd.FieldValue(k, valkey.BinaryString(v))
	}

	resp := a.client.Do(ctx, cmd.Build())
	if resp.Error() != nil {
		return "", resp.Error()
	}
	return resp.ToString()
}

// XRead reads entries from a stream.
func (a *Adapter) XRead(ctx context.Context, key string, start string, count int64) ([]kvx.StreamEntry, error) {
	var cmd valkey.Completed
	if count > 0 {
		cmd = a.client.B().Xread().Count(count).Block(0).Streams().Key(key).Id(start).Build()
	} else {
		cmd = a.client.B().Xread().Block(0).Streams().Key(key).Id(start).Build()
	}

	resp := a.client.Do(ctx, cmd)
	if resp.Error() != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return nil, nil
		}
		return nil, resp.Error()
	}

	// Parse XREAD response using AsXRead
	xreadResult, err := resp.AsXRead()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.StreamEntry, 0)
	for _, streamEntries := range xreadResult {
		for _, entry := range streamEntries {
			values := make(map[string][]byte)
			for f, v := range entry.FieldValues {
				values[f] = []byte(v)
			}
			entries = append(entries, kvx.StreamEntry{
				ID:     entry.ID,
				Values: values,
			})
		}
	}

	return entries, nil
}

// XRange reads entries in a range.
func (a *Adapter) XRange(ctx context.Context, key string, start, stop string) ([]kvx.StreamEntry, error) {
	resp := a.client.Do(ctx, a.client.B().Xrange().Key(key).Start(start).End(stop).Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}

	// Parse XRANGE response using AsXRange
	xrangeEntries, err := resp.AsXRange()
	if err != nil {
		return nil, err
	}

	entries := make([]kvx.StreamEntry, len(xrangeEntries))
	for i, entry := range xrangeEntries {
		values := make(map[string][]byte)
		for f, v := range entry.FieldValues {
			values[f] = []byte(v)
		}
		entries[i] = kvx.StreamEntry{
			ID:     entry.ID,
			Values: values,
		}
	}

	return entries, nil
}

// XLen gets the number of entries in a stream.
func (a *Adapter) XLen(ctx context.Context, key string) (int64, error) {
	resp := a.client.Do(ctx, a.client.B().Xlen().Key(key).Build())
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	return resp.AsInt64()
}

// XTrim trims the stream to approximately maxLen entries.
func (a *Adapter) XTrim(ctx context.Context, key string, maxLen int64) error {
	return a.client.Do(ctx, a.client.B().Xtrim().Key(key).Maxlen().Threshold(strconv.FormatInt(maxLen, 10)).Build()).Error()
}
